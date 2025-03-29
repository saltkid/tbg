package main

import (
	"fmt"
	"os"
)

type ConfigCommand struct {
	Interval *uint16
	Port     *uint16
	Profile  *string
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
}

func (cmd *ConfigCommand) ValidateValue(val *string) error {
	if val == nil || *val == "" {
		return nil
	}
	return fmt.Errorf("'config' takes no arguments. got: '%s'", *val)
}

func (cmd *ConfigCommand) ValidateFlag(f Flag) error {
	switch f.Type {
	case IntervalFlag:
		val, err := ValidateInterval(f.Value)
		if err != nil {
			return err
		}
		cmd.Interval = val
	case PortFlag:
		val, err := ValidatePort(f.Value)
		if err != nil {
			return err
		}
		cmd.Port = val
	case ProfileFlag:
		val, err := ValidateProfile(f.Value)
		if err != nil {
			return err
		}
		cmd.Profile = val
	default:
		return fmt.Errorf("invalid flag for 'config': '%s'", f.Type)
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
	isEditingConfig := cmd.Profile != nil || cmd.Interval != nil || cmd.Port != nil
	if isEditingConfig {
		config.EditConfig(configPath, cmd.Interval, cmd.Port, cmd.Profile)
	} else {
		config.Log(configPath)
	}
	return nil
}
