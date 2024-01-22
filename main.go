package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

func main() {
	command, err := ParseArgs(os.Args[1:])
	if err != nil {
		log.Fatalln(err)
	}

	LogParsedArgs(command)
	switch command.name {
	case "run":
		Run(command)
	case "add":
		Add(command)
	case "config":
		Config(command)
	default:
		log.Fatalln("invalid command")
	}
}

func Run(command *Command) {
	// TODO
	log.Println("run")
	return
}

func Add(command *Command) {
	// TODO
	log.Println("add")
	return
}

func Config(command *Command) {
	var configPath string
	if command.value == "default" {
		configPath, _ = filepath.Abs("config.yaml")
	} else {
		configPath, _ = filepath.Abs(command.value)
	}

	yamlFile, _ := os.ReadFile(configPath)
	contents := ConfigFile{}
	err := yaml.Unmarshal(yamlFile, &contents)
	if err != nil {
		log.Fatalln(err)
	}

	LogConfig(contents, configPath)
}

func ParseArgs(args []string) (*Command, error) {
	// default command if no args
	if len(args) < 1 {
		return RUN_CMD, nil
	}

	// commands that allow no args
	if len(args) < 2 {
		tmp, err := ToCommand(args[0])
		if err != nil {
			return nil, err
		}

		switch tmp.name {
		case "run":
			return RUN_CMD, nil
		case "config":
			return CONFIG_CMD, nil
		default:
			return nil, fmt.Errorf("only 'run' command is allowed to be empty. got '%s'", args[0])
		}
	}

	// parse command name
	tmpCMD := new(Command)
	tmp, err := ToCommand(args[0])
	if err != nil {
		return nil, err
	}
	switch tmp.name {
	case "run":
		tmpCMD = RUN_CMD
	case "add":
		tmpCMD = ADD_CMD
	case "config":
		tmpCMD = CONFIG_CMD
	default:
		return nil, fmt.Errorf("start with invalid command: '%s'", args[0])
	}

	// parse command value
	flagArgs := args[1:]
	for i, arg := range flagArgs {
		if IsValidFlagName(arg) || IsValidCommandName(arg) {
			flagArgs = flagArgs[i:]
			break
		} else {
			err := tmpCMD.validateValue(arg)
			if err != nil {
				return nil, err
			}
			tmpCMD.value = arg

			if i+1 < len(args) {
				flagArgs = flagArgs[i+1:]
			}
			break
		}
	}

	// parse flags
	flags := make([]CLI_Arg, 0)
	var tmpArg CLI_Arg
	tmpValue := ""
	start := true

	for i, arg := range flagArgs {
		if arg == "" {
			continue
		}

		if start {
			if IsValidFlagName(arg) {
				tmp, err := ToFlag(arg)
				if err != nil {
					return nil, err
				}
				tmpArg = tmp
				start = false

			} else if IsValidCommandName(arg) {
				tmp, err := ToCommand(arg)
				if err != nil {
					return nil, err
				}
				tmpArg = tmp
				start = false

			} else {
				return nil, fmt.Errorf("'%s' is an invalid flag for '%s'", arg, tmpCMD.name)
			}

			if i+1 == len(flagArgs) {
				flags = append(flags, tmpArg)
			}

		} else {
			if IsValidFlagName(arg) {
				_, isFlag := tmpArg.(*Flag)
				_, isCommand := tmpArg.(*Command)
				if isFlag && tmpValue != "" {
					tmpArg.(*Flag).value = tmpValue
				} else if isCommand && tmpValue != "" {
					tmpArg.(*Command).value = tmpValue
				}

				flags = append(flags, tmpArg)

				tmp, err := ToFlag(arg)
				if err != nil {
					return nil, err
				}
				tmpArg = tmp
				tmpValue = ""

			} else if IsValidCommandName(arg) {
				_, isFlag := tmpArg.(*Flag)
				_, isCommand := tmpArg.(*Command)
				if isFlag && tmpValue != "" {
					tmpArg.(*Flag).value = tmpValue
				} else if isCommand && tmpValue != "" {
					tmpArg.(*Command).value = tmpValue
				}

				flags = append(flags, tmpArg)

				tmp, err := ToCommand(arg)
				if err != nil {
					return nil, err
				}
				tmpArg = tmp
				tmpValue = ""

			} else {
				if tmpValue == "" {
					tmpValue = arg

				} else {
					_, isFlag := tmpArg.(*Flag)
					_, isCommand := tmpArg.(*Command)
					var errCtx string
					if isFlag {
						errCtx = tmpArg.(*Flag).name
					} else if isCommand {
						errCtx = tmpArg.(*Command).name
					}

					return nil, fmt.Errorf("multiple values for flag: '%s' {%s, %s}", errCtx, tmpValue, arg)
				}
			}

			if i+1 == len(flagArgs) {
				_, isFlag := tmpArg.(*Flag)
				_, isCommand := tmpArg.(*Command)
				if isFlag && tmpValue != "" {
					tmpArg.(*Flag).value = tmpValue
				} else if isCommand && tmpValue != "" {
					tmpArg.(*Command).value = tmpValue
				}

				flags = append(flags, tmpArg)
			}
		}
	}

	// validate flags

	// if any flag is a command, validate the following flags with the subcommand
	subCMD := tmpCMD
	for _, flag := range flags {
		var err error
		_, isFlag := flag.(*Flag)
		_, isCommand := flag.(*Command)
		if isFlag {
			err = subCMD.validateFlag(flag.(*Flag).name, flag.(*Flag).value)
		} else if isCommand {
			err = subCMD.validateFlag(flag.(*Command).name, flag.(*Command).value)
			subCMD = flag.(*Command)
		}

		if err != nil {
			return nil, err
		}
		tmpCMD.flags = append(tmpCMD.flags, flag)
	}
	return tmpCMD, nil
}
