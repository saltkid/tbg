package main

import (
	"fmt"
	"strconv"
	"strings"
)

func ValidateProfile(val *string) (*string, error) {
	if val == nil {
		return nil, fmt.Errorf("--profile must have an argument. got none")
	}
	if *val == "default" {
		return val, nil
	}
	list, num, isList := strings.Cut(*val, "-")
	if list != "list" {
		return nil, fmt.Errorf("invalid arg '%s' for --profile: must be 'list' followed by a dash then number", *val)
	}
	if !isList {
		return nil, fmt.Errorf("invalid arg '%s' for --profile: list and number must be separated by '-'", *val)
	}
	_, err := strconv.Atoi(num)
	if err != nil {
		return nil, fmt.Errorf("invalid arg '%s' for --profile: %s", *val, err.Error())
	}
	return val, nil
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

func ValidateRandom(val *string) (*bool, error) {
	ret := true
	if val == nil {
		return &ret, nil
	}
	switch *val {
	case "":
		return &ret, nil
	default:
		return nil, fmt.Errorf("'--random' flag does not take any arguments. got '%s'", *val)
	}
}
