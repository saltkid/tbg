package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type NextImageCommand struct {
	Alignment *string
	Opacity   *float32
	Stretch   *string
	Port      *uint16
}

func (cmd *NextImageCommand) Type() CommandType { return NextImageCommandType }

func (r *NextImageCommand) String() {
	fmt.Println("Next Image Command:", r.Type())
}

func (cmd *NextImageCommand) ValidateValue(val *string) error {
	if val == nil || *val == "" {
		return nil
	}
	return fmt.Errorf("'next-image' takes no args. got: '%s'", *val)
}

func (cmd *NextImageCommand) ValidateFlag(f Flag) error {
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
	case PortFlag:
		val, err := ValidatePort(f.Value)
		if err != nil {
			return err
		}
		cmd.Port = val
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

func (cmd *NextImageCommand) ValidateSubCommand(sc Command) error {
	switch sc.Type() {
	case NoCommandType:
		return nil
	default:
		return fmt.Errorf("'next-image' takes no sub commands. got: '%s'", sc.Type())
	}
}

type NextImageRequestBody struct {
	Alignment *string  `json:"alignment,omitempty"`
	Opacity   *float32 `json:"opacity,omitempty"`
	Stretch   *string  `json:"stretch,omitempty"`
}

func (cmd *NextImageCommand) Execute() error {
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
	nextImageArgs := NextImageRequestBody{
		Alignment: cmd.Alignment,
		Stretch:   cmd.Stretch,
		Opacity:   cmd.Opacity,
	}
	reqBody, err := json.Marshal(nextImageArgs)
	if err != nil {
		return fmt.Errorf("failed to marshal request body: %s", err)
	}
	url := fmt.Sprintf("http://127.0.0.1:%d/next-image", Option(cmd.Port).UnwrapOr(config.Port))
	resp, err := http.Post(url, "application/json", bytes.NewReader(reqBody))
	if err != nil {
		return err
	}
	resp.Body.Close()
	return nil
}
