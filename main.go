package main

import (
	"fmt"
	"log"
	"os"
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
	case "config", "add", "remove", "edit":
		err := command.run(command)
		if err != nil {
			log.Fatalln(err)
		}
	default:
		log.Fatalln("invalid command")
	}
}

func Run(command *Command) {
	// TODO
	log.Println("run")
	return
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
			_ = CONFIG_CMD.validateValue("default") // to create default config.yaml if not created yet
			return CONFIG_CMD, nil
		default:
			return nil, fmt.Errorf("command '%s' requires an argument. got none", args[0])
		}
	}

	// parse command name
	tmpCMD, err := ToCommand(args[0])
	if err != nil {
		return nil, err
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
				if _, isFlag := tmpArg.(*Flag); isFlag && tmpValue != "" {
					tmpArg.(*Flag).value = tmpValue

				} else if _, isCommand := tmpArg.(*Command); isCommand && tmpValue != "" {
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
				if _, isFlag := tmpArg.(*Flag); isFlag && tmpValue != "" {
					tmpArg.(*Flag).value = tmpValue
				} else if _, isCommand := tmpArg.(*Command); isCommand && tmpValue != "" {
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
					var errCtx string
					if _, isFlag := tmpArg.(*Flag); isFlag {
						errCtx = tmpArg.(*Flag).name

					} else if _, isCommand := tmpArg.(*Command); isCommand {
						errCtx = tmpArg.(*Command).name
					}

					return nil, fmt.Errorf("multiple values for flag: '%s' {%s, %s}", errCtx, tmpValue, arg)
				}
			}

			if i+1 == len(flagArgs) {
				if _, isFlag := tmpArg.(*Flag); isFlag && tmpValue != "" {
					tmpArg.(*Flag).value = tmpValue

				} else if _, isCommand := tmpArg.(*Command); isCommand && tmpValue != "" {
					tmpArg.(*Command).value = tmpValue
				}

				flags = append(flags, tmpArg)
			}
		}
	}

	// validate flags
	for _, flag := range flags {
		if _, isFlag := flag.(*Flag); isFlag {
			err = tmpCMD.validateFlag(flag.(*Flag).name, flag.(*Flag).value)
			if err != nil {
				return nil, err
			}

			_, ok := tmpCMD.flags[flag.(*Flag).name]
			if ok {
				return nil, fmt.Errorf("multiple flags with the same name: '%s'", flag.(*Flag).name)
			}
			tmpCMD.flags[flag.(*Flag).name] = flag

		} else if _, isCommand := flag.(*Command); isCommand {
			err = tmpCMD.validateFlag(flag.(*Command).name, flag.(*Command).value)
			if err != nil {
				return nil, err
			}

			_, ok := tmpCMD.flags[flag.(*Command).name]
			if ok {
				return nil, fmt.Errorf("multiple flags with the same name: '%s'", flag.(*Command).name)
			}
			tmpCMD.flags[flag.(*Command).name] = flag
		}
	}
	return tmpCMD, nil
}
