package main

import (
	"fmt"
	"os"
	"path/filepath"
)

type RemoveCommand struct {
	Path      string
	Alignment bool
	Stretch   bool
	Opacity   bool
}

func (cmd *RemoveCommand) Type() CommandType { return RemoveCommandType }

func (cmd *RemoveCommand) String() {
	fmt.Println("Remove Command")
	fmt.Println("Flags:")
	if cmd.Alignment {
		fmt.Println(" ", AlignmentFlag)
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
	absPath, err := filepath.Abs(*val)
	if err != nil {
		return fmt.Errorf("Failed to get absolute path of %s: %s", *val, err)
	}
	cmd.Path = filepath.ToSlash(absPath)
	return nil
}

func (cmd *RemoveCommand) ValidateFlag(f Flag) error {
	switch f.Type {
	case AlignmentFlag:
		if f.Value != nil {
			return fmt.Errorf("'%s' for 'remove' does not take any arguments. got %s", AlignmentFlag, *f.Value)
		}
		cmd.Alignment = true
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
	configPath, err := ConfigPath()
	if err != nil {
		return err
	}
	yamlFile, err := os.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("Failed to read config file %s: %s", configPath, err)
	}
	configContents := new(Config)
	err = configContents.Unmarshal(yamlFile)
	if err != nil {
		return err
	}
	err = configContents.RemovePath(configPath, cmd.Path, cmd.Alignment, cmd.Stretch, cmd.Opacity)
	if err != nil {
		return err
	}
	return nil
}
