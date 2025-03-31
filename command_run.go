package main

import (
	"fmt"
	"os"
)

type RunCommand struct {
	Profile   *string
	Interval  *uint16
	Alignment *string
	Stretch   *string
	Opacity   *float32
	Port      *uint16
}

func (cmd *RunCommand) Type() CommandType { return RunCommandType }

func (cmd *RunCommand) String() {
	fmt.Println("Run Command")
	fmt.Println("Flags:")
	if cmd.Alignment != nil {
		fmt.Println(" ", AlignmentFlag, *cmd.Alignment)
	}
	if cmd.Interval != nil {
		fmt.Println(" ", IntervalFlag, *cmd.Interval)
	}
	if cmd.Opacity != nil {
		fmt.Println(" ", OpacityFlag, *cmd.Opacity)
	}
	if cmd.Port != nil {
		fmt.Println(" ", ProfileFlag, *cmd.Profile)
	}
	if cmd.Profile != nil {
		fmt.Println(" ", ProfileFlag, *cmd.Profile)
	}
	if cmd.Stretch != nil {
		fmt.Println(" ", StretchFlag, *cmd.Stretch)
	}
}

func (cmd *RunCommand) ValidateValue(val *string) error {
	if val == nil || *val == "" {
		return nil
	}
	return fmt.Errorf("'run' takes no arguments. got: '%s'", *val)
}

func (cmd *RunCommand) ValidateFlag(f Flag) error {
	switch f.Type {
	case AlignmentFlag:
		val, err := ValidateAlignment(f.Value)
		if err != nil {
			return err
		}
		cmd.Alignment = val
	case IntervalFlag:
		val, err := ValidateInterval(f.Value)
		if err != nil {
			return err
		}
		cmd.Interval = val
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
	case ProfileFlag:
		val, err := ValidateProfile(f.Value)
		if err != nil {
			return err
		}
		cmd.Profile = val
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

func (cmd *RunCommand) ValidateSubCommand(sc Command) error {
	switch sc.Type() {
	case NoCommandType:
		return nil
	default:
		return fmt.Errorf("'run' takes no sub commands. got: '%s'", sc.Type())
	}
}

func (cmd *RunCommand) Execute() error {
	configPath, err := ConfigPath()
	if err != nil {
		return err
	}
	yamlFile, err := os.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("Failed to read config at %s: %s", shrinkHome(configPath), err)
	}
	config := new(Config)
	err = config.Unmarshal(yamlFile)
	if err != nil {
		return err
	}
	alignment, stretch, opacity := config.determineExecutionFlags(cmd)
	backgroundState, err := NewTbgState(config, configPath, alignment, stretch, opacity)
	if err != nil {
		return err
	}
	return backgroundState.Start(cmd.Port)
}

func (config *Config) determineExecutionFlags(cmd *RunCommand) (string, string, float32) {
	config.Profile = Option(cmd.Profile).Or(config.Profile).val
	config.Interval = Option(cmd.Interval).Or(config.Interval).val
	return Option(cmd.Alignment).UnwrapOr(DefaultAlignment),
		Option(cmd.Stretch).UnwrapOr(DefaultStretch),
		Option(cmd.Opacity).UnwrapOr(DefaultOpacity)
}
