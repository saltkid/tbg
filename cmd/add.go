package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/saltkid/tbg/config"
	"github.com/saltkid/tbg/flag"
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
		if d.IsDir() && d.Name() != filepath.Base(absPath) {
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

func AddValidateFlag(f *flag.Flag) error {
	switch f.Type {
	case flag.Alignment, flag.Opacity, flag.Stretch:
		return f.ValidateValue(f.Value)
	default:
		return fmt.Errorf("unexpected error: unknown flag: '%s'", f.Type.ToString())
	}
}

func AddValidateSubCmd(c *Cmd) error {
	switch c.Type {
	case Config:
		if c.Value == "" {
			return fmt.Errorf("'config' subcommand requires a config file path")
		}
		return c.ValidateValue(c.Value)
	default:
		return fmt.Errorf("unexpected error: unknown sub command: '%s'", c.Type.ToString())
	}
}

func AddExecute(c *Cmd) error {
	toAdd, _ := filepath.Abs(c.Value)

	// check if flags are set by user
	align := ExtractFlagValue(flag.Alignment, c.Flags)
	opacity := ExtractFlagValue(flag.Opacity, c.Flags)
	stretch := ExtractFlagValue(flag.Stretch, c.Flags)

	// check if config subcommand is set by user
	specifiedConfig := ExtractSubCmdValue(Config, c.SubCmds)
	var configPath string
	var configContents config.Config
	if specifiedConfig == "default" || specifiedConfig == "" {
		configPath, _ = filepath.Abs(config.DefaultConfigPath())
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

	if specifiedConfig == "" {
		// read default config to check if using user config or default
		defaultContents, _ := configContents.(*config.DefaultConfig)

		if defaultContents.UserConfig == "" {
			// using default config
			err = defaultContents.AddPath(toAdd, configPath, align, stretch, opacity)
		} else {
			// using user config
			err = defaultContents.AddPath(toAdd, defaultContents.UserConfig, align, stretch, opacity)
		}
	} else {
		err = configContents.AddPath(toAdd, configPath, align, stretch, opacity)
	}

	if err != nil {
		return err
	}
	return nil
}
