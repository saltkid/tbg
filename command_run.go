package cmd

import (
	"fmt"
	"os"

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
	case flag.Profile, flag.Interval, flag.Alignment, flag.Opacity, flag.Stretch, flag.Random:
		return f.ValidateValue(f.Value)
	default:
		return fmt.Errorf("invalid flag for 'run': '%s'", f.Type.ToString())
	}
}

func RunValidateSubCmd(sc *Cmd) error {
	switch sc.Type {
	case None:
		return nil
	default:
		return fmt.Errorf("'run' takes no sub commands. got: '%s'", sc.Type.ToString())
	}
}

func RunExecute(c *Cmd) error {
	// get flags if set by user
	profile := ExtractFlagValue(flag.Profile, c.Flags)
	interval := ExtractFlagValue(flag.Interval, c.Flags)
	alignment := ExtractFlagValue(flag.Alignment, c.Flags)
	opacity := ExtractFlagValue(flag.Opacity, c.Flags)
	stretch := ExtractFlagValue(flag.Stretch, c.Flags)
	random := ExtractFlagValue(flag.Random, c.Flags)

	configPath, err := config.ConfigPath()
	if err != nil {
		return err
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

	err = configContents.ChangeBgImage(configPath, profile, interval, alignment, stretch, opacity, random)
	if err != nil {
		return err
	}
	return nil
}
