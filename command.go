package main

import (
	"fmt"
)

type Command interface {
	// returns the command type of the command struct
	Type() CommandType
	// prints out debug information
	Debug()
	//
	ValidateValue(val *string) error
	ValidateFlag(f Flag) error
	ValidateSubCommand(cmd Command) error
	Execute() error
}

// converts a string to a command struct (case sensitive)
func ToCommand(s string) (Command, error) {
	switch s {
	case "run":
		return new(RunCommand), nil
	case "config":
		return new(ConfigCommand), nil
	case "add":
		return new(AddCommand), nil
	case "remove":
		return new(RemoveCommand), nil
	case "help":
		return new(HelpCommand), nil
	case "version":
		return new(VersionCommand), nil
	default:
		return nil, fmt.Errorf("unknown command: %s", s)
	}
}

type CommandType uint8

const (
	NoCommandType CommandType = iota
	RunCommandType
	ConfigCommandType
	AddCommandType
	RemoveCommandType
	HelpCommandType
	VersionCommandType
)

func (c CommandType) String() string {
	switch c {
	case NoCommandType:
		return "none"
	case RunCommandType:
		return "run"
	case ConfigCommandType:
		return "config"
	case AddCommandType:
		return "add"
	case RemoveCommandType:
		return "remove"
	case HelpCommandType:
		return "help"
	case VersionCommandType:
		return "version"
	default:
		return fmt.Sprintf("UNKNOWN COMMAND '%d'", c)
	}
}

// converts a command type to a command struct
func (c CommandType) ToCommand() (Command, error) {
	switch c {
	case RunCommandType:
		return new(RunCommand), nil
	case ConfigCommandType:
		return new(ConfigCommand), nil
	case AddCommandType:
		return new(AddCommand), nil
	case RemoveCommandType:
		return new(RemoveCommand), nil
	case HelpCommandType:
		return new(HelpCommand), nil
	case VersionCommandType:
		return new(VersionCommand), nil
	default: // case: NoCommandType
		return nil, fmt.Errorf("Cannot convert CommandType NoCommandType to a Command!")
	}
}
