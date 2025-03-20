package main

import (
	"fmt"
	"strconv"
)

func ValidateAlignment(val *string) (*string, error) {
	if val == nil {
		return nil, fmt.Errorf("--interval must have an argument. got none")
	}
	switch *val {
	case "topLeft", "top", "topRight", "left", "center", "right", "bottomLeft", "bottom", "bottomRight":
		return val, nil
	default:
		return nil, fmt.Errorf("invalid arg '%s' for --alignment: unknown alignment", *val)
	}
}

func ValidateInterval(val *string) (*uint16, error) {
	if val == nil {
		return nil, fmt.Errorf("--interval must have an argument. got none")
	}
	num, err := strconv.Atoi(*val)
	if err != nil {
		return nil, fmt.Errorf("invalid arg '%s' for --interval: %s", *val, err.Error())
	}
	if num < 1 {
		return nil, fmt.Errorf("invalid arg '%s' for --interval: must be greater than 0", *val)
	}
	ret := uint16(num)
	return &ret, nil
}

func ValidateOpacity(val *string) (*float32, error) {
	if val == nil {
		return nil, fmt.Errorf("--interval must have an argument. got none")
	}
	num, err := strconv.ParseFloat(*val, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid arg '%s' for --opacity: %s", *val, err.Error())
	}
	if num < 0 || num > 1 {
		return nil, fmt.Errorf("invalid arg '%s' for --opacity: must be between 0 and 1", *val)
	}
	ret := float32(num)
	return &ret, nil
}

func ValidateProfile(val *string) (*string, error) {
	if val == nil {
		return nil, fmt.Errorf("--profile must have an argument. got none")
	}
	valNum, err := strconv.Atoi(*val)
	if err == nil && valNum < 1 {
		return nil, fmt.Errorf("invalid arg '%d' for --profile: profile indices start at 1.", valNum)
	}
	return val, nil
}

func ValidateStretch(val *string) (*string, error) {
	if val == nil {
		return nil, fmt.Errorf("--stretch must have an argument. got none")
	}
	switch *val {
	case "none", "fill", "uniform", "uniformToFill":
		return val, nil
	default:
		return nil, fmt.Errorf("invalid arg '%s' for --stretch: unknown stretch", *val)
	}
}
