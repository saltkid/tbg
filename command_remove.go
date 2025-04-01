package main

import (
	"fmt"
)

type RemoveCommand struct {
	// raw input of user which may or may not have ~ and environment
	// variables, both of which will be kept unexpanded.
	Path string
	// normalized path from user input which expands both ~ and environment
	// variables.
	CleanPath string
	Alignment bool
	// path to a custom config path
	Config  *string
	Opacity bool
	Stretch bool
}

func (cmd *RemoveCommand) Type() CommandType { return RemoveCommandType }

func (cmd *RemoveCommand) String() {
	fmt.Println("Remove Command")
	fmt.Println("Flags:")
	if cmd.Alignment {
		fmt.Println(" ", AlignmentFlag)
	}
	if cmd.Config != nil {
		fmt.Println(" ", ConfigFlag, *cmd.Config)
	}
	if cmd.Opacity {
		fmt.Println(" ", OpacityFlag)
	}
	if cmd.Stretch {
		fmt.Println(" ", StretchFlag)
	}
}

func (cmd *RemoveCommand) ValidateValue(val *string) error {
	if val == nil {
		return fmt.Errorf("'remove' must have an argument. got none")
	}
	absPath, err := NormalizePath(*val)
	if err != nil {
		return fmt.Errorf("Failed to normalize path %s: %s", *val, err)
	}
	cmd.Path = *val
	cmd.CleanPath = absPath
	return nil
}

func (cmd *RemoveCommand) ValidateFlag(f Flag) error {
	switch f.Type {
	case AlignmentFlag:
		if f.Value != nil {
			return fmt.Errorf("'%s' for 'remove' does not take any arguments. got %s", AlignmentFlag, *f.Value)
		}
		cmd.Alignment = true
	case ConfigFlag:
		val, err := ValidateConfig(f.Value)
		if err != nil {
			return err
		}
		cmd.Config = val
	case OpacityFlag:
		if f.Value != nil {
			return fmt.Errorf("'%s' for 'remove' does not take any arguments. got %s", OpacityFlag, *f.Value)
		}
		cmd.Opacity = true
	case StretchFlag:
		if f.Value != nil {
			return fmt.Errorf("'%s' for 'remove' does not take any arguments. got %s", StretchFlag, *f.Value)
		}
		cmd.Stretch = true
	default:
		return fmt.Errorf("invalid flag for 'remove': '%s'", f.Type)
	}
	return nil
}

func (cmd *RemoveCommand) ValidateSubCommand(sc Command) error {
	switch sc.Type() {
	case NoCommandType:
		return nil
	default:
		return fmt.Errorf("'remove' takes no sub commands. got: '%s'", sc.Type())
	}
}

func (cmd *RemoveCommand) Execute() error {
	config, configPath, err := ConfigInit(cmd.Config)
	if err != nil {
		return err
	}
	return config.RemovePath(configPath, cmd.Path, cmd.CleanPath, cmd.Alignment, cmd.Opacity, cmd.Stretch)
}
