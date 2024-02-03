package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/saltkid/tbg/config"
	"github.com/saltkid/tbg/flag"
	"gopkg.in/yaml.v3"
)

func RemoveValidateValue(val string) error {
	absPath, err := filepath.Abs(val)
	if err != nil {
		return fmt.Errorf("Failed to get absolute path of %s: %s", val, err)
	}
	if d, err := os.Stat(absPath); os.IsNotExist(err) {
		return fmt.Errorf("%s does not exist: %s", val, err.Error())
	} else if !d.IsDir() {
		return fmt.Errorf("%s is not a directory", val)
	}
	return nil
}

func RemoveValidateFlag(t flag.FlagType) error {
	switch t {
	case flag.None:
		return nil
	default:
		return fmt.Errorf("'remove' takes no flags. got type: %d", t)
	}
}

func RemoveValidateSubCmd(t CmdType) error {
	switch t {
	case None:
		return nil
	default:
		return fmt.Errorf("'remove' takes no sub commands. got type: %d", t)
	}
}

func RemoveExecute(c *Cmd) error {
	absPath, _ := filepath.Abs(c.Value)

	// check if config subcommand is set by user
	specifiedConfig := ExtractSubCmdValue(Config, c.SubCmds)
	var configPath string
	var configContents config.Config
	if specifiedConfig == "default" || specifiedConfig == "" {
		configPath, _ = filepath.Abs(config.DefaultConfigPath())
		configContents = &config.DefaultConfig{}
	} else {
		configPath, _ = filepath.Abs(specifiedConfig)
		configContents = &config.UserConfig{}
	}

	yamlFile, err := os.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("Failed to read config: %s", err)
	}
	err = configContents.Unmarshal(yamlFile)
	if err != nil {
		return err
	}

	if configContents.IsDefaultConfig() && specifiedConfig == "default" {
		err = RemoveFromDefaultConfig(configContents, absPath, configPath)
		if err != nil {
			return err
		}
	} else if configContents.IsUserConfig() {
		err = RemoveFromUserConfig(configContents, absPath, configPath)
		if err != nil {
			return err
		}
	} else {
		// read default config if using user config or default
		defaultContents, _ := configContents.(*config.DefaultConfig)

		// using default config
		if defaultContents.UserConfig == "" {
			err = RemoveFromDefaultConfig(configContents, absPath, configPath)
			if err != nil {
				return err
			}
		} else {
			// recheck if user config path in default config is valid
			err = ConfigValidateValue(defaultContents.UserConfig)
			if err != nil {
				return err
			}

			userConfigPath, _ := filepath.Abs(defaultContents.UserConfig)
			yamlFile, err := os.ReadFile(userConfigPath)
			if err != nil {
				return fmt.Errorf("Failed to read user config: %s", err)
			}
			userContents := &config.UserConfig{}
			err = userContents.Unmarshal(yamlFile)
			if err != nil {
				return err
			}
			err = RemoveFromUserConfig(userContents, absPath, userConfigPath)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func RemoveFromUserConfig(contents config.Config, absPath string, configPath string) error {
	userContents, ok := contents.(*config.UserConfig)
	if !ok {
		return fmt.Errorf("unexpected error: contents is not a user config")
	}
	removed := make(map[string]struct{})
	for i, path := range userContents.ImageColPaths {
		purePath, _, _ := strings.Cut(path, "|")
		purePath = strings.TrimSpace(purePath)

		if strings.EqualFold(absPath, purePath) {
			removed[absPath] = struct{}{}
			userContents.ImageColPaths = append(userContents.ImageColPaths[:i], userContents.ImageColPaths[i+1:]...)
			break
		}
	}
	template := config.UserTemplate(configPath)
	template.YamlContents, _ = yaml.Marshal(userContents)
	err := template.WriteFile()
	if err != nil {
		return fmt.Errorf("error writing to user config: %s", err.Error())
	}
	if len(removed) == 0 {
		removed["no changes made"] = struct{}{}
	}
	userContents.Log(configPath).LogRemoved(removed)
	return nil
}

func RemoveFromDefaultConfig(contents config.Config, absPath string, configPath string) error {
	defaultContents, ok := contents.(*config.DefaultConfig)
	if !ok {
		return fmt.Errorf("unexpected error: contents is not a default config")
	}
	removed := make(map[string]struct{})
	for i, path := range defaultContents.ImageColPaths {
		purePath, _, _ := strings.Cut(path, "|")
		purePath = strings.TrimSpace(purePath)

		if strings.EqualFold(absPath, purePath) {
			removed[absPath] = struct{}{}
			defaultContents.ImageColPaths = append(defaultContents.ImageColPaths[:i], defaultContents.ImageColPaths[i+1:]...)
			break
		}
	}
	template := config.DefaultTemplate(configPath)
	template.YamlContents, _ = yaml.Marshal(defaultContents)
	err := template.WriteFile()
	if err != nil {
		return fmt.Errorf("error writing to default config: %s", err.Error())
	}
	if len(removed) == 0 {
		removed["no changes made"] = struct{}{}
	}
	defaultContents.Log(configPath).LogRemoved(removed)
	return nil
}
