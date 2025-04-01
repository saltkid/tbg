package main

import (
	"fmt"
)

type ConfigCommand struct {
	// path to a custom config path
	Config   *string
	Interval *uint16
	Port     *uint16
	Profile  *string
}

func (cmd *ConfigCommand) Type() CommandType { return ConfigCommandType }

func (cmd *ConfigCommand) String() {
	fmt.Println("Config Command:", cmd.Type())
	fmt.Println("Config File:", Option(cmd.Config).UnwrapOr("default config"))
	fmt.Println("Flags:")
	if cmd.Interval != nil {
		fmt.Println(" ", IntervalFlag, *cmd.Interval)
	}
	if cmd.Port != nil {
		fmt.Println(" ", PortFlag, *cmd.Port)
	}
	if cmd.Profile != nil {
		fmt.Println(" ", ProfileFlag, *cmd.Profile)
	}
}

func (cmd *ConfigCommand) ValidateValue(val *string) error {
	if val == nil || *val == "" {
		// use default config
		return nil
	}
	absPath, err := ValidateConfig(val)
	if err != nil {
		return err
	}
	cmd.Config = absPath
	return nil
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
	config, configPath, err := ConfigInit(cmd.Config)
	if err != nil {
		return err
	}
	isEditingConfig := cmd.Profile != nil || cmd.Interval != nil || cmd.Port != nil
	if isEditingConfig {
		return config.EditConfig(configPath, cmd.Interval, cmd.Port, cmd.Profile)
	}
	config.Log(configPath)
	return nil
}
