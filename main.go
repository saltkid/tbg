package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func main() {
	cmd, err := ParseArgs(os.Args[1:])
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(cmd)
	// if cmd.IsNone() {
	rename()
	// }
}

type Command struct {
	name         string
	values       []string
	isValidValue func(string) error
}

// define commands here
var CLI_CMDS = []Command{ADD_CMD}
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

				if !d.IsDir() && d.Name() != filepath.Base(path) && IsImageFile(d.Name()) {
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
)

func IsValidArgName(a string) bool {
	for _, ARG := range CLI_CMDS {
		if a != ARG.name {
			return false
		}
	}
	return true
}

func ParseArgs(rawArgs []string) (Command, error) {
	var parsedCommand Command
	for i, arg := range rawArgs {
		if !IsValidArgName(arg) {
			return Command{}, fmt.Errorf("'%s' is not a valid command/switch/flag", arg)
		}

		if arg == ADD_CMD.name {
			newArg := ADD_CMD
			for _, argOpt := range rawArgs[i+1:] {
				if IsValidArgName(argOpt) {
					break
				}

				err := newArg.isValidValue(argOpt)
				if err != nil {
					return Command{}, err
				}
				newArg.values = append(newArg.values, argOpt)
			}
			parsedCommand = newArg
		}
	}

	return parsedCommand, nil
}
