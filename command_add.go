package main

import (
	"fmt"
	"os"
	"path/filepath"
)

type AddCommand struct {
	// raw input of user which may or may not have ~ and environment
	// variables, both of which will be kept unexpanded.
	Path string
	// normalized path from user input which expands both ~ and environment
	// variables.
	CleanPath string
	// path to a custom config path
	Config    *string
	Alignment *string
	Stretch   *string
	Opacity   *float32
}

func (cmd *AddCommand) Type() CommandType { return AddCommandType }

func (cmd *AddCommand) String() {
	fmt.Println("Add Command:", cmd.Type())
	fmt.Println("Path:", cmd.Path)
	fmt.Println("Flags:")
	if cmd.Alignment != nil {
		fmt.Println(" ", AlignmentFlag, *cmd.Alignment)
	}
	if cmd.Config != nil {
		fmt.Println(" ", ConfigFlag, *cmd.Config)
	}
	if cmd.Opacity != nil {
		fmt.Println(" ", OpacityFlag, *cmd.Opacity)
	}
	if cmd.Stretch != nil {
		fmt.Println(" ", StretchFlag, *cmd.Stretch)
	}
}

func (cmd *AddCommand) ValidateValue(val *string) error {
	if val == nil || *val == "" {
		return fmt.Errorf("'add' must have an argument. got none")
	}
	absPath, err := NormalizePath(*val)
	if err != nil {
		return fmt.Errorf("Failed to normalize path %s: %s", *val, err)
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
		if IsImageFile(filepath.Join(absPath, d.Name())) {
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
	cmd.Path = *val
	cmd.CleanPath = absPath
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
	case ConfigFlag:
		val, err := ValidateConfig(f.Value)
		if err != nil {
			return err
		}
		cmd.Config = val
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
	config, configPath, err := ConfigInit(cmd.Config)
	if err != nil {
		return err
	}
	return config.AddPath(configPath, cmd.Path, cmd.CleanPath, cmd.Alignment, cmd.Opacity, cmd.Stretch)
}
