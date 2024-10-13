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
	Random    *bool
}

func (cmd *RunCommand) Type() CommandType { return RunCommandType }

func (cmd *RunCommand) Debug() {
	fmt.Println("Run Command")
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
	if cmd.Random != nil {
		fmt.Println(" ", RandomFlag, *cmd.Random)
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
	case RandomFlag:
		val, err := ValidateRandom(f.Value)
		if err != nil {
			return err
		}
		cmd.Random = val
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
		return fmt.Errorf("Failed to read config file %s: %s", configPath, err)
	}
	config := new(Config)
	err = config.Unmarshal(yamlFile)
	if err != nil {
		return err
	}
	err = configContents.ChangeBgImage(configPath, profile, interval, alignment, stretch, opacity, random)
	if err != nil {
		return err
	}
	return nil
}

func (config *Config) determineExecutionFlags(cmd *RunCommand) bool {
	config.Profile = Option(cmd.Profile).UnwrapOr(config.Profile)
	config.Interval = Option(cmd.Interval).UnwrapOr(config.Interval)
	config.Alignment = Option(cmd.Alignment).UnwrapOr(config.Alignment)
	config.Stretch = Option(cmd.Stretch).UnwrapOr(config.Stretch)
	config.Opacity = Option(cmd.Opacity).UnwrapOr(config.Opacity)
	return Option(cmd.Random).UnwrapOr(false)
}
