package main

import (
	"fmt"
	"os"
	"path/filepath"
)

type AddCommand struct {
	Path      string
	Alignment *string
	Stretch   *string
	Opacity   *float32
}

func (cmd *AddCommand) Type() CommandType { return AddCommandType }

func (r *AddCommand) Debug() {
	fmt.Println("Add Command:", r.Type())
	fmt.Println("Path:", r.Path)
	fmt.Println("Flags:")
	if r.Alignment != nil {
		fmt.Println(" ", AlignmentFlag, *r.Alignment)
	}
	if r.Stretch != nil {
		fmt.Println(" ", StretchFlag, *r.Stretch)
	}
	if r.Opacity != nil {
		fmt.Println(" ", OpacityFlag, *r.Opacity)
	}
}

func (cmd *AddCommand) ValidateValue(val *string) error {
	if val == nil {
		return fmt.Errorf("'add' must have an argument. got none")
	}
	absPath, err := filepath.Abs(*val)
	if err != nil {
		return fmt.Errorf("Failed to get absolute path of %s: %s", *val, err)
	}
	if _, err := os.Stat(absPath); os.IsNotExist(err) {
		return fmt.Errorf("%s does not exist: %s", *val, err.Error())
	}
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
		if IsImageFile(d.Name()) {
			hasImageFile = true
			return filepath.SkipAll
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("Failed to walk directory %s: %s", *val, err)
	}
	if !hasImageFile {
		return fmt.Errorf("No image files found in %s", *val)
	}
	cmd.Path = filepath.ToSlash(absPath)
	return nil
}

func (cmd *AddCommand) ValidateFlag(f Flag) error {
	switch f.Type {
	case AlignmentFlag:
		val, err := ValidateAlignment(f.Value)
		if err != nil {
			return err
		}
		cmd.Alignment = val
	case OpacityFlag:
		val, err := ValidateOpacity(f.Value)
		if err != nil {
			return err
		}
		cmd.Opacity = val
	case StretchFlag:
		val, err := ValidateStretch(f.Value)
		if err != nil {
			return err
		}
		cmd.Stretch = val
	default:
		return fmt.Errorf("invalid flag for 'add': '%s'", f.Type)
	}
	return nil
}

func (cmd *AddCommand) ValidateSubCommand(sc Command) error {
	switch sc.Type() {
	case NoCommandType:
		return nil
	default:
		return fmt.Errorf("'add' takes no sub commands. got: '%s'", sc.Type())
	}
}

func (cmd *AddCommand) Execute() error {
	configPath, err := ConfigPath()
	if err != nil {
		return err
	}
	yamlFile, err := os.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("Failed to read config file %s: %s", configPath, err)
	}
	config := new(Config)
	err = config.Unmarshal(yamlFile)
	if err != nil {
		return err
	}
	err = config.AddPath(configPath, cmd.Path, cmd.Alignment, cmd.Stretch, cmd.Opacity)
	if err != nil {
		return err
	}
	return nil
}
