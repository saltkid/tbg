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


		}
	}
}

func readUserInput(keysEvents <-chan keyboard.KeyEvent, done chan<- struct{},
	nextDir chan<- struct{}, nextImage chan<- struct{},
	prevDir chan<- struct{}, prevImage chan<- struct{},
	randomDir chan<- struct{}, randomImage chan<- struct{}) {
	for {
		event := <-keysEvents
		if event.Err != nil {
			panic(event.Err)
		}

		if keyboard.Key(event.Rune) == keyboard.Key('q') {
			fmt.Println("Exiting...")
			close(done)
			return

		} else if keyboard.Key(event.Rune) == keyboard.Key('R') {
			randomDir <- struct{}{}

		} else if keyboard.Key(event.Rune) == keyboard.Key('r') {
			randomImage <- struct{}{}

		} else if keyboard.Key(event.Rune) == keyboard.Key('N') {
			nextDir <- struct{}{}

		} else if keyboard.Key(event.Rune) == keyboard.Key('P') {
			prevDir <- struct{}{}

		} else if keyboard.Key(event.Rune) == keyboard.Key('n') {
			nextImage <- struct{}{}

		} else if keyboard.Key(event.Rune) == keyboard.Key('p') {
			prevImage <- struct{}{}

		} else if keyboard.Key(event.Rune) == keyboard.Key('c') {
			commandList()

		} else {
			fmt.Printf("invalid key '%c' ('h' for help)\n", event.Rune)
		}
	}
}

func fetchImages(dir string) ([]string, error) {
	images := make([]string, 0)
	err := filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() && d.Name() != filepath.Base(dir) {
			return filepath.SkipDir
		}
		if IsImageFile(d.Name()) {
			images = append(images, path)
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("Failed to walk directory %s: %s", dir, err)
	}
	return images, nil
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
