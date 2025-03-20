package main

import (
	"fmt"
	"strings"
)

type Token struct {
	id     uint8
	value  *string
	isFlag bool
}

func TokenizeArgs(args []string) ([]Token, error) {
	tokens := make([]Token, 0)
	var tmpTok Token
	for i, arg := range args {
		tokenIsEmpty := tmpTok.id == 0
		if tokenIsEmpty {
			if strings.HasPrefix(arg, "-") {
				if tmpFlag, err := ToFlag(arg); err != nil {
					return nil, err
				} else {
					tmpTok = Token{
						id:     uint8(tmpFlag.Type),
						isFlag: true,
					}
				}
			} else if tmpCmd, err := ToCommand(arg); err != nil {
				// is neither
				return nil, err
			} else {
				tmpTok = Token{
					id:     uint8(tmpCmd.Type()),
					isFlag: false,
				}
			}
		} else {
			// there already is a command/flag
			// find a value for it
			if strings.HasPrefix(arg, "-") {
				if tmpFlag, err := ToFlag(arg); err != nil {
					return nil, err
				} else {
					// encountered flag instead of value
					tokens = append(tokens, tmpTok)
					tmpTok = Token{
						id:     uint8(tmpFlag.Type),
						isFlag: true,
					}
				}
			} else if tmpCmd, err := ToCommand(arg); err != nil {
				// is value
				tmpTok.value = &arg
				tokens = append(tokens, tmpTok)
				tmpTok = Token{}
			} else {
				// encountered command instead of value
				tokens = append(tokens, tmpTok)
				tmpTok = Token{
					id:     uint8(tmpCmd.Type()),
					isFlag: false,
				}
			}
		}
		lastItem := i == len(args)-1
		tokenIsNotEmpty := tmpTok.id != 0
		if lastItem && tokenIsNotEmpty {
			tokens = append(tokens, tmpTok)
		}
	}
	return tokens, nil
}

func ParseArgs(tokens []Token) (Command, error) {
	if len(tokens) == 0 {
		return new(HelpCommand), nil
	}
	var mainCommand Command
	for i, tok := range tokens {
		isCommandToken := i == 0
		if isCommandToken {
			if tok.isFlag {
				return nil, fmt.Errorf("Must start with a valid command. got flag: '%s'", FlagType(tok.id))
			} else {
				mainCommand = CommandType(tok.id).ToCommand()
				if err := mainCommand.ValidateValue(tok.value); err != nil {
					return nil, err
				}
			}
		} else {
			if tok.isFlag {
				flag := Flag{
					Type:  FlagType(tok.id),
					Value: tok.value,
				}
				if err := mainCommand.ValidateFlag(flag); err != nil {
					return nil, err
				}
			} else {
				subCmd := CommandType(tok.id).ToCommand()
				if err := mainCommand.ValidateSubCommand(subCmd); err != nil {
					return nil, err
				}
			}
		}
	}
	return mainCommand, nil
}

func LogTokens(tokens []Token) {
	fmt.Printf("Tokens (%d):\n", len(tokens))
	for _, token := range tokens {
		if token.isFlag {
			fmt.Print("| ", FlagType(token.id), ": ")
			if token.value != nil {
				fmt.Println(*token.value)
			} else {
				fmt.Println("no arg")
			}
		} else {
			fmt.Print("| ", FlagType(token.id), ": ")
			if token.value != nil {
				fmt.Println("", *token.value)
			} else {
				fmt.Println("no arg")
			}
		}
	}
	fmt.Println()
}
