package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/saltkid/tbg/config"
	"github.com/saltkid/tbg/flag"
)

func RemoveValidateValue(val string) error {
	_, err := filepath.Abs(val)
	if err != nil {
		return fmt.Errorf("Failed to get absolute path of %s: %s", val, err)
	}
	return nil
}

func RemoveValidateFlag(f *flag.Flag) error {
	switch f.Type {
	case flag.Alignment, flag.Opacity, flag.Stretch:
		if f.Value != "" {
			return fmt.Errorf("'remove' flags don't take any values. flag '%s' has value: '%s'", f.Type.ToString(), f.Value)
		}
		return nil
	default:
		return fmt.Errorf("invalid flag for 'remove': '%s'", f.Type.ToString())
	}
}

func RemoveValidateSubCmd(c *Cmd) error {
	switch c.Type {
	case None:
		return nil
	default:
		return fmt.Errorf("'remove' takes no sub commands. got: '%s'", c.Type.ToString())
	}
}

func RemoveExecute(c *Cmd) error {
	toRemove, _ := filepath.Abs(c.Value)
	toRemove = filepath.ToSlash(toRemove)

	// check if flags are set by user (empty if not)
	align := ExtractFlagValue(flag.Alignment, c.Flags)
	opacity := ExtractFlagValue(flag.Opacity, c.Flags)
	stretch := ExtractFlagValue(flag.Stretch, c.Flags)

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

	err = configContents.RemovePath(toRemove, configPath, align, stretch, opacity)
	if err != nil {
		return err
	}
	return nil
}
