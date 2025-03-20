package main

import (
	"fmt"
)

type FlagType uint8

const (
	NoFlag FlagType = iota
	AlignmentFlag
	IntervalFlag
	OpacityFlag
	ProfileFlag
	StretchFlag
)

func (f FlagType) String() string {
	switch f {
	case AlignmentFlag:
		return "--alignment"
	case IntervalFlag:
		return "--interval"
	case NoFlag:
		return "none"
	case OpacityFlag:
		return "--opacity"
	case ProfileFlag:
		return "--profile"
	case StretchFlag:
		return "--stretch"
	default:
		return "unknown"
	}
}

type Flag struct {
	Type  FlagType
	Value *string
}

func ToFlag(s string) (*Flag, error) {
	switch s {
	case "--alignment", "-a":
		return &Flag{Type: AlignmentFlag}, nil
	case "--interval", "-i":
		return &Flag{Type: IntervalFlag}, nil
	case "--opacity", "-o":
		return &Flag{Type: OpacityFlag}, nil
	case "--profile", "-p":
		return &Flag{Type: ProfileFlag}, nil
	case "--stretch", "-s":
		return &Flag{Type: StretchFlag}, nil
	default:
		return nil, fmt.Errorf("unknown flag: %s", s)
	}
}
