package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/saltkid/tbg/config"
	"github.com/saltkid/tbg/flag"
	"github.com/saltkid/tbg/utils"
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
		if utils.IsImageFile(d.Name()) {
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

	// check if flags are set by user (empty if not)
	align := ExtractFlagValue(flag.Alignment, c.Flags)
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

	err = configContents.AddPath(toAdd, configPath, align, stretch, opacity)
	if err != nil {
		return err
	}
	return nil
}

func AddHelp(verbose bool) {
	fmt.Printf("%-60s%s", "  add path/to/images/dir",
		"Adds a path containing images to currently used config\n")
	if verbose {
		fmt.Println("\n\n  Path to images dir should have at least one image")
		fmt.Println("  file under it. All subdirectories will be ignored.")
		fmt.Println("\n  You can specify which config to add to using the 'config' subcommand.")
		fmt.Println("  If you do not specify a config, the currently used config will be used.")
		fmt.Println("\n  You can specify alignment, stretch, and opacity using flags.")
		fmt.Println("\n  Examples:")
		fmt.Println("  1. tbg add path/to/images/dir")
		fmt.Println("     This is how it would look like in the config:")
		fmt.Println("      ----------------------")
		fmt.Println("      | image_col_paths:")
		fmt.Println("      |   - /path/to/images/dir")
		fmt.Println("      |")
		fmt.Println("      | other fields...")
		fmt.Println("      ----------------------")
		fmt.Println("\n  2. tbg add path/to/images/dir --alignment center --opacity 0.5 --stretch uniform")
		fmt.Println("     This is how it would look like in the config:")
		fmt.Println("      ----------------------")
		fmt.Println("      | image_col_paths:")
		fmt.Println("      |   - /path/to/images/dir | center uniform 0.5")
		fmt.Println("      |")
		fmt.Println("      | other fields...")
		fmt.Println("      ----------------------")
		fmt.Println("\n  3. tbg add path/to/images/dir --alignment top")
		fmt.Println("     This is how it would look like in the config:")
		fmt.Println("      ----------------------")
		fmt.Println("      | image_col_paths:")
		fmt.Println("      |   - /path/to/images/dir | top fill 0.1")
		fmt.Println("      |")
		fmt.Println("      | default_alignment: right")
		fmt.Println("      | default_stretch: fill")
		fmt.Println("      | default_alignment: 0.1")
		fmt.Println("      |")
		fmt.Println("      | other fields...")
		fmt.Println("      ----------------------")
		fmt.Println("     Notice that even though there is only one flag specified, there are 3 after |")
		fmt.Print("     This is because it inerited the default values for flags that were not specified\n\n")
	}
}
