package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"gopkg.in/yaml.v3"
	// "gopkg.in/yaml.v3"
)

type Command struct {
	name          string
	values        []string
	validateValue func(string) error
	flags         []Flag
	validateFlag  func(string, string) error
}

// define commands here
var CLI_CMDS = []Command{
	ADD_CMD, CONFIG_CMD,
}
var (
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
	}
)

type Flag struct {
	name          string
	short         string
	value         string
	validateValue func(string) error
}

// define flags here
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
				return fmt.Errorf("invalid value for --opacity: %s; must a float between 0-1", num)
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
)

func IsValidArgName(a string) bool {
	for _, ARG := range CLI_CMDS {
		if a != ARG.name {
			return false
		}
	}
	return true
}
