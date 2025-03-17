package main

import (
	"fmt"
)

type FlagType uint8

const (
	NoFlag FlagType = iota
	ProfileFlag
	IntervalFlag
	AlignmentFlag
	OpacityFlag
	StretchFlag
	RandomFlag
)

func (f FlagType) String() string {
	switch f {
	case NoFlag:
		return "none"
	case ProfileFlag:
		return "--profile"
	case IntervalFlag:
		return "--interval"
	case AlignmentFlag:
		return "--alignment"
	case OpacityFlag:
		return "--opacity"
	case StretchFlag:
		return "--stretch"
	case RandomFlag:
		return "--random"
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
	case "--profile", "-p":
		return &Flag{Type: ProfileFlag}, nil
	case "--interval", "-i":
		return &Flag{Type: IntervalFlag}, nil
	case "--alignment", "-a":
		return &Flag{Type: AlignmentFlag}, nil
	case "--opacity", "-o":
		return &Flag{Type: OpacityFlag}, nil
	case "--stretch", "-s":
		return &Flag{Type: StretchFlag}, nil
	case "--random", "-r":
		return &Flag{Type: RandomFlag}, nil
	default:
		return nil, fmt.Errorf("unknown flag: %s", s)
	}
}
