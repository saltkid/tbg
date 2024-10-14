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
	Error        chan error
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
			Error:        make(chan error),
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
	go tbg.readUserInput()
	return tbg.Wait()
}

func (tbg *TbgState) Init() error {
	if len(tbg.Config.Paths) == 0 {
		return fmt.Errorf(`config at "%s" has no paths`, tbg.ConfigPath)
	}
	ShuffleFrom(0, tbg.Paths)
	return tbg.UpdateCurrentPathState()
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

func (tbg *TbgState) readUserInput() {
	keysEvents, err := keyboard.GetKeys(10)
	if err != nil {
		tbg.Events.Error <- err
	}
	defer func() {
		_ = keyboard.Close()
	}()
	for {
		event := <-keysEvents
		if event.Err != nil {
			tbg.Events.Error <- event.Err
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


// Handles events emitted by various TbgState methods
func (tbg *TbgState) Wait() error {
	for {
		ticker := time.Tick(time.Duration(tbg.Config.Interval) * time.Minute)
		// ticker := time.Tick(time.Second * 5) // for debug purposes
		select {
		case <-ticker:
			// in a go routine since NextImage can emit TbgState.Events.Errors
			// and we want to catch that same event in this same loop
			go tbg.NextImage()
		case <-tbg.Events.Done:
			fmt.Println("Goodbye!")
			return nil
		case err := <-tbg.Events.Error:
			return err
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

// emits TbgState.Events.ImageChanged
//
// may emit TbgState.Events.Error through TbgState.NextPath()
func (tbg *TbgState) NextImage() {
	fmt.Println("using next image...")
	tbg.ImageIndex++
	switch tbg.ImageIndex {
	case uint16(len(tbg.Images)):
		fmt.Print("no more images. going to next path: ")
		tbg.NextPath()
	default:
		tbg.Events.ImageChanged <- struct{}{}
	}
}

// emits TbgState.Events.ImageChanged
//
// may emit TbgState.Events.Error through TbgState.PreviousPath()
func (tbg *TbgState) PreviousImage() {
	fmt.Println("using previous image...")
	switch tbg.ImageIndex {
	case 0:
		fmt.Print("no more images. going to previous path: ")
		tbg.PreviousPath()
	default:
		tbg.ImageIndex--
		tbg.Events.ImageChanged <- struct{}{}
	}
}

// emits TbgState.Events.ImageChanged
func (tbg *TbgState) RandomizeImages() {
	fmt.Println("randomizing from current image up to last image...")
	fmt.Println("(previous images will not be randomized so you can go back)")
	ShuffleFrom(int(tbg.ImageIndex), tbg.Images)
	tbg.Events.ImageChanged <- struct{}{}
}

// emits TbgState.Events.ImageChanged
//
// may emit TbgState.Events.Error
func (tbg *TbgState) NextPath() {
	fmt.Println("using next dir...")
	tbg.ImageIndex = 0
	tbg.PathIndex++
	if tbg.PathIndex >= uint16(len(tbg.Config.Paths)) {
		fmt.Println("no more next dirs. going to first dir again: ", tbg.Config.Paths[0].Path)
		ShuffleFrom(0, tbg.Paths)
		tbg.PathIndex = 0
	}
	err := tbg.UpdateCurrentPathState()
	if err != nil {
		tbg.Events.Error <- err
	}
	tbg.Events.ImageChanged <- struct{}{}
}

// emits TbgState.Events.ImageChanged
//
// emits TbgState.Events.Error
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
	err := tbg.UpdateCurrentPathState()
	if err != nil {
		tbg.Events.Error <- err
	}
	tbg.Events.ImageChanged <- struct{}{}
}

// emits TbgState.Events.ImageChanged
//
// may emit TbgState.Events.Error
func (tbg *TbgState) RandomizePaths() {
	fmt.Println("randomizing from current path up to last path...")
	fmt.Println("(previous paths will not be randomized so you can go back)")
	ShuffleFrom(int(tbg.PathIndex), tbg.Config.Paths)
	err := tbg.UpdateCurrentPathState()
	if err != nil {
		tbg.Events.Error <- err
	}
	tbg.Events.ImageChanged <- struct{}{}
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
