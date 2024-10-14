package main

import (
	"fmt"
	"github.com/eiannone/keyboard"
	"time"
)

type TbgState struct {
	Images               []string
	ImageIndex           uint16
	Paths                []ImagesPath
	PathIndex            uint16
	Config               *Config
	ConfigPath           string
	CurrentPathAlignment string
	CurrentPathStretch   string
	CurrentPathOpacity   float32
	Random               bool
	Events               *TbgEvents
	Settings             *WTSettings
}

func (tbg *TbgState) String() string {
	return fmt.Sprint(`TbgState
  ConfigPath: `, tbg.ConfigPath, `
  Config: `, tbg.Config, `
  CurrentPathAlignment: `, tbg.CurrentPathAlignment, `
  CurrentPathStretch: `, tbg.CurrentPathStretch, `
  CurrentPathOpacity: `, tbg.CurrentPathOpacity, `
  Images: `, tbg.Images[0], `...`, tbg.Images[len(tbg.Images)-1], `
  ImageIndex: `, tbg.ImageIndex, `
  PathIndex: `, tbg.PathIndex,
	)
}

type TbgEvents struct {
	Done         chan struct{}
	ImageChanged chan struct{}
}

func NewBackgroundState(config *Config, configPath string, randomFlag bool) (*TbgState, error) {
	wtSettings, err := NewWTSettings()
	if err != nil {
		return nil, err
	}
	return &TbgState{
		Images:               make([]string, 2),
		ImageIndex:           0,
		Paths:                config.Paths,
		PathIndex:            0,
		Config:               config,
		ConfigPath:           configPath,
		CurrentPathAlignment: config.Alignment,
		CurrentPathStretch:   config.Stretch,
		CurrentPathOpacity:   config.Opacity,
		Random:               randomFlag,
		Events: &TbgEvents{
			Done:         make(chan struct{}),
			ImageChanged: make(chan struct{}),
		},
		Settings: wtSettings,
	}, nil
}

func (tbg *TbgState) Start() error {
	err := tbg.Init()
	if err != nil {
		return fmt.Errorf("Failed to initialize tbg: %s", err.Error())
	}
	tbg.Settings.Write(tbg.Images[tbg.ImageIndex],
		tbg.Config.Profile,
		tbg.CurrentPathAlignment,
		tbg.CurrentPathStretch,
		tbg.CurrentPathOpacity,
	)
	tbg.Config.Log(tbg.ConfigPath).RunSettings(
		tbg.Images[tbg.ImageIndex],
		tbg.Config.Profile,
		tbg.Config.Interval,
		tbg.CurrentPathAlignment,
		tbg.CurrentPathStretch,
		tbg.CurrentPathOpacity,
	)
	keysEvents, err := keyboard.GetKeys(10)
	if err != nil {
		return err
	}
	defer func() {
		_ = keyboard.Close()
	}()
	go tbg.readUserInput(keysEvents)
	return tbg.Wait()
}

func (tbg *TbgState) Init() error {
	if len(tbg.Config.Paths) == 0 {
		return fmt.Errorf(`config at "%s" has no paths`, tbg.ConfigPath)
	}
	ShuffleFrom(0, tbg.Paths)
	err := tbg.UpdateCurrentPathState()
	return err
}

func (tbg *TbgState) UpdateCurrentPathState() error {
	currentPath := tbg.Config.Paths[tbg.PathIndex]
	tbg.CurrentPathAlignment = Option(currentPath.Alignment).UnwrapOr(tbg.Config.Alignment)
	tbg.CurrentPathStretch = Option(currentPath.Stretch).UnwrapOr(tbg.Config.Stretch)
	tbg.CurrentPathOpacity = Option(currentPath.Opacity).UnwrapOr(tbg.Config.Opacity)
	var err error
	tbg.Images, err = currentPath.Images()
	ShuffleFrom(0, tbg.Images)
	return err
}

func (tbg *TbgState) readUserInput(keysEvents <-chan keyboard.KeyEvent) {
	for {
		event := <-keysEvents
		if event.Err != nil {
			panic(event.Err)
		}
		switch keyboard.Key(event.Rune) {
		case keyboard.Key('c'):
			commandList()
		case keyboard.Key('n'):
			tbg.NextImage()
			tbg.Events.ImageChanged <- struct{}{}
		case keyboard.Key('N'):
			tbg.NextPath()
			tbg.Events.ImageChanged <- struct{}{}
		case keyboard.Key('p'):
			tbg.PreviousImage()
			tbg.Events.ImageChanged <- struct{}{}
		case keyboard.Key('P'):
			tbg.PreviousPath()
			tbg.Events.ImageChanged <- struct{}{}
		case keyboard.Key('q'):
			fmt.Println("Exiting...")
			close(tbg.Events.Done)
			return
		case keyboard.Key('r'):
			tbg.RandomizeImages()
			tbg.Events.ImageChanged <- struct{}{}
		case keyboard.Key('R'):
			tbg.RandomizePaths()
			tbg.Events.ImageChanged <- struct{}{}
		case keyboard.Key('d'):
			fmt.Println(tbg)
		default:
			fmt.Printf("invalid key '%c' ('c' for list of [c]ommand)\n", event.Rune)
		}
	}
}

func (tbg *TbgState) Wait() error {
	ticker := time.Tick(time.Duration(tbg.Config.Interval) * time.Minute)
	// ticker := time.Tick(time.Second * 10) // for debug purposes
	for {
		select {
		case <-ticker:
			tbg.NextImage()
		case <-tbg.Events.Done:
			fmt.Println("Goodbye!")
			return nil
		case <-tbg.Events.ImageChanged:
			tbg.Settings.Write(tbg.Images[tbg.ImageIndex],
				tbg.Config.Profile,
				tbg.CurrentPathAlignment,
				tbg.CurrentPathStretch,
				tbg.CurrentPathOpacity,
			)
			tbg.Config.Log(tbg.ConfigPath).RunSettings(
				tbg.Images[tbg.ImageIndex],
				tbg.Config.Profile,
				tbg.Config.Interval,
				tbg.CurrentPathAlignment,
				tbg.CurrentPathStretch,
				tbg.CurrentPathOpacity,
			)
		}
	}
}
func (tbg *TbgState) NextImage() {
	fmt.Println("using next image...")
	tbg.ImageIndex++
	if tbg.ImageIndex == uint16(len(tbg.Images)) {
		fmt.Print("no more images. going to next path: ")
		tbg.NextPath()
		return
	}
}
func (tbg *TbgState) PreviousImage() {
	fmt.Println("using previous image...")
	switch tbg.ImageIndex {
	case 0:
		fmt.Print("no more images. going to previous path: ")
		tbg.PreviousPath()
	default:
		tbg.ImageIndex--
	}
}
func (tbg *TbgState) RandomizeImages() {
	fmt.Println("randomizing from current image up to last image...")
	fmt.Println("(previous images will not be randomized so you can go back)")
	ShuffleFrom(int(tbg.ImageIndex), tbg.Images)
}
func (tbg *TbgState) NextPath() {
	fmt.Println("using next dir...")
	tbg.ImageIndex = 0
	tbg.PathIndex++
	if tbg.PathIndex >= uint16(len(tbg.Config.Paths)) {
		fmt.Println("no more next dirs. going to first dir again: ", tbg.Config.Paths[0].Path)
		ShuffleFrom(0, tbg.Paths)
		tbg.PathIndex = 0
	}
	tbg.UpdateCurrentPathState()
}
func (tbg *TbgState) PreviousPath() {
	fmt.Println("using previous path...")
	tbg.ImageIndex = 0
	switch tbg.PathIndex {
	case 0:
		fmt.Println("no more previous dirs. going to last dir again: ", tbg.Config.Paths[len(tbg.Config.Paths)-1].Path)
		ShuffleFrom(0, tbg.Paths)
		tbg.PathIndex = uint16(len(tbg.Config.Paths) - 1)
	default:
		tbg.PathIndex--
	}
	tbg.UpdateCurrentPathState()
}
func (tbg *TbgState) RandomizePaths() {
	fmt.Println("randomizing from current path up to last path...")
	fmt.Println("(previous paths will not be randomized so you can go back)")
	ShuffleFrom(int(tbg.PathIndex), tbg.Config.Paths)
	tbg.UpdateCurrentPathState()
}

func commandList() {
	fmt.Print(`
q: [q]uit
n: [n]ext image
p: [p]revious image
N: [N]ext dir
P: [P]revious dir
r: [r]andomize images (current to last; previous unaffected)
R: [R]andomize dirs (current to last; previous unaffected)
c: [c]ommand list
`)
}
