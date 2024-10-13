package flag

import (
	"fmt"
	"strconv"
	"strings"
)

type FlagType uint8

const (
	None FlagType = iota
	Profile
	Interval
	Alignment
	Opacity
	Stretch
	Random
)

func (f FlagType) ToString() string {
	switch f {
	case None:
		return "none"
	case Profile:
		return "--profile"
	case Interval:
		return "--interval"
	case Alignment:
		return "--alignment"
	case Opacity:
		return "--opacity"
	case Stretch:
		return "--stretch"
	case Random:
		return "--random"
	default:
		return "unknown"
	}
}

type Flag struct {
	Type  FlagType
	Value string
}

func ToFlag(s string) (*Flag, error) {
	switch s {
	case "--profile", "-p":
		return &Flag{Type: Profile}, nil
	case "--interval", "-i":
		return &Flag{Type: Interval}, nil
	case "--alignment", "-a":
		return &Flag{Type: Alignment}, nil
	case "--opacity", "-o":
		return &Flag{Type: Opacity}, nil
	case "--stretch", "-s":
		return &Flag{Type: Stretch}, nil
	case "--random", "-r":
		return &Flag{Type: Random}, nil
	default:
		return nil, fmt.Errorf("unknown flag: %s", s)
	}
}

func (f *Flag) IsNone() bool {
	return f.Type == None
}

func (f *Flag) ValidateValue(val string) error {
	switch f.Type {
	case Profile:
		return ValidateProfile(val)
	case Interval:
		return ValidateInterval(val)
	case Alignment:
		return ValidateAlignment(val)
	case Opacity:
		return ValidateOpacity(val)
	case Stretch:
		return ValidateStretch(val)
	case Random:
		return ValidateRandom(val)
	case None:
		return nil
	default:
		return fmt.Errorf("unexpected error: unknown flag: %d", f.Type)
	}
}

// putting flags value validation in here for now
func ValidateProfile(val string) error {
	if val == "default" {
		return nil
	}
	list, num, isList := strings.Cut(val, "-")
	if list != "list" {
		return fmt.Errorf("invalid arg '%s' for --profile: must be 'list' followed by a dash then number", val)
	}
	if !isList {
		return fmt.Errorf("invalid arg '%s' for --profile: list and number must be separated by '-'", val)
	}
	_, err := strconv.Atoi(num)
	if err != nil {
		return fmt.Errorf("invalid arg '%s' for --profile: %s", val, err.Error())
	}
	return nil
}

func ValidateInterval(val string) error {
	num, err := strconv.Atoi(val)
	if err != nil {
		return fmt.Errorf("invalid arg '%s' for --interval: %s", val, err.Error())
	}
	if num < 1 {
		return fmt.Errorf("invalid arg '%s' for --interval: must be greater than 0", val)
	}
	return nil
}

func ValidateAlignment(val string) error {
	switch val {
	case "topLeft", "top", "topRight", "left", "center", "right", "bottomLeft", "bottom", "bottomRight":
		return nil
	default:
		return fmt.Errorf("invalid arg '%s' for --alignment: unknown alignment", val)
	}
}

func ValidateOpacity(val string) error {
	num, err := strconv.ParseFloat(val, 64)
	if err != nil {
		return fmt.Errorf("invalid arg '%s' for --opacity: %s", val, err.Error())
	}
	if num < 0 || num > 1 {
		return fmt.Errorf("invalid arg '%s' for --opacity: must be between 0 and 1", val)
	}
	return nil
}

func ValidateStretch(val string) error {
	switch val {
	case "none", "fill", "uniform", "uniformToFill":
		return nil
	default:
		return fmt.Errorf("invalid arg '%s' for --stretch: unknown stretch", val)
	}
}

func ValidateRandom(val string) error {
	switch val {
	case "":
		return nil
	default:
		return fmt.Errorf("'--random' flag does not take any arguments. got '%s'", val)
	}
}
