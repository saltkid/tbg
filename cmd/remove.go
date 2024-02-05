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

func RemoveValidateFlag(f *flag.Flag) error {
	switch f.Type {
	case flag.None:
		return f.ValidateValue(f.Value)
	default:
		return fmt.Errorf("'remove' takes no flags. got: '%s'", f.Type.ToString())
	}
}

func RemoveValidateSubCmd(c *Cmd) error {
	switch c.Type {
	case Config:
		if c.Value == "" {
			return fmt.Errorf("'config' subcommand requires a config file path")
		}
		return c.ValidateValue(c.Value)
	default:
		return fmt.Errorf("unexpected error: unknown sub command type: '%s'", c.Type.ToString())
	}
}

func RemoveExecute(c *Cmd) error {
	toRemove, _ := filepath.Abs(c.Value)

	// check if config subcommand is set by user (empty if not)
	specifiedConfig := ExtractSubCmdValue(Config, c.SubCmds)
	var configPath string
	var err error
	if specifiedConfig == "default" {
		configPath, err = config.DefaultConfigPath()
	} else if specifiedConfig == "" {
		configPath, err = config.UsedConfig()
	} else {
		configPath, err = filepath.Abs(specifiedConfig)
	}
	if err != nil {
		return fmt.Errorf("Failed to get config path: %s", err)
	}

	yamlFile, err := os.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("Failed to read config file %s: %s", configPath, err)
	}
	configContents := &config.Config{}
	err = configContents.Unmarshal(yamlFile)
	if err != nil {
		return err
	}

	err = configContents.RemovePath(toRemove, configPath)
	if err != nil {
		return err
	}
	return nil
}
