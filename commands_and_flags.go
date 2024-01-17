package main

import (
	"fmt"
	"os"
	"path/filepath"
	// "gopkg.in/yaml.v3"
)

type Command struct {
	name         string
	values       []string
	isValidValue func(string) error
}

// define commands here
var CLI_CMDS = []Command{
	ADD_CMD,
}
var (
	ADD_CMD = Command{
		name: "add",
		isValidValue: func(path string) error {
			// check if exists
			if _, err := os.Stat(path); os.IsNotExist(err) {
				return fmt.Errorf("%s does not exist: %s", path, err.Error())
			}

			// check if has any image files
			imgFileCount := 0
			err := filepath.WalkDir(path, func(p string, d os.DirEntry, err error) error {
				if err != nil {
					return err
				}

				if d.IsDir() && d.Name() != filepath.Base(path) {
					return filepath.SkipDir
				}

				if !d.IsDir() && IsImageFile(d.Name()) {
					imgFileCount++
				}

				return nil
			})

			if err != nil {
				return fmt.Errorf("error reading %s: %s", path, err.Error())
			}

			if imgFileCount < 1 {
				return fmt.Errorf("no image files in %s", path)
			}

			return nil
		},
	}

	CONFIG_CMD = Command{
		name: "config",
		isValidValue: func(s string) error {
			if s == "" || s == "default" {
				return nil
			}

			// check if exists
			if _, err := os.Stat(s); os.IsNotExist(err) {
				return fmt.Errorf("%s does not exist: %s", s, err.Error())
			}

			// check if has a config file (.yaml)
			configCount := 0
			err := filepath.WalkDir(s, func(p string, d os.DirEntry, err error) error {
				if err != nil {
					return err
				}

				if d.IsDir() && d.Name() != filepath.Base(s) {
					return filepath.SkipDir
				}

				if !d.IsDir() && filepath.Ext(p) == ".yaml" {
					// read file
					// unmarshal yaml
					fmt.Println("found:", d.Name())
					configCount++
				}

				return nil
			})

			if err != nil {
				return fmt.Errorf("error reading %s: %s", s, err.Error())
			}

			if configCount == 0 {
				return fmt.Errorf("no config files found in %s", s)

			} else if configCount > 1 {
				return fmt.Errorf("multiple config files found in %s", s)
			}

			return nil
		},
	}
)

func IsValidArgName(a string) bool {
	for _, ARG := range CLI_CMDS {
		if a != ARG.name {
			return false
		}
	}
	return true
}
