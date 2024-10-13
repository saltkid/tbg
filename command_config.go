package main

import (
	"fmt"
	"os"
)

type ConfigCommand struct {
	Profile   *string
	Interval  *uint16
	Alignment *string
	Stretch   *string
	Opacity   *float32
}

func (cmd *ConfigCommand) Type() CommandType { return ConfigCommandType }

func (cmd *ConfigCommand) Debug() {
	fmt.Println("Config Command:", cmd.Type())
	fmt.Println("Flags:")
	if cmd.Profile != nil {
		fmt.Println(" ", ProfileFlag, *cmd.Profile)
	}
	if cmd.Interval != nil {
		fmt.Println(" ", IntervalFlag, *cmd.Interval)
	}
	if cmd.Alignment != nil {
		fmt.Println(" ", AlignmentFlag, *cmd.Alignment)
	}
	if cmd.Stretch != nil {
		fmt.Println(" ", StretchFlag, *cmd.Stretch)
	}
	if cmd.Opacity != nil {
		fmt.Println(" ", OpacityFlag, *cmd.Opacity)
	}
}

func (cmd *ConfigCommand) ValidateValue(val *string) error {
	if val == nil || *val == "" {
		return nil
	}
	return fmt.Errorf("'run' takes no arguments. got: '%s'", *val)
}

func (cmd *ConfigCommand) ValidateFlag(f Flag) error {
	switch f.Type {
	case ProfileFlag:
		val, err := ValidateProfile(f.Value)
		if err != nil {
			return err
		}
		cmd.Profile = val
	case IntervalFlag:
		val, err := ValidateInterval(f.Value)
		if err != nil {
			return err
		}
		cmd.Interval = val
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
		return fmt.Errorf("invalid flag for 'run': '%s'", f.Type)
	}
	return nil
}

func (cmd *ConfigCommand) ValidateSubCommand(sc Command) error {
	switch sc.Type() {
	case NoCommandType:
		return nil
	default:
		return fmt.Errorf("'config' takes no sub commands. got: '%s'", sc.Type())
	}
}

func (cmd *ConfigCommand) Execute() error {
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
	isEditingConfig := cmd.Profile != nil || cmd.Interval != nil || cmd.Alignment != nil || cmd.Stretch != nil || cmd.Opacity != nil
	if isEditingConfig {
		config.EditConfig(configPath, cmd.Profile, cmd.Interval, cmd.Alignment, cmd.Stretch, cmd.Opacity)
	} else {
		config.Log(configPath)
	}
	return nil
}
