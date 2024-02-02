package flag

import "fmt"

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
		return nil
	case Interval:
		return nil
	case Alignment:
		return nil
	case Opacity:
		return nil
	case Stretch:
		return nil
	default:
		return fmt.Errorf("unexpected error: unknown flag: %d", f.Type)
	}
}
