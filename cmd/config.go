package cmd

import (
	"fmt"
	"os"

	"github.com/saltkid/tbg/config"
	"github.com/saltkid/tbg/flag"
)

func ConfigValidateValue(val string) error {
	// default config
	switch val {
	case "edit", "":
		return nil
	default:
		return fmt.Errorf("invalid arg for 'config': '%s'", val)
	}
}

func ConfigValidateFlag(f *flag.Flag) error {
	switch f.Type {
	case flag.Profile, flag.Interval, flag.Alignment, flag.Opacity, flag.Stretch:
		return f.ValidateValue(f.Value)
	default:
		return fmt.Errorf("invalid flag for 'config': '%s'", f.Type.ToString())
	}
}

func ConfigValidateSubCmd(c *Cmd) error {
	switch c.Type {
	case None:
		return nil
	default:
		return fmt.Errorf("'config' takes no sub commands. got: '%s'", c.Type.ToString())
	}
}

func ConfigExecute(c *Cmd) error {
	switch c.Value {
	// print currently used config
	case "":
		configPath, err := config.ConfigPath()
		if err != nil {
			return err
		}
		yamlFile, err := os.ReadFile(configPath)
		if err != nil {
			return fmt.Errorf("Failed to read config at %s: %s", configPath, err)
		}
		configContents := &config.Config{}
		err = configContents.Unmarshal(yamlFile)
		if err != nil {
			return fmt.Errorf("Failed to unmarshal default config.yaml: %s", err)
		}
		configContents.Log(configPath)

	// edit config fields
	case "edit":

	default:
		return fmt.Errorf("unexpected error: invalid arg for 'config' after validation: '%s'", c.Value)
	}
	return nil
}
