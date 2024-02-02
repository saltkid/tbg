package cmd

import (
	"fmt"
	"github.com/saltkid/tbg/flag"
)

type CmdType uint8

const (
	Run CmdType = iota
	Config
	Add
	Remove
	Edit
	Help
	Version
)

type Cmd struct {
	Type    CmdType
	Value   string
	SubCmds []Cmd
	Flags   []flag.Flag
}

func ToCommand(s string) (*Cmd, error) {
	switch s {
	case "run":
		return &Cmd{
			Type:    Run,
			SubCmds: make([]Cmd, 0),
			Flags:   make([]flag.Flag, 0),
		}, nil
	case "config":
		return &Cmd{
			Type:    Config,
			SubCmds: make([]Cmd, 0),
			Flags:   make([]flag.Flag, 0),
		}, nil
	case "add":
		return &Cmd{
			Type:    Add,
			SubCmds: make([]Cmd, 0),
			Flags:   make([]flag.Flag, 0),
		}, nil
	case "remove":
		return &Cmd{
			Type:    Remove,
			SubCmds: make([]Cmd, 0),
			Flags:   make([]flag.Flag, 0),
		}, nil
	case "edit":
		return &Cmd{
			Type:    Edit,
			SubCmds: make([]Cmd, 0),
			Flags:   make([]flag.Flag, 0),
		}, nil
	case "help":
		return &Cmd{
			Type:    Help,
			SubCmds: make([]Cmd, 0),
			Flags:   make([]flag.Flag, 0),
		}, nil
	case "version":
		return &Cmd{
			Type:    Version,
			SubCmds: make([]Cmd, 0),
			Flags:   make([]flag.Flag, 0),
		}, nil
	default:
		return nil, fmt.Errorf("unknown command: %s", s)
	}
}

func (c *Cmd) ValidateValue(val string) error {
	// TODO
	switch c.Type {
	case Run:
		return nil
	case Config:
		return nil
	case Add:
		return nil
	case Remove:
		return nil
	case Edit:
		return nil
	case Help:
		return nil
	case Version:
		return nil
	default:
		return fmt.Errorf("unknown command: %s", val)
	}
}

func (c *Cmd) ValidateFlag(f flag.Flag) error {
	// TODO
	switch c.Type {
	case Run:
		return nil
	case Config:
		return nil
	case Add:
		return nil
	case Remove:
		return nil
	case Edit:
		return nil
	case Help:
		return nil
	case Version:
		return nil
	default:
		return fmt.Errorf("unexpected error: unknown command type: %d", c.Type)
	}
}

func (c *Cmd) Execute() error {
	// TODO
	switch c.Type {
	case Run:
		return nil
	case Config:
		return nil
	case Add:
		return nil
	case Remove:
		return nil
	case Edit:
		return nil
	case Help:
		return nil
	case Version:
		return nil
	default:
		return fmt.Errorf("unexpected error: unknown command type: %d", c.Type)
	}
}
