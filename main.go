package main

import (
	"fmt"
	"log"
	"os"
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
