package main

import (
	"fmt"
	"os"
)

func RunValidateValue(val string) error {
	if val == "" {
		return nil
	}
	return fmt.Errorf("'run' takes no arguments. got: '%s'", val)
}

func RunValidateFlag(f *Flag) error {
	switch f.Type {
	case Profile, Interval, Alignment, Opacity, Stretch, Random:
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
	profile := ExtractFlagValue(Profile, c.Flags)
	interval := ExtractFlagValue(Interval, c.Flags)
	alignment := ExtractFlagValue(Alignment, c.Flags)
	opacity := ExtractFlagValue(Opacity, c.Flags)
	stretch := ExtractFlagValue(Stretch, c.Flags)
	random := ExtractFlagValue(Random, c.Flags)

	configPath, err := ConfigPath()
	if err != nil {
		return err
	}
	yamlFile, err := os.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("Failed to read config file %s: %s", configPath, err)
	}
	configContents := &Config{}
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
