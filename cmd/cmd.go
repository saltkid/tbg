package cmd

import (
	"fmt"
	"github.com/saltkid/tbg/flag"
)

type CmdType uint8

const (
	_ CmdType = iota
	Run
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
	SubCmds map[CmdType]Cmd
	Flags   map[flag.FlagType]flag.Flag
}

func ToCommand(s string) (*Cmd, error) {
	switch s {
	case "run":
		return &Cmd{
			Type:    Run,
			SubCmds: make(map[CmdType]Cmd, 0),
			Flags:   make(map[flag.FlagType]flag.Flag, 0),
		}, nil
	case "config":
		return &Cmd{
			Type:    Config,
			SubCmds: make(map[CmdType]Cmd, 0),
			Flags:   make(map[flag.FlagType]flag.Flag, 0),
		}, nil
	case "add":
		return &Cmd{
			Type:    Add,
			SubCmds: make(map[CmdType]Cmd, 0),
			Flags:   make(map[flag.FlagType]flag.Flag, 0),
		}, nil
	case "remove":
		return &Cmd{
			Type:    Remove,
			SubCmds: make(map[CmdType]Cmd, 0),
			Flags:   make(map[flag.FlagType]flag.Flag, 0),
		}, nil
	case "edit":
		return &Cmd{
			Type:    Edit,
			SubCmds: make(map[CmdType]Cmd, 0),
			Flags:   make(map[flag.FlagType]flag.Flag, 0),
		}, nil
	case "help":
		return &Cmd{
			Type:    Help,
			SubCmds: make(map[CmdType]Cmd, 0),
			Flags:   make(map[flag.FlagType]flag.Flag, 0),
		}, nil
	case "version":
		return &Cmd{
			Type:    Version,
			SubCmds: make(map[CmdType]Cmd, 0),
			Flags:   make(map[flag.FlagType]flag.Flag, 0),
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
		return ConfigValidateValue(val)
	case Add:
		return AddValidateValue(val)
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
		return ConfigValidateFlag(f)
	case Add:
		return AddValidateFlag(f)
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
		return ConfigExecute(c)
	case Add:
		return AddExecute(c)
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
