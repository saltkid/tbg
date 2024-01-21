package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

// common interface to allow slice of either Command/Flag in flags
// since both Command and Flag can be a valid flag for a Command
type CLI_Arg interface {
	ValidateValue(string) error
}

type Command struct {
	name          string
	value         string
	validateValue func(string) error
	flags         []CLI_Arg
	validateFlag  func(string, string) error
}

func (c Command) ValidateValue(s string) error {
	return c.validateValue(s)
}

// define commands here
var CLI_CMDS = []Command{
	RUN_CMD,
	ADD_CMD,
	CONFIG_CMD,
}
var (
	RUN_CMD = Command{
		name: "run",
		validateValue: func(s string) error {
			switch s {
			case "":
				return nil
			default:
				return fmt.Errorf("run takes no args. got '%s'", s)
			}
		},
		validateFlag: func(flagName string, flagValue string) error {
			if flagValue == "" {
				return fmt.Errorf("missing argument for flag '%s'", flagName)
			}

			switch flagName {
			case TARGET_FLAG.name, TARGET_FLAG.short:
				return TARGET_FLAG.validateValue(flagValue)
			case INTERVAL_FLAG.name, INTERVAL_FLAG.short:
				return INTERVAL_FLAG.validateValue(flagValue)
			case CONFIG_CMD.name:
				return CONFIG_CMD.validateValue(flagValue)
			default:
				return fmt.Errorf("invalid flag for 'add': '%s'", flagName)
			}

		},
	}

	ADD_CMD = Command{
		name: "add",
		validateValue: func(path string) error {
			// check if exists
			if _, err := os.Stat(path); os.IsNotExist(err) {
				return fmt.Errorf("%s does not exist: %s", path, err.Error())
			}

			// check if has any image files
			imgFileCount := 0
			err := filepath.WalkDir(path, func(p string, d os.DirEntry, err error) error {
				if err != nil {
					return err
				}

				if d.IsDir() && d.Name() != filepath.Base(path) {
					return filepath.SkipDir
				}

				if !d.IsDir() && IsImageFile(d.Name()) {
					imgFileCount++
				}

				return nil
			})

			if err != nil {
				return fmt.Errorf("error reading %s: %s", path, err.Error())
			}

			if imgFileCount < 1 {
				return fmt.Errorf("no image files in %s", path)
			}

			return nil
		},
		validateFlag: func(flagName string, flagValue string) error {
			if flagValue == "" {
				return fmt.Errorf("missing argument for flag '%s'", flagName)
			}

			switch flagName {
			case ALIGN_FLAG.name, ALIGN_FLAG.short:
				return ALIGN_FLAG.validateValue(flagValue)
			case OPACITY_FLAG.name, OPACITY_FLAG.short:
				return OPACITY_FLAG.validateValue(flagValue)
			case STRETCH_FLAG.name, STRETCH_FLAG.short:
				return STRETCH_FLAG.validateValue(flagValue)
			default:
				return fmt.Errorf("invalid flag for 'add': '%s'", flagName)
			}
		},
	}

	CONFIG_CMD = Command{
		name: "config",
		validateValue: func(s string) error {
			if s == "" || s == "default" {
				return nil
			}

			// check if exists
			if _, err := os.Stat(s); os.IsNotExist(err) {
				return fmt.Errorf("%s does not exist: %s", s, err.Error())
			}

			// check if has a config file (.yaml)
			configCount := 0
			err := filepath.WalkDir(s, func(p string, d os.DirEntry, err error) error {
				if err != nil {
					return err
				}

				if d.IsDir() && d.Name() != filepath.Base(s) {
					return filepath.SkipDir
				}

				if !d.IsDir() && filepath.Ext(p) == ".yaml" {
					yamlFile, err := os.ReadFile(s)
					if err != nil {
						return err
					}

					contents := Config{}
					err = yaml.Unmarshal(yamlFile, &contents)
					if err != nil {
						return err
					}
					configCount++
				}

				return nil
			})

			if err != nil {
				return fmt.Errorf("error reading %s: %s", s, err.Error())
			}

			if configCount == 0 {
				return fmt.Errorf("no config files found in %s", s)

			} else if configCount > 1 {
				return fmt.Errorf("multiple config files found in %s", s)
			}

			return nil
		},
		validateFlag: func(flagName string, flagValue string) error {
			switch flagName {
			case CREATE_FLAG.name, CREATE_FLAG.short:
				return CREATE_FLAG.validateValue(flagValue)
			default:
				return fmt.Errorf("invalid flag '%s' for config", flagName)
			}
		},
	}
)

type Flag struct {
	name          string
	short         string
	value         string
	validateValue func(string) error
}

func (f Flag) ValidateValue(s string) error {
	return f.validateValue(s)
}

// define flags here
var CLI_FLAGS = []Flag{
	ALIGN_FLAG,
	OPACITY_FLAG,
	STRETCH_FLAG,
	TARGET_FLAG,
	INTERVAL_FLAG,
}
var (
	ALIGN_FLAG = Flag{
		name:  "--alignment",
		short: "-a",
		validateValue: func(s string) error {
			switch s {
			case "top", "t", "top-right", "tr", "top-left", "tl", "center", "left", "right", "bottom", "b", "bottom-right", "br", "bottom-left", "bl":
				return nil
			default:
				return fmt.Errorf("invalid value for --alignment: %s", s)
			}
		},
	}

	OPACITY_FLAG = Flag{
		name:  "--opacity",
		short: "-o",
		validateValue: func(s string) error {
			num, err := strconv.Atoi(s)
			if err != nil {
				return err
			}

			if num > 1 || num < 0 {
				return fmt.Errorf("invalid value for --opacity: %d; must a float between 0-1", num)
			}
			return nil
		},
	}

	STRETCH_FLAG = Flag{
		name:  "--stretch",
		short: "-s",
		validateValue: func(s string) error {
			switch s {
			case "fill", "none", "uniform", "uniform-fill":
				return nil
			default:
				return fmt.Errorf("invalid value for --stretch: %s", s)
			}
		},
	}

	TARGET_FLAG = Flag{
		name:  "--target",
		short: "-t",
		validateValue: func(s string) error {
			if s == "default" {
				return nil

			} else if strings.HasPrefix(s, "list-") {
				// check if list- is followed by number
				numPart, _ := strings.CutPrefix(s, "list-")
				if numPart == "" {
					return fmt.Errorf("no number found after 'list-' for --target")
				}

				_, err := strconv.Atoi(numPart)
				if err != nil {
					return fmt.Errorf("invalid number '%s' after 'list-' for --target; error: %s", numPart, err.Error())
				}
				return nil

			} else {
				return fmt.Errorf("invalid value '%s' for --target", s)
			}
		},
	}

	INTERVAL_FLAG = Flag{
		name:  "--interval",
		short: "-i",
		validateValue: func(f string) error {
			if f == "" {
				return fmt.Errorf("missing value for --interval")
			}

			_, err := strconv.ParseFloat(f, 64)
			if err != nil {
				return fmt.Errorf("invalid float value '%s' for --interval; error: %s", f, err.Error())
			}
			return nil
		},
	}

	CREATE_FLAG = Flag{
		name:  "--create",
		short: "-c",
		validateValue: func(s string) error {
			switch s {
			case "":
				return nil
			default:
				return fmt.Errorf("--create takes no args. got '%s'", s)
			}
		},
	}
)

func IsValidArgName(a string) bool {
	validNames := make(map[string]struct{})

	for _, cmd := range CLI_CMDS {
		validNames[cmd.name] = struct{}{}
	}
	for _, flag := range CLI_FLAGS {
		validNames[flag.name] = struct{}{}
	}

	_, exists := validNames[a]
	return exists
}
