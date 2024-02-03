package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/saltkid/tbg/config"
	"github.com/saltkid/tbg/flag"
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
	toRemove, _ := filepath.Abs(c.Value)

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

	if specifiedConfig == "" {
		// read default config to check if using user config or default
		defaultContents, _ := configContents.(*config.DefaultConfig)

		// using default config
		if defaultContents.UserConfig == "" {
			err = defaultContents.RemovePath(toRemove, configPath)
		} else {
			err = defaultContents.RemovePath(toRemove, defaultContents.UserConfig)
		}
	} else {
		err = configContents.RemovePath(toRemove, configPath)
	}

	if err != nil {
		return err
	}
	return nil
}
