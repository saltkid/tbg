package main

import (
	"fmt"
)

type RunCommand struct {
	Alignment *string
	// path to a custom config file
	Config   *string
	Interval *uint16
	Opacity  *float32
	Port     *uint16
	Profile  *string
	Stretch  *string
}

func (cmd *RunCommand) Type() CommandType { return RunCommandType }

func (cmd *RunCommand) String() {
	fmt.Println("Run Command")
	fmt.Println("Flags:")
	if cmd.Alignment != nil {
		fmt.Println(" ", AlignmentFlag, *cmd.Alignment)
	}
	if cmd.Config != nil {
		fmt.Println(" ", ConfigFlag, *cmd.Config)
	}
	if cmd.Interval != nil {
		fmt.Println(" ", IntervalFlag, *cmd.Interval)
	}
	if cmd.Opacity != nil {
		fmt.Println(" ", OpacityFlag, *cmd.Opacity)
	}
	if cmd.Port != nil {
		fmt.Println(" ", PortFlag, *cmd.Port)
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
	case ConfigFlag:
		val, err := ValidateConfig(f.Value)
		if err != nil {
			return err
		}
		cmd.Config = val
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
	config, configPath, err := ConfigInit(cmd.Config)
	if err != nil {
		return err
	}
	config.Profile = Option(cmd.Profile).Or(config.Profile).val
	config.Interval = Option(cmd.Interval).Or(config.Interval).val
	config.Port = Option(cmd.Port).Or(config.Port).val
	tbgState, err := NewTbgState(
		config,
		configPath,
		cmd.Alignment,
		cmd.Opacity,
		cmd.Stretch,
	)
	if err != nil {
		return err
	}
	return tbgState.Start()
}
