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
	run           func(*Command) error
}

func (c *Command) ValidateValue(s string) error {
	return c.validateValue(s)
}

// define commands here
var CLI_CMDS = []*Command{
	RUN_CMD,
	ADD_CMD,
	CONFIG_CMD,
}

func ToCommand(s string) (*Command, error) {
	switch s {
	case RUN_CMD.name:
		return RUN_CMD, nil
	case ADD_CMD.name:
		return ADD_CMD, nil
	case CONFIG_CMD.name:
		return CONFIG_CMD, nil
	default:
		return &Command{}, fmt.Errorf("'%s' is not a valid command", s)
	}
}

var (
	RUN_CMD = &Command{
		name:  "run",
		value: "",
		validateValue: func(s string) error {
			switch s {
			case "":
				return nil
			default:
				return fmt.Errorf("run takes no args. got '%s'", s)
			}
		},
		validateFlag: func(flagName string, flagValue string) error {
			switch flagName {
			case TARGET_FLAG.name, TARGET_FLAG.short:
				return TARGET_FLAG.validateValue(flagValue)
			case INTERVAL_FLAG.name, INTERVAL_FLAG.short:
				return INTERVAL_FLAG.validateValue(flagValue)
			case CONFIG_CMD.name:
				return CONFIG_CMD.validateValue(flagValue)
			default:
				return fmt.Errorf("invalid flag for 'run': '%s'", flagName)
			}

		},
	}

	ADD_CMD = &Command{
		name:  "add",
		value: "",
		validateValue: func(path string) error {
			if path == "" {
				return fmt.Errorf("no path provided for 'add'")
			}

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

	CONFIG_CMD = &Command{
		name:  "config",
		value: "",
		validateValue: func(s string) error {
			// default config
			if s == "" || s == "default" {
				configPath, err := filepath.Abs("config.yaml")
				if err != nil {
					return err
				}
				// check if exists first, then create if needed
				if _, err := os.Stat(configPath); os.IsNotExist(err) {
					err := DefaultTemplate(configPath).WriteFile()
					if err != nil {
						return fmt.Errorf("error creating default config: %s", err.Error())
					}
				}

			} else {
				// user config
				configPath, err := filepath.Abs(s)
				if err != nil {
					return err
				}
				if _, err := os.Stat(s); os.IsNotExist(err) {
					// create parent dirs if needed first before writing to file
					err := os.MkdirAll(filepath.Dir(configPath), os.ModePerm)
					if err != nil {
						return fmt.Errorf("error creating parent dirs of user config %s: %s", configPath, err.Error())
					}

					err = UserTemplate(configPath).WriteFile()
					if err != nil {
						return fmt.Errorf("error creating default config: %s", err.Error())
					}
				}
			}

			return nil
		},
		validateFlag: func(flagName string, flagValue string) error {
			switch flagName {
			case CREATE_FLAG.name, CREATE_FLAG.short:
				return CREATE_FLAG.validateValue(flagValue)
			default:
				return fmt.Errorf("invalid flag '%s' for 'config'", flagName)
			}
		},
		run: func(cmd *Command) error {
			var configPath string
			if cmd.value == "" {
				// print used config
				configPath, _ = filepath.Abs("config.yaml")
				yamlFile, err := os.ReadFile(configPath)
				if err != nil {
					return err
				}

				contents := DefaultConfig{}
				err = yaml.Unmarshal(yamlFile, &contents)
				if err != nil {
					return fmt.Errorf("error reading config: %s", err.Error())
				}

				// print user config if set
				if contents.UseUserConfig {
					userConfigPath, err := filepath.Abs(contents.UserConfig)
					if err != nil {
						return err
					}

					yamlFile, err = os.ReadFile(userConfigPath)
					if err != nil {
						return err
					}

					contents := UserConfig{}
					err = yaml.Unmarshal(yamlFile, &contents)
					if err != nil {
						return fmt.Errorf("error reading config: %s", err.Error())
					}
					contents.Log(userConfigPath)

				} else {
					contents.Log(configPath)
				}

			} else if cmd.value == "default" {
				// read default config
				configPath, _ = filepath.Abs("config.yaml")
				yamlFile, err := os.ReadFile(configPath)
				if err != nil {
					return err
				}

				contents := DefaultConfig{}
				err = yaml.Unmarshal(yamlFile, &contents)
				if err != nil {
					return fmt.Errorf("error reading config: %s", err.Error())
				}
				contents.UseUserConfig = false

				// log to console
				contents.Log(configPath)

				// update default config to use default config
				template := DefaultTemplate(configPath)
				template.yamlContents, _ = yaml.Marshal(contents)
				err = template.WriteFile()
				if err != nil {
					return fmt.Errorf("error writing config: %s", err.Error())
				}

			} else {
				// read config
				configPath, _ = filepath.Abs(cmd.value)
				yamlFile, err := os.ReadFile(configPath)
				if err != nil {
					return err
				}

				contents := UserConfig{}
				err = yaml.Unmarshal(yamlFile, &contents)
				if err != nil {
					return fmt.Errorf("error reading config: %s", err.Error())
				}

				// log to console
				contents.Log(configPath)

				// edit default config to use user config
				defaultContents := DefaultConfig{
					UseUserConfig: true,
					UserConfig:    configPath,
					ImageColPaths: contents.ImageColPaths,
					Interval:      contents.Interval,
					Target:        contents.Target,
					Alignment:     contents.Alignment,
					Opacity:       contents.Opacity,
					Stretch:       contents.Stretch,
				}
				defaultConfigPath, _ := filepath.Abs("config.yaml")
				template := DefaultTemplate(defaultConfigPath)
				template.yamlContents, _ = yaml.Marshal(defaultContents)
				err = template.WriteFile()
				if err != nil {
					return fmt.Errorf("error writing config: %s", err.Error())
				}
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

func (f *Flag) ValidateValue(s string) error {
	return f.validateValue(s)
}

// define flags here
var CLI_FLAGS = []*Flag{
	ALIGN_FLAG,
	OPACITY_FLAG,
	STRETCH_FLAG,
	TARGET_FLAG,
	INTERVAL_FLAG,
	CREATE_FLAG,
}

func ToFlag(s string) (*Flag, error) {
	switch s {
	case ALIGN_FLAG.name, ALIGN_FLAG.short:
		return ALIGN_FLAG, nil
	case OPACITY_FLAG.name, OPACITY_FLAG.short:
		return OPACITY_FLAG, nil
	case STRETCH_FLAG.name, STRETCH_FLAG.short:
		return STRETCH_FLAG, nil
	case TARGET_FLAG.name, TARGET_FLAG.short:
		return TARGET_FLAG, nil
	case INTERVAL_FLAG.name, INTERVAL_FLAG.short:
		return INTERVAL_FLAG, nil
	case CREATE_FLAG.name, CREATE_FLAG.short:
		return CREATE_FLAG, nil
	default:
		return &Flag{}, fmt.Errorf("'%s' is not a valid flag", s)
	}
}

var (
	ALIGN_FLAG = &Flag{
		name:  "--alignment",
		short: "-a",
		value: "center",
		validateValue: func(s string) error {
			switch s {
			case "top", "t", "top-right", "tr", "top-left", "tl", "center", "left", "right", "bottom", "b", "bottom-right", "br", "bottom-left", "bl":
				return nil
			case "":
				return fmt.Errorf("missing value for --alignment")
			default:
				return fmt.Errorf("invalid value for --alignment: '%s'", s)
			}
		},
	}

	OPACITY_FLAG = &Flag{
		name:  "--opacity",
		short: "-o",
		value: "0.1",
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

	STRETCH_FLAG = &Flag{
		name:  "--stretch",
		short: "-s",
		value: "uniform",
		validateValue: func(s string) error {
			switch s {
			case "fill", "none", "uniform", "uniform-fill":
				return nil
			case "":
				return fmt.Errorf("missing value for --stretch")
			default:
				return fmt.Errorf("invalid value for --stretch: '%s'", s)
			}
		},
	}

	TARGET_FLAG = &Flag{
		name:  "--target",
		short: "-t",
		value: "default",
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
				return fmt.Errorf("invalid value for --target: '%s'", s)
			}
		},
	}

	INTERVAL_FLAG = &Flag{
		name:  "--interval",
		short: "-i",
		value: "30",
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

	CREATE_FLAG = &Flag{
		name:  "--create",
		short: "-c",
		value: "",
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

func IsValidCommandName(s string) bool {
	validCommands := make(map[string]struct{})
	for _, cmd := range CLI_CMDS {
		validCommands[cmd.name] = struct{}{}
	}

	_, exists := validCommands[s]
	return exists
}

func IsValidFlagName(s string) bool {
	validFlags := make(map[string]struct{})
	for _, flag := range CLI_FLAGS {
		validFlags[flag.name] = struct{}{}
	}

	_, exists := validFlags[s]
	return exists
}
