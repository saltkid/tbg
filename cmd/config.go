package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/saltkid/tbg/config"
	"github.com/saltkid/tbg/flag"
	"gopkg.in/yaml.v3"
)

func ConfigValidateValue(val string) error {
	// default config
	var configPath string
	var err error
	switch val {
	default:
		configPath, err = filepath.Abs(val)
		if err != nil {
			return fmt.Errorf("Failed to get absolute path of %s: %s", val, err)
		}
		if filepath.Ext(configPath) != ".yaml" && filepath.Ext(configPath) != ".yml" {
			configPath = configPath + ".yaml"
			fmt.Printf("Creating \"%s\" instead because \"%s\" does not have .yaml or .yml extension.\n", configPath, val)
		}
		if _, err := os.Stat(configPath); os.IsNotExist(err) {
			err := os.MkdirAll(filepath.Dir(configPath), os.ModePerm)
			if err != nil {
				return fmt.Errorf("Failed to create parent dirs of %s: %s", configPath, err)
			}
			err = config.NewConfigTemplate(configPath).WriteFile()
			if err != nil {
				return fmt.Errorf("Failed to create user config: %s", err)
			}
		}

	case "edit":
		return nil
	case "default":
		configPath, err = config.DefaultConfigPath()
	case "":
		configPath, err = config.UsedConfig()
	}

	if err != nil {
		return err
	}
	return nil
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
		configPath, _ := config.UsedConfig()
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

	// set currently used config to default config
	case "default":
		profilePath, _ := config.TbgProfilePath()
		configPath := filepath.Join(filepath.Dir(profilePath), "config.yaml")
		// edit tbg profile to use default config
		yamlFile, err := os.ReadFile(profilePath)
		if err != nil {
			return fmt.Errorf("Failed to read tbg profile: %s", err)
		}
		profileContents := &config.TbgProfile{}
		err = profileContents.Unmarshal(yamlFile)
		if err != nil {
			return fmt.Errorf("Failed to unmarshal tbg profile: %s", err)
		}
		profileContents.UsedConfig = configPath
		UpdateTbgProfile(profileContents, profilePath, configPath, yamlFile)

	// set currently used config to user config passed by user
	default:
		configPath, _ := filepath.Abs(c.Value)
		// edit tbg profile config to use this config
		profilePath, err := config.TbgProfilePath()
		if err != nil {
			return err
		}
		yamlFile, err := os.ReadFile(profilePath)
		if err != nil {
			return fmt.Errorf("Failed to read tbg profile: %s", err)
		}
		profileContents := &config.TbgProfile{}
		err = profileContents.Unmarshal(yamlFile)
		if err != nil {
			return fmt.Errorf("Failed to unmarshal tbg profile: %s", err)
		}
		profileContents.UsedConfig = configPath
		UpdateTbgProfile(profileContents, profilePath, configPath, yamlFile)
	}
	return nil
}

func UpdateTbgProfile(profileContents *config.TbgProfile, profilePath string, configPath string, yamlFile []byte) error {
	// write to tbg profile
	updatedProfile, err := yaml.Marshal(profileContents)
	if err != nil {
		return fmt.Errorf("Failed to marshal tbg profile: %s", err)
	}
	err = os.WriteFile(profilePath, updatedProfile, 0644)

	// log the user passed config
	yamlFile, err = os.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("Failed to read user config at %s: %s", configPath, err)
	}
	configContents := &config.Config{}
	err = configContents.Unmarshal(yamlFile)
	if err != nil {
		return fmt.Errorf("Failed to unmarshal user config at %s: %s", configPath, err)
	}
	configContents.Log(configPath)

	return nil
}
