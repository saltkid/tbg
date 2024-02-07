package cmd

import (
	"fmt"
	"github.com/saltkid/tbg/flag"
)

type CmdType uint8

const (
	None CmdType = iota
	Run
	Config
	Add
	Remove
	Edit
	Help
	Version
)

func (c CmdType) ToString() string {
	switch c {
	case Run:
		return "run"
	case Config:
		return "config"
	case Add:
		return "add"
	case Remove:
		return "remove"
	case Edit:
		return "edit"
	case Help:
		return "help"
	case Version:
		return "version"
	default:
		return ""
	}
}

type Cmd struct {
	Type    CmdType
	Value   string
	SubCmds map[CmdType]*Cmd
	Flags   map[flag.FlagType]*flag.Flag
}

func ToCommand(s string) (*Cmd, error) {
	switch s {
	case "run":
		return &Cmd{
			Type:    Run,
			SubCmds: make(map[CmdType]*Cmd, 0),
			Flags:   make(map[flag.FlagType]*flag.Flag, 0),
		}, nil
	case "config":
		return &Cmd{
			Type:    Config,
			SubCmds: make(map[CmdType]*Cmd, 0),
			Flags:   make(map[flag.FlagType]*flag.Flag, 0),
		}, nil
	case "add":
		return &Cmd{
			Type:    Add,
			SubCmds: make(map[CmdType]*Cmd, 0),
			Flags:   make(map[flag.FlagType]*flag.Flag, 0),
		}, nil
	case "remove":
		return &Cmd{
			Type:    Remove,
			SubCmds: make(map[CmdType]*Cmd, 0),
			Flags:   make(map[flag.FlagType]*flag.Flag, 0),
		}, nil
	case "edit":
		return &Cmd{
			Type:    Edit,
			SubCmds: make(map[CmdType]*Cmd, 0),
			Flags:   make(map[flag.FlagType]*flag.Flag, 0),
		}, nil
	case "help":
		return &Cmd{
			Type:    Help,
			SubCmds: make(map[CmdType]*Cmd, 0),
			Flags:   make(map[flag.FlagType]*flag.Flag, 0),
		}, nil
	case "version":
		return &Cmd{
			Type:    Version,
			SubCmds: make(map[CmdType]*Cmd, 0),
			Flags:   make(map[flag.FlagType]*flag.Flag, 0),
		}, nil
	default:
		return nil, fmt.Errorf("unknown command: %s", s)
	}
}

func (c *Cmd) IsNone() bool {
	return c.Type == None
}

func (c *Cmd) ValidateValue(val string) error {
	switch c.Type {
	case Run:
		return RunValidateValue(val)
	case Config:
		return ConfigValidateValue(val)
	case Add:
		return AddValidateValue(val)
	case Remove:
		return RemoveValidateValue(val)
	case Edit:
		return EditValidateValue(val)
	case Help:
		return HelpValidateValue(val)
	case Version:
		return nil
	case None:
		return nil
	default:
		return fmt.Errorf("unknown command: %s", val)
	}
}

func (c *Cmd) ValidateFlag(f *flag.Flag) error {
	switch c.Type {
	case Run:
		return RunValidateFlag(f)
	case Config:
		return ConfigValidateFlag(f)
	case Add:
		return AddValidateFlag(f)
	case Remove:
		return RemoveValidateFlag(f)
	case Edit:
		return EditValidateFlag(f)
	case Help:
		return HelpValidateFlag(f)
	case Version:
		return nil
	case None:
		return nil
	default:
		return fmt.Errorf("unexpected error: unknown command type: %d", c.Type)
	}
}

func (c *Cmd) ValidateSubCmd(sc *Cmd) error {
	switch c.Type {
	case Run:
		return RunValidateSubCmd(sc)
	case Config:
		return ConfigValidateSubCmd(sc)
	case Add:
		return AddValidateSubCmd(sc)
	case Remove:
		return RemoveValidateSubCmd(sc)
	case Edit:
		return EditValidateSubCmd(sc)
	case Help:
		return HelpValidateSubCmd(sc)
	case Version:
		return nil
	case None:
		return nil
	default:
		return fmt.Errorf("unexpected error: unknown command type: %d", sc.Type)
	}
}

func (c *Cmd) Execute() error {
	switch c.Type {
	case Run:
		return RunExecute(c)
	case Config:
		return ConfigExecute(c)
	case Add:
		return AddExecute(c)
	case Remove:
		return RemoveExecute(c)
	case Edit:
		return EditExecute(c)
	case Help:
		return HelpExecute(c)
	case Version:
		return VersionExecute()
	case None:
		return nil
	default:
		return fmt.Errorf("unexpected error: unknown command type: %d", c.Type)
	}
}
