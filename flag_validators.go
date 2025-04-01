package main

import (
	"fmt"
	"os"
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
		return nil, fmt.Errorf(`invalid arg '%s' for --alignment: unknown alignment
[topLeft top topRight left center right bottomLeft bottom bottomRight]`, *val)
	}
}

// validates only that the file exists
func ValidateConfig(val *string) (*string, error) {
	if val == nil {
		return nil, fmt.Errorf("--interval must have an argument. got none")
	}
	absPath, err := NormalizePath(*val)
	if err != nil {
		return nil, fmt.Errorf("Failed to normalize path %s: %s", *val, err)
	}
	if _, err := os.Stat(absPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("%s does not exist: %s", *val, err.Error())
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

func ValidatePort(val *string) (*uint16, error) {
	if val == nil {
		return nil, fmt.Errorf("--port must have an argument. got none")
	}
	valNum, err := strconv.Atoi(*val)
	if err != nil {
		return nil, fmt.Errorf("invalid arg '%s' for --port: %s", *val, err)
	}
	if valNum < 1 {
		return nil, fmt.Errorf("invalid arg '%d' for --port: must be positive.", valNum)
	}
	if valNum > 65535 {
		return nil, fmt.Errorf("invalid arg '%d' for --port: ports can only be up to 65535.", valNum)
	}
	ret := uint16(valNum)
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
		return nil, fmt.Errorf(`invalid arg '%s' for --stretch: unknown stretch
[fill uniform uniformToFill none]`, *val)
	}
}
