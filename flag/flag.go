package flag

import (
	"fmt"
	"strconv"
	"strings"
)

type FlagType uint8

const (
	Profile FlagType = iota
	Interval
	Alignment
	Opacity
	Stretch
)

type Flag struct {
	Type  FlagType
	Value string
}

func ToFlag(s string) (*Flag, error) {
	switch s {
	case "profile":
		return &Flag{Type: Profile}, nil
	case "interval":
		return &Flag{Type: Interval}, nil
	case "alignment":
		return &Flag{Type: Alignment}, nil
	case "opacity":
		return &Flag{Type: Opacity}, nil
	case "stretch":
		return &Flag{Type: Stretch}, nil
	default:
		return nil, fmt.Errorf("unknown flag: %s", s)
	}
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
	_, err := strconv.Atoi(val)
	if err != nil {
		return fmt.Errorf("invalid arg '%s' for --interval: %s", val, err.Error())
	}
	return nil
}

func ValidateAlignment(val string) error {
	switch val {
	case "top-right", "tr", "top-left", "tl", "top", "t", "left", "l", "center", "c", "right", "r", "bottom-right", "br", "bottom-left", "bl", "bottom", "b":
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
	case "none", "fill", "uniform", "uniform-fill":
		return nil
	default:
		return fmt.Errorf("invalid arg '%s' for --stretch: unknown stretch", val)
	}
}