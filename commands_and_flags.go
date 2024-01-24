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
	flags         map[string]CLI_Arg
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
	REMOVE_CMD,
	CONFIG_CMD,
}

func ToCommand(s string) (*Command, error) {
	switch s {
	case RUN_CMD.name:
		return RUN_CMD, nil
	case ADD_CMD.name:
		return ADD_CMD, nil
	case REMOVE_CMD.name:
		return REMOVE_CMD, nil
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
		flags: make(map[string]CLI_Arg),
		validateFlag: func(flagName string, flagValue string) error {
			switch flagName {
			case PROFILE_FLAG.name, PROFILE_FLAG.short:
				return PROFILE_FLAG.validateValue(flagValue)
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
			absPath, err := filepath.Abs(path)
			if err != nil {
				return err
			}
			if _, err := os.Stat(absPath); os.IsNotExist(err) {
				return fmt.Errorf("%s does not exist: %s", absPath, err.Error())
			}

			// check if has any image files
			imgFileCount := 0
			err = filepath.WalkDir(absPath, func(p string, d os.DirEntry, err error) error {
				if err != nil {
					return err
				}

				if d.IsDir() && d.Name() != filepath.Base(absPath) {
					return filepath.SkipDir
				}

				if !d.IsDir() && IsImageFile(d.Name()) {
					imgFileCount++
				}

				return nil
			})

			if err != nil {
				return fmt.Errorf("error reading %s: %s", absPath, err.Error())
			}

			if imgFileCount < 1 {
				return fmt.Errorf("no image files in %s", absPath)
			}

			return nil
		},
		flags: make(map[string]CLI_Arg),
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
		run: func(c *Command) error {
			absPath, _ := filepath.Abs(strings.ToLower(c.value))

			// get values of flags if set
			alignment, ok := c.flags[ALIGN_FLAG.name].(*Flag)
			var alignVal string
			if ok {
				alignVal = alignment.value
			}
			opacity, ok := c.flags[OPACITY_FLAG.name].(*Flag)
			var opacityVal string
			if ok {
				opacityVal = opacity.value
			}
			stretch, ok := c.flags[STRETCH_FLAG.name].(*Flag)
			var stretchVal string
			if ok {
				stretchVal = stretch.value
			}

			// only set flags after path if at least one is set
			if alignVal != "" || opacityVal != "" || stretchVal != "" {
				if stretchVal == "" {
					stretchVal = "uniform"
				}
				if alignVal == "" {
					alignVal = "center"
				}
				if opacityVal == "" {
					opacityVal = "0.1"
				}
				absPath = fmt.Sprintf("%s | %s %s %s", absPath, alignVal, stretchVal, opacityVal)
			}

			// read config to determine whether to append these values to default config or user config
			err := CONFIG_CMD.validateValue("default")
			if err != nil {
				return err
			}
			configPath, _ := filepath.Abs("config.yaml")
			yamlFile, err := os.ReadFile(configPath)
			if err != nil {
				return err
			}
			contents := DefaultConfig{}
			err = yaml.Unmarshal(yamlFile, &contents)
			if err != nil {
				return fmt.Errorf("error reading config: %s", err.Error())
			}

			if contents.UseUserConfig {
				err := CONFIG_CMD.validateValue(contents.UserConfig)
				if err != nil {
					return err
				}
				userConfigPath, _ := filepath.Abs(contents.UserConfig)
				yamlFile, err = os.ReadFile(userConfigPath)
				if err != nil {
					return err
				}
				contents := UserConfig{}
				err = yaml.Unmarshal(yamlFile, &contents)
				if err != nil {
					return fmt.Errorf("error reading config: %s", err.Error())
				}

				// check if exists before appending
				for _, path := range contents.ImageColPaths {
					if path == absPath {
						return fmt.Errorf("%s already exists in user config", absPath)
					}
				}
				contents.ImageColPaths = append(contents.ImageColPaths, absPath)

				template := UserTemplate(userConfigPath)
				template.yamlContents, _ = yaml.Marshal(contents)
				err = template.WriteFile()
				if err != nil {
					return fmt.Errorf("error writing config: %s", err.Error())
				}
				contents.Log(userConfigPath)

			} else {
				for _, path := range contents.ImageColPaths {
					if path == absPath {
						return fmt.Errorf("%s already exists in default config", absPath)
					}
				}
				contents.ImageColPaths = append(contents.ImageColPaths, absPath)

				template := DefaultTemplate(configPath)
				template.yamlContents, _ = yaml.Marshal(contents)
				err = template.WriteFile()
				if err != nil {
					return fmt.Errorf("error writing config: %s", err.Error())
				}

				contents.Log(configPath)
			}

			return nil
		},
	}

	REMOVE_CMD = &Command{
		name:  "remove",
		value: "",
		validateValue: func(path string) error {
			if path == "" {
				return fmt.Errorf("no path provided for 'remove'")
			}

			return nil
		},
		flags: make(map[string]CLI_Arg),
		validateFlag: func(flagName string, flagValue string) error {
			switch flagName {
			case "":
				return CONFIG_CMD.validateValue("default")
			case CONFIG_CMD.name:
				if flagValue == "" {
					return fmt.Errorf("no config provided for 'remove'. please remove empty 'config' flag to use currently set config")
				}
				return CONFIG_CMD.validateValue(flagValue)
			default:
				return fmt.Errorf("invalid flag for 'remove': '%s'", flagName)
			}
		},
		run: func(c *Command) error {
			absPath, _ := filepath.Abs(strings.ToLower(c.value))

			// read config to determine whether to delete values on default config or user config
			var configPath string
			var contents Config
			configFlag, hasConfigFlag := c.flags[CONFIG_CMD.name].(*Command)

			if hasConfigFlag {
				if configFlag.value == "default" {
					configPath, _ = filepath.Abs("config.yaml")
					contents = &DefaultConfig{}
				} else {
					configPath, _ = filepath.Abs(configFlag.value)
					contents = &UserConfig{}
				}
			} else {
				configPath, _ = filepath.Abs("config.yaml")
				contents = &DefaultConfig{}
			}

			// read config file
			yamlFile, err := os.ReadFile(configPath)
			if err != nil {
				return err
			}

			userContents, IsUserConfig := contents.(*UserConfig)
			defaultContents, IsDefaultConfig := contents.(*DefaultConfig)
			if IsUserConfig {
				// remove in user config
				err = yaml.Unmarshal(yamlFile, &userContents)
				if err != nil {
					return fmt.Errorf("error reading config: %s", err.Error())
				}

				// delete matched path
				removedDir := ""
				for i, path := range userContents.ImageColPaths {
					path = strings.TrimSpace(strings.Split(path, "|")[0])
					if strings.EqualFold(path, absPath) {
						removedDir = path
						userContents.ImageColPaths = append(userContents.ImageColPaths[:i], userContents.ImageColPaths[i+1:]...)
						break
					}
				}
				if removedDir == "" {
					return fmt.Errorf("%s not found in user config %s", absPath, configPath)
				}
				removedDir = fmt.Sprintf(`"%s" in user config %s`, removedDir, configPath)

				template := UserTemplate(configPath)
				template.yamlContents, _ = yaml.Marshal(contents)
				err = template.WriteFile()
				if err != nil {
					return fmt.Errorf("error writing config: %s", err.Error())
				}
				contents.Log(configPath).LogRemoved(removedDir)

			} else if IsDefaultConfig && hasConfigFlag {
				// remove in default config
				err = yaml.Unmarshal(yamlFile, &defaultContents)
				if err != nil {
					return fmt.Errorf("error reading config: %s", err.Error())
				}
				// delete matched path
				removedDir := ""
				for i, path := range defaultContents.ImageColPaths {
					path = strings.TrimSpace(strings.Split(path, "|")[0])
					println(path)
					if strings.EqualFold(path, absPath) {
						removedDir = path
						defaultContents.ImageColPaths = append(defaultContents.ImageColPaths[:i], defaultContents.ImageColPaths[i+1:]...)
						break
					}
				}
				if removedDir == "" {
					return fmt.Errorf("%s not found in default config", c.value)
				}
				removedDir = fmt.Sprintf(`"%s" in default config`, removedDir)

				template := DefaultTemplate(configPath)
				template.yamlContents, _ = yaml.Marshal(defaultContents)
				err = template.WriteFile()
				if err != nil {
					return fmt.Errorf("error writing config: %s", err.Error())
				}

				contents.Log(configPath).LogRemoved(removedDir)

			} else {
				// removed in currently used config (can be default or user)
				err = yaml.Unmarshal(yamlFile, &defaultContents)
				if err != nil {
					return fmt.Errorf("error reading config: %s", err.Error())
				}

				if defaultContents.UseUserConfig {
					err := CONFIG_CMD.validateValue(defaultContents.UserConfig)
					if err != nil {
						return err
					}
					userConfigPath, _ := filepath.Abs(defaultContents.UserConfig)
					yamlFile, err = os.ReadFile(userConfigPath)
					if err != nil {
						return err
					}
					contents := UserConfig{}
					err = yaml.Unmarshal(yamlFile, &contents)
					if err != nil {
						return fmt.Errorf("error reading config: %s", err.Error())
					}

					// delete matched path
					removedDir := ""
					for i, path := range contents.ImageColPaths {
						path = strings.TrimSpace(strings.Split(path, "|")[0])
						if strings.EqualFold(path, absPath) {
							removedDir = path
							contents.ImageColPaths = append(contents.ImageColPaths[:i], contents.ImageColPaths[i+1:]...)
							break
						}
					}
					if removedDir == "" {
						return fmt.Errorf("%s not found in user config", c.value)
					}

					template := UserTemplate(userConfigPath)
					template.yamlContents, _ = yaml.Marshal(contents)
					err = template.WriteFile()
					if err != nil {
						return fmt.Errorf("error writing config: %s", err.Error())
					}
					contents.Log(userConfigPath).LogRemoved(removedDir)

				} else {
					// delete matched path
					removedDir := ""
					for i, path := range defaultContents.ImageColPaths {
						path = strings.TrimSpace(strings.Split(path, "|")[0])
						if strings.EqualFold(path, absPath) {
							removedDir = path
							defaultContents.ImageColPaths = append(defaultContents.ImageColPaths[:i], defaultContents.ImageColPaths[i+1:]...)
							break
						}
					}
					if removedDir == "" {
						return fmt.Errorf("%s not found in default config", c.value)
					}

					template := DefaultTemplate(configPath)
					template.yamlContents, _ = yaml.Marshal(contents)
					err = template.WriteFile()
					if err != nil {
						return fmt.Errorf("error writing config: %s", err.Error())
					}

					contents.Log(configPath).LogRemoved(removedDir)
				}
			}

			return nil
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
		flags: make(map[string]CLI_Arg),
		validateFlag: func(flagName string, flagValue string) error {
			switch flagName {
			case "":
				return nil
			default:
				return fmt.Errorf("'config' has no flags. got '%s'", flagName)
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

				// update default config to use default config
				template := DefaultTemplate(configPath)
				template.yamlContents, _ = yaml.Marshal(contents)
				err = template.WriteFile()
				if err != nil {
					return fmt.Errorf("error writing config: %s", err.Error())
				}
				// log to console
				contents.Log(configPath)

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

				// edit default config to use user config
				defaultConfigPath, _ := filepath.Abs("config.yaml")
				yamlFile, err = os.ReadFile(defaultConfigPath)
				if err != nil {
					return err
				}

				defaultContents := DefaultConfig{}
				err = yaml.Unmarshal(yamlFile, &defaultContents)
				if err != nil {
					return fmt.Errorf("error reading default config to update it: %s", err.Error())
				}

				defaultContents.UseUserConfig = true
				defaultContents.UserConfig = configPath

				template := DefaultTemplate(defaultConfigPath)
				template.yamlContents, _ = yaml.Marshal(defaultContents)
				err = template.WriteFile()
				if err != nil {
					return fmt.Errorf("error writing config: %s", err.Error())
				}

				// log to console
				contents.Log(configPath)
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
	PROFILE_FLAG,
	INTERVAL_FLAG,
}

func ToFlag(s string) (*Flag, error) {
	switch s {
	case ALIGN_FLAG.name, ALIGN_FLAG.short:
		return ALIGN_FLAG, nil
	case OPACITY_FLAG.name, OPACITY_FLAG.short:
		return OPACITY_FLAG, nil
	case STRETCH_FLAG.name, STRETCH_FLAG.short:
		return STRETCH_FLAG, nil
	case PROFILE_FLAG.name, PROFILE_FLAG.short:
		return PROFILE_FLAG, nil
	case INTERVAL_FLAG.name, INTERVAL_FLAG.short:
		return INTERVAL_FLAG, nil
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
			num, err := strconv.ParseFloat(s, 64)
			if err != nil {
				return err
			}

			if num > 1 || num < 0 {
				return fmt.Errorf("invalid value for --opacity: %v; must a float between 0-1", num)
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

	PROFILE_FLAG = &Flag{
		name:  "--profile",
		short: "-p",
		value: "default",
		validateValue: func(s string) error {
			if s == "default" {
				return nil

			} else if strings.HasPrefix(s, "list-") {
				// check if list- is followed by number
				numPart, _ := strings.CutPrefix(s, "list-")
				if numPart == "" {
					return fmt.Errorf("no number found after 'list-' for --profile")
				}

				_, err := strconv.Atoi(numPart)
				if err != nil {
					return fmt.Errorf("invalid number '%s' after 'list-' for --profile; error: %s", numPart, err.Error())
				}
				return nil

			} else {
				return fmt.Errorf("invalid value for --profile: '%s'", s)
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
		validFlags[flag.short] = struct{}{}
	}

	_, exists := validFlags[s]
	return exists
}
