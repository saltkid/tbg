package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
)

type SetImageCommand struct {
	Path      string
	Alignment *string
	Opacity   *float32
	Stretch   *string
}

func (cmd *SetImageCommand) Type() CommandType { return SetImageCommandType }

func (r *SetImageCommand) Debug() {
	fmt.Println("Set Image Command:", r.Type())
}

func (cmd *SetImageCommand) ValidateValue(val *string) error {
	if val == nil {
		return fmt.Errorf("'set-image' must have an argument. got none")
	}
	absPath, err := filepath.Abs(*val)
	if err != nil {
		return fmt.Errorf("Failed to get absolute path of %s: %s", *val, err)
	}
	if _, err := os.Stat(absPath); os.IsNotExist(err) {
		return fmt.Errorf("%s does not exist: %s", *val, err.Error())
	}
	if !IsImageFile(absPath) {
		return fmt.Errorf("Not an image file: %s", *val)
	}
	cmd.Path = filepath.ToSlash(absPath)
	return nil
}

func (cmd *SetImageCommand) ValidateFlag(f Flag) error {
	switch f.Type {
	case AlignmentFlag:
		val, err := ValidateAlignment(f.Value)
		if err != nil {
			return err
		}
		cmd.Alignment = val
	case OpacityFlag:
		val, err := ValidateOpacity(f.Value)
		if err != nil {
			return err
		}
		cmd.Opacity = val
	case StretchFlag:
		val, err := ValidateStretch(f.Value)
		if err != nil {
			return err
		}
		cmd.Stretch = val
	default:
		return fmt.Errorf("invalid flag for 'next-image': '%s'", f.Type)
	}
	return nil
}

func (cmd *SetImageCommand) ValidateSubCommand(sc Command) error {
	switch sc.Type() {
	case NoCommandType:
		return nil
	default:
		return fmt.Errorf("'next-image' takes no sub commands. got: '%s'", sc.Type())
	}
}

type SetImageRequestBody struct {
	Path      string   `json:"path"`
	Alignment *string  `json:"alignment,omitempty"`
	Opacity   *float32 `json:"opacity,omitempty"`
	Stretch   *string  `json:"stretch,omitempty"`
}

func (cmd *SetImageCommand) Execute() error {
	configPath, err := ConfigPath()
	if err != nil {
		return err
	}
	yamlFile, err := os.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("Failed to read config file %s: %s", configPath, err)
	}
	config := new(Config)
	err = config.Unmarshal(yamlFile)
	if err != nil {
		return err
	}
	setImageArgs := SetImageRequestBody{
		Path:      cmd.Path,
		Alignment: cmd.Alignment,
		Stretch:   cmd.Stretch,
		Opacity:   cmd.Opacity,
	}
	reqBody, err := json.Marshal(setImageArgs)
	if err != nil {
		return fmt.Errorf("failed to marshal request body: %s", err)
	}
	url := fmt.Sprintf("http://127.0.0.1:%d/set-image", config.Port)
	resp, err := http.Post(url, "application/json", bytes.NewReader(reqBody))
	if err != nil {
		return err
	}
	resp.Body.Close()
	return nil
}
