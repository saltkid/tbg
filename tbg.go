package main

import (
	"encoding/json"
	"fmt"
	"math/rand/v2"
	"net/http"
	"strconv"
	"time"
)

type TbgState struct {
	// Slice of images under a path to randomly choose the next image from.
	//
	// Dynamically generated by ImagesPath.images() in TbgState.updateCurrentPathState()
	Images []string
	// slice of paths defined in the tbg config (.tbg.yml)
	Paths []ImagesPath
	// tbg config where paths, interval, and profile information is from
	Config *Config
	// Used for logging the config with the current execution state
	ConfigPath string
	// Events for TbgState goroutines to communicate with each other
	Events *TbgEvents
	// Used to call the WTSettings.Write() method to update WT's settings.json
	// with the current background image
	Settings *WTSettings
}

func (tbg *TbgState) String() string {
	return fmt.Sprint(`TbgState
  ConfigPath: `, tbg.ConfigPath, `
  Config: `, tbg.Config, `
  Images: `, tbg.Images[0], `...`, tbg.Images[len(tbg.Images)-1],
	)
}

// Events for TbgState goroutines to communicate with each other
type TbgEvents struct {
	Done      chan struct{}
	NextImage chan NextImageEvent
	SetImage  chan SetImageEvent
	// all TbgState errors must be routed here. The only method that's allowed
	// to return an error is TbgState.eventHandler() which handles the errors
	// as well
	Error chan error
}

type NextImageEvent struct {
	Alignment *string
	Opacity   *float32
	Stretch   *string
}

type SetImageEvent struct {
	Path      string
	Alignment *string
	Opacity   *float32
	Stretch   *string
}

func NewTbgState(config *Config, configPath string, alignment string, stretch string, opacity float32) (*TbgState, error) {
	wtSettings, err := NewWTSettings()
	if err != nil {
		return nil, err
	}
	return &TbgState{
		Images:     make([]string, 2),
		Paths:      config.Paths,
		Config:     config,
		ConfigPath: configPath,
		Events: &TbgEvents{
			Done:      make(chan struct{}),
			NextImage: make(chan NextImageEvent),
			SetImage:  make(chan SetImageEvent),
			Error:     make(chan error),
		},
		Settings: wtSettings,
	}, nil
}

func (tbg *TbgState) Start(port *uint16) error {
	if len(tbg.Config.Paths) == 0 {
		return fmt.Errorf(`config at "%s" has no paths`, tbg.ConfigPath)
	}
	go tbg.imageUpdateTicker()
	go tbg.startServer(port)
	return tbg.eventHandler()
}

// Creates a ticker that emits a NextImage Event every *interval* minutes where
// interval is defined in the tbg config (.tbg.yml)
func (tbg *TbgState) imageUpdateTicker() {
	ticker := time.Tick(time.Duration(tbg.Config.Interval) * time.Minute)
	for {
		select {
		case <-ticker:
			tbg.Events.NextImage <- NextImageEvent{
				Alignment: nil,
				Opacity:   nil,
				Stretch:   nil,
			}
		}
	}
}

// may emit TbgState.Events.Error (e.g. port is taken)
func (tbg *TbgState) startServer(port *uint16) {
	http.HandleFunc("POST /next-image", func(w http.ResponseWriter, r *http.Request) {
		var reqBody NextImageRequestBody
		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
			tbg.Events.Error <- fmt.Errorf("Failed to decode request body: %s", err)
			return
		}
		tbg.Events.NextImage <- NextImageEvent{
			Alignment: reqBody.Alignment,
			Opacity:   reqBody.Opacity,
			Stretch:   reqBody.Stretch,
		}
		fmt.Fprint(w, "next-image: changed image successfully")
	})
	http.HandleFunc("POST /set-image", func(w http.ResponseWriter, r *http.Request) {
		var reqBody SetImageRequestBody
		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
			tbg.Events.Error <- fmt.Errorf("Failed to decode request body: %s", err)
			return
		}
		tbg.Events.SetImage <- SetImageEvent{
			Path:      reqBody.Path,
			Alignment: reqBody.Alignment,
			Opacity:   reqBody.Opacity,
			Stretch:   reqBody.Stretch,
		}
		fmt.Fprint(w, "set-image: changed image successfully")
	})
	http.HandleFunc("POST /quit", func(w http.ResponseWriter, _ *http.Request) {
		fmt.Fprint(w, "quit: stopped server successfully. Goodbye!")
		close(tbg.Events.Done)
	})
	tbgPort := ":" + strconv.FormatUint(uint64(Option(port).UnwrapOr(tbg.Config.Port)), 10)
	err := http.ListenAndServe(tbgPort, nil)
	if err != nil {
		tbg.Events.Error <- err
	}
}

// Handles events emitted by various TbgState methods.
func (tbg *TbgState) eventHandler() error {
	for {
		select {
		case <-tbg.Events.Done:
			fmt.Println("Goodbye!")
			return nil
		case err := <-tbg.Events.Error:
			return err
		case evt := <-tbg.Events.NextImage:
			if err := tbg.changeToRandomImage(evt.Alignment, evt.Opacity, evt.Stretch); err != nil {
				return err
			}
		case evt := <-tbg.Events.SetImage:
			err := tbg.setImage(
				evt.Path,
				Option(evt.Alignment).UnwrapOr(DefaultAlignment),
				Option(evt.Opacity).UnwrapOr(DefaultOpacity),
				Option(evt.Stretch).UnwrapOr(DefaultStretch),
			)
			if err != nil {
				return err
			}
		}
	}
}

// Changes the background image to a randomly chosen image from images in dirs
// under "paths" in the tbg config file (.tbg.yml)
func (tbg *TbgState) changeToRandomImage(
	alignment *string,
	opacity *float32,
	stretch *string,
) error {
	currentImage, currentAlignment, currentOpacity, currentStretch, err := tbg.randomImage()
	if err != nil {
		return err
	}
	currentAlignment = Option(alignment).UnwrapOr(currentAlignment)
	currentOpacity = Option(opacity).UnwrapOr(currentOpacity)
	currentStretch = Option(stretch).UnwrapOr(currentStretch)
	return tbg.setImage(currentImage, currentAlignment, currentOpacity, currentStretch)
}

// Sets the passed in image path with its properties as the current background
// image
func (tbg *TbgState) setImage(
	imagePath string,
	alignment string,
	opacity float32,
	stretch string,
) error {
	err := tbg.Settings.Write(
		imagePath,
		tbg.Config.Profile,
		alignment,
		opacity,
		stretch,
	)
	if err != nil {
		return err
	}
	tbg.Config.Log(tbg.ConfigPath).RunSettings(
		imagePath,
		tbg.Config.Profile,
		tbg.Config.Interval,
		alignment,
		opacity,
		stretch,
	)
	return nil
}

// Selects a random image from dirs in "paths" field set in tbg config
// (.tbg.yml)
func (tbg *TbgState) randomImage() (string, string, float32, string, error) {
	randomPath := tbg.Config.Paths[uint16(rand.IntN(len(tbg.Config.Paths)))]
	var err error
	tbg.Images, err = randomPath.Images()
	if err != nil {
		return "", "", 0.0, "", err
	}
	return tbg.Images[uint16(rand.IntN(len(tbg.Images)))],
		Option(randomPath.Alignment).UnwrapOr(DefaultAlignment),
		Option(randomPath.Opacity).UnwrapOr(DefaultOpacity),
		Option(randomPath.Stretch).UnwrapOr(DefaultStretch),
		nil
}
