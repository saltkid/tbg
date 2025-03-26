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
	case "next-image":
		return new(NextImageCommand), nil
	case "quit":
		return new(QuitCommand), nil
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
	NextImageCommandType
	QuitCommandType
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
	case NextImageCommandType:
		return "next-image"
	case QuitCommandType:
		return "quit"
	default:
		return fmt.Sprintf("UNKNOWN COMMAND '%d'", c)
	}
}

// converts a command type to a command struct
//
// This is guaranteed to not return nil if the token's CommandType come from
// TokenizeArgs
func (c CommandType) ToCommand() Command {
	switch c {
	case RunCommandType:
		return new(RunCommand)
	case ConfigCommandType:
		return new(ConfigCommand)
	case AddCommandType:
		return new(AddCommand)
	case RemoveCommandType:
		return new(RemoveCommand)
	case HelpCommandType:
		return new(HelpCommand)
	case VersionCommandType:
		return new(VersionCommand)
	case NextImageCommandType:
		return new(NextImageCommand)
	case QuitCommandType:
		return new(QuitCommand)
	default: // case: NoCommandType
		return nil
	}
}
