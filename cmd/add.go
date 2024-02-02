package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/saltkid/tbg/config"
	"github.com/saltkid/tbg/flag"
	"gopkg.in/yaml.v3"
)

func AddValidateValue(val string) error {
	absPath, err := filepath.Abs(val)
	if err != nil {
		return fmt.Errorf("Failed to get absolute path of %s: %s", val, err)
	}

	if _, err := os.Stat(absPath); os.IsNotExist(err) {
		return fmt.Errorf("%s does not exist: %s", val, err.Error())
	}

	// path must have at least one image file
	hasImageFile := false
	err = filepath.WalkDir(absPath, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		// only search depth 1
		if d.IsDir() {
			return filepath.SkipDir
		}
		// find at least one
		if isImageFile(d.Name()) {
			hasImageFile = true
			return filepath.SkipAll
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("Failed to walk directory %s: %s", val, err)
	}
	if !hasImageFile {
		return fmt.Errorf("No image files found in %s", val)
	}
	return nil
}

func AddValidateFlag(f flag.Flag) error {
	switch f.Type {
	case flag.Alignment, flag.Opacity, flag.Stretch:
		return nil
	default:
		return fmt.Errorf("unexpected error: unknown flag: %d", f.Type)
	}
}

func AddExecute(c *Cmd) error {
	absPath, _ := filepath.Abs(c.Value)

	// check if flags are set by user
	align := ExtractFlagValue(flag.Alignment, c.Flags)
	opacity := ExtractFlagValue(flag.Opacity, c.Flags)
	stretch := ExtractFlagValue(flag.Stretch, c.Flags)
	// set flags after path only if at least one is set
	if align != "" || opacity != "" || stretch != "" {
		if align == "" {
			align = "center"
		}
		if opacity == "" {
			opacity = "0.1"
		}
		if stretch == "" {
			stretch = "uniform"
		}
		absPath = fmt.Sprintf("%s | %s %s %s", absPath, align, opacity, stretch)
	}

	// check if config subcommand is set by user
	specifiedConfig := ExtractSubCmdValue(Config, c.SubCmds)
	var configPath string
	var configContents config.Config
	if specifiedConfig == "default" || specifiedConfig == "" {
		configPath, _ = filepath.Abs("config.yaml")
		configContents = &config.DefaultConfig{}
	} else {
		configPath, _ = filepath.Abs(specifiedConfig)
		configContents = &config.UserConfig{}
	}

	yamlFile, err := os.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("Failed to read config file %s: %s", configPath, err)
	}
	err = configContents.Unmarshal(yamlFile)
	if err != nil {
		return err
	}

	if configContents.IsDefaultConfig() && specifiedConfig == "default" {
		err = AddToDefaultConfig(configContents, absPath, configPath)
		if err != nil {
			return err
		}
	} else if configContents.IsUserConfig() {
		err = AddToUserConfig(configContents, absPath, configPath)
		if err != nil {
			return err
		}
	} else {
		// read default config to check if using user config or not
		defaultContents, _ := configContents.(*config.DefaultConfig)

		// using default config
		if defaultContents.UserConfig == "" {
			err = AddToUserConfig(configContents, absPath, configPath)
			if err != nil {
				return err
			}

		} else {
			// recheck if user config path in default config is valid
			// TODO
			err = nil // ConfigValidateValue(defaultContents.UserConfig)
			if err != nil {
				return err
			}

			userConfigPath, _ := filepath.Abs(defaultContents.UserConfig)
			yamlFile, err = os.ReadFile(userConfigPath)
			if err != nil {
				return fmt.Errorf("Failed to read user config file %s: %s", userConfigPath, err)
			}
			contents := &config.UserConfig{}
			err = contents.Unmarshal(yamlFile)
			if err != nil {
				return err
			}
			err = AddToUserConfig(contents, absPath, userConfigPath)
			if err != nil {
				return err
			}
		}

	}

	return nil
}

func AddToUserConfig(contents config.Config, absPath string, configPath string) error {
	userContents, ok := contents.(*config.UserConfig)
	if !ok {
		return fmt.Errorf("unexpected error: contents is not a user config")
	}

	for _, path := range userContents.ImageColPaths {
		pureAbsPath, _, _ := strings.Cut(absPath, "|")
		pureAbsPath = strings.TrimSpace(pureAbsPath)
		purePath, _, _ := strings.Cut(path, "|")

		if pureAbsPath == purePath {
			return fmt.Errorf("%s already in user config", absPath)
		}
	}
	userContents.ImageColPaths = append(userContents.ImageColPaths, absPath)

	template := config.DefaultTemplate(configPath)
	template.YamlContents, _ = yaml.Marshal(userContents)
	err := template.WriteFile()
	if err != nil {
		return fmt.Errorf("error writing to user config: %s", err.Error())
	}
	userContents.Log(configPath)
	return nil
}

func AddToDefaultConfig(contents config.Config, absPath string, configPath string) error {
	defaultContents, ok := contents.(*config.DefaultConfig)
	if !ok {
		return fmt.Errorf("unexpected error: contents is not a default config")
	}

	for _, path := range defaultContents.ImageColPaths {
		pureAbsPath, _, _ := strings.Cut(absPath, "|")
		pureAbsPath = strings.TrimSpace(pureAbsPath)
		purePath, _, _ := strings.Cut(path, "|")

		if pureAbsPath == purePath {
			return fmt.Errorf("%s already in default config", absPath)
		}
	}
	defaultContents.ImageColPaths = append(defaultContents.ImageColPaths, absPath)

	template := config.DefaultTemplate(configPath)
	template.YamlContents, _ = yaml.Marshal(defaultContents)
	err := template.WriteFile()
	if err != nil {
		return fmt.Errorf("error writing to default config: %s", err.Error())
	}
	defaultContents.Log(configPath)
	return nil
}
