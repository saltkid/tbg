package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/saltkid/tbg/config"
	"github.com/saltkid/tbg/flag"
)

func EditValidateValue(val string) error {
	if val == "fields" || val == "all" {
		return nil
	}
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

func EditValidateFlag(f *flag.Flag) error {
	switch f.Type {
	case flag.Profile, flag.Interval, flag.Alignment, flag.Opacity, flag.Stretch:
		return f.ValidateValue(f.Value)
	default:
		return fmt.Errorf("unexpected error: unknown flag: '%s'", f.Type.ToString())
	}
}

func EditValidateSubCmd(c *Cmd) error {
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

func EditExecute(c *Cmd) error {
	toEdit, _ := filepath.Abs(c.Value)

	// check if flags are set by user
	profile := ExtractFlagValue(flag.Profile, c.Flags)
	interval := ExtractFlagValue(flag.Interval, c.Flags)
	alignment := ExtractFlagValue(flag.Alignment, c.Flags)
	opacity := ExtractFlagValue(flag.Opacity, c.Flags)
	stretch := ExtractFlagValue(flag.Stretch, c.Flags)

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

	err = configContents.EditPath(toEdit, configPath, profile, interval, alignment, stretch, opacity)
	if err != nil {
		return err
	}
	return nil
}
