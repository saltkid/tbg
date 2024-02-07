package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/saltkid/tbg/config"
	"github.com/saltkid/tbg/flag"
)

func RunValidateValue(val string) error {
	if val == "" {
		return nil
	}
	return fmt.Errorf("'run' takes no arguments. got: '%s'", val)
}

func RunValidateFlag(f *flag.Flag) error {
	switch f.Type {
	case flag.Profile, flag.Interval, flag.Alignment, flag.Opacity, flag.Stretch:
		return f.ValidateValue(f.Value)
	default:
		return fmt.Errorf("unexpected error: unknown flag: '%s'", f.Type.ToString())
	}
}

func RunValidateSubCmd(sc *Cmd) error {
	switch sc.Type {
	case Config:
		if sc.Value == "" {
			return fmt.Errorf("'config' subcommand requires a config file path")
		}
		return sc.ValidateValue(sc.Value)
	default:
		return fmt.Errorf("unexpected error: unknown sub command: '%s'", sc.Type.ToString())
	}
}

func RunExecute(c *Cmd) error {
	// get flags if set by user
	profile := ExtractFlagValue(flag.Profile, c.Flags)
	interval := ExtractFlagValue(flag.Interval, c.Flags)
	alignment := ExtractFlagValue(flag.Alignment, c.Flags)
	opacity := ExtractFlagValue(flag.Opacity, c.Flags)
	stretch := ExtractFlagValue(flag.Stretch, c.Flags)

	// check if config subcommand is set by user (empty if not)
	specifiedConfig := ExtractSubCmdValue(Config, c.SubCmds)
	var configPath string
	var err error
	if specifiedConfig == nil {
		configPath, err = config.UsedConfig()
	} else if *specifiedConfig == "default" {
		configPath, err = config.DefaultConfigPath()
	} else {
		configPath, err = filepath.Abs(*specifiedConfig)
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

	err = configContents.ChangeBgImage(configPath, profile, interval, alignment, stretch, opacity)
	if err != nil {
		return err
	}
	return nil
}
