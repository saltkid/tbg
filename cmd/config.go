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
	if val == "default" || val == "" {
		configPath, err := filepath.Abs(config.DefaultConfigPath())
		if err != nil {
			return fmt.Errorf("Failed to get absolute path of config.yaml: %s", err)
		}

		if _, err := os.Stat(configPath); os.IsNotExist(err) {
			err = config.DefaultTemplate(configPath).WriteFile()
			if err != nil {
				return fmt.Errorf("Failed to create default config: %s", err)
			}
		}

		return nil
	}

	// user config
	configPath, err := filepath.Abs(val)
	if err != nil {
		return fmt.Errorf("Failed to get absolute path of %s: %s", val, err)
	}
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		err := os.MkdirAll(filepath.Dir(configPath), os.ModePerm)
		if err != nil {
			return fmt.Errorf("Failed to create parent dirs of %s: %s", configPath, err)
		}
		err = config.UserTemplate(configPath).WriteFile()
		if err != nil {
			return fmt.Errorf("Failed to create user config: %s", err)
		}
	}
	return nil
}

func ConfigValidateFlag(f flag.Flag) error {
	switch f.Type {
	case 0: // none
		return nil
	default:
		return fmt.Errorf("'config' has no flags. got type %d", f.Type)
	}
}

func ConfigExecute(c *Cmd) error {
	switch c.Value {
	// print currently used config
	case "":
		configPath, _ := filepath.Abs(config.DefaultConfigPath())
		yamlFile, err := os.ReadFile(configPath)
		if err != nil {
			return fmt.Errorf("Failed to read default config.yaml: %s", err)
		}
		defaultContents := &config.DefaultConfig{}
		err = defaultContents.Unmarshal(yamlFile)
		if err != nil {
			return fmt.Errorf("Failed to unmarshal default config.yaml: %s", err)
		}

		// print user config if set
		if defaultContents.UserConfig != "" {
			userConfigPath, _ := filepath.Abs(defaultContents.UserConfig)
			yamlFile, err := os.ReadFile(userConfigPath)
			if err != nil {
				return fmt.Errorf("Failed to read user config.yaml: %s", err)
			}
			userContents := &config.UserConfig{}
			err = userContents.Unmarshal(yamlFile)
			if err != nil {
				return fmt.Errorf("Failed to unmarshal user config.yaml: %s", err)
			}
			userContents.Log(userConfigPath)
		} else {
			defaultContents.Log(configPath)
		}

	// set currently used config to default config
	case "default":
		configPath, _ := filepath.Abs(config.DefaultConfigPath())
		yamlFile, err := os.ReadFile(configPath)
		if err != nil {
			return fmt.Errorf("Failed to read default config.yaml: %s", err)
		}
		defaultContents := &config.DefaultConfig{}
		err = defaultContents.Unmarshal(yamlFile)
		if err != nil {
			return fmt.Errorf("Failed to unmarshal default config.yaml: %s", err)
		}

		// unset user config since we're using default config now
		defaultContents.UserConfig = ""
		UpdateDefaultConfig(defaultContents, configPath)

	// set currently used config to user config passed by user
	default:
		configPath, _ := filepath.Abs(c.Value)
		yamlFile, err := os.ReadFile(configPath)
		if err != nil {
			return fmt.Errorf("Failed to read user config at %s: %s", c.Value, err)
		}
		userContents := &config.UserConfig{}
		err = userContents.Unmarshal(yamlFile)
		if err != nil {
			return fmt.Errorf("Failed to unmarshal user config at %s: %s", c.Value, err)
		}
		userContents.Log(configPath)

		// edit default config to use this user config
		configPath, _ = filepath.Abs(config.DefaultConfigPath())
		yamlFile, err = os.ReadFile(configPath)
		if err != nil {
			return fmt.Errorf("Failed to read default config.yaml: %s", err)
		}
		defaultContents := &config.DefaultConfig{}
		err = defaultContents.Unmarshal(yamlFile)
		if err != nil {
			return fmt.Errorf("Failed to unmarshal default config.yaml: %s", err)
		}
		defaultContents.UserConfig = configPath
		UpdateDefaultConfig(defaultContents, configPath)
	}
	return nil
}

func UpdateDefaultConfig(contents *config.DefaultConfig, configPath string) error {
	template := config.DefaultTemplate(configPath)
	template.YamlContents, _ = yaml.Marshal(contents)
	err := template.WriteFile()
	if err != nil {
		return fmt.Errorf("error writing to config: %s", err.Error())
	}
	contents.Log(configPath)
	return nil
}
