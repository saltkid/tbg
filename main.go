package main

import (
	"fmt"
	"os"
)

func main() {
	tokens, err := TokenizeArgs(os.Args[1:])
	if err != nil {
		fmt.Println(err)
		return
	}

	command, err := ParseArgs(tokens)
	if err != nil {
		fmt.Println(err)
		return
	}

	if command.IsNone() {
		command.Type = Help
	}

	Cls()
	err = command.Execute()
	if err != nil {
		fmt.Println(err)
		return
	}
}

type Token struct {
	id     uint8
	value  string
	isCmd  bool
	isFlag bool
}

func TokenizeArgs(args []string) ([]Token, error) {
	tokens := make([]Token, 0)
	var tmpTok Token

	for i, arg := range args {
		// id == 0 means empty token

		if tmpTok.id == 0 {
			if tmpCmd, err := ToCommand(arg); err == nil {
				// is command
				tmpTok = Token{
					id:     uint8(tmpCmd.Type),
					isCmd:  true,
					isFlag: false,
				}
			} else if tmpFlag, err := ToFlag(arg); err == nil {
				// is flag
				tmpTok = Token{
					id:     uint8(tmpFlag.Type),
					isCmd:  false,
					isFlag: true,
				}
			} else {
				// is neither
				return nil, fmt.Errorf("'%s' is not a valid command or flag", arg)
			}

			// last item so append token if not empty
			if i == len(args)-1 && tmpTok.id != 0 {
				tokens = append(tokens, tmpTok)
			}
		} else {
			if tmpCmd, err := ToCommand(arg); err == nil {
				// encountered command instead of value
				tokens = append(tokens, tmpTok)
				tmpTok = Token{
					id:     uint8(tmpCmd.Type),
					isCmd:  true,
					isFlag: false,
				}
			} else if tmpFlag, err := ToFlag(arg); err == nil {
				// encountered flag instead of value
				tokens = append(tokens, tmpTok)
				tmpTok = Token{
					id:     uint8(tmpFlag.Type),
					isCmd:  false,
					isFlag: true,
				}
			} else {
				// is value
				tmpTok.value = arg
				tokens = append(tokens, tmpTok)
				tmpTok = Token{}
			}

			// last item so append token if not empty
			if i == len(args)-1 && tmpTok.id != 0 {
				tokens = append(tokens, tmpTok)
			}
		}
	}
	return tokens, nil
}

func ParseArgs(tokens []Token) (*Cmd, error) {
	mainCommand := Cmd{
		Type:    None,
		SubCmds: make(map[CmdType]*Cmd, 0),
		Flags:   make(map[FlagType]*Flag, 0),
	}

	for i, tok := range tokens {
		if i == 0 && tok.isFlag {
			return nil, fmt.Errorf("must start with a valid command. got flag: '%s'", FlagType(tok.id).ToString())
		}

		if mainCommand.IsNone() {
			mainCommand.Type = CmdType(tok.id)
			err := mainCommand.ValidateValue(tok.value)
			if err != nil {
				return nil, err
			}
			mainCommand.Value = tok.value

		} else {
			if tok.isCmd {
				subCmd := &Cmd{
					Type:  CmdType(tok.id),
					Value: tok.value,
				}
				err := mainCommand.ValidateSubCmd(subCmd)
				if err != nil {
					return nil, err
				}
				mainCommand.SubCmds[subCmd.Type] = subCmd
			} else if tok.isFlag {
				flag := &Flag{
					Type:  FlagType(tok.id),
					Value: tok.value,
				}
				err := mainCommand.ValidateFlag(flag)
				if err != nil {
					return nil, err
				}
				mainCommand.Flags[flag.Type] = flag
			}
		}
	}
	return &mainCommand, nil
}

func LogTokens(tokens []Token) {
	fmt.Println("Tokens:")
	for _, token := range tokens {
		if token.isCmd {
			fmt.Println("|", CmdType(token.id).ToString(), token.value)
		} else if token.isFlag {
			fmt.Println("|", FlagType(token.id).ToString(), token.value)
		}
	}
}

func LogArgs(mainCmd *Cmd) {
	fmt.Println("Main Command:", CmdType(mainCmd.Type).ToString())
	fmt.Println("       Value:", mainCmd.Value)
	fmt.Println("Sub Commands:")
	for _, c := range mainCmd.SubCmds {
		fmt.Println("|", CmdType(c.Type).ToString(), c.Value)
	}
	fmt.Println("Flags:")
	for _, f := range mainCmd.Flags {
		fmt.Println("|", FlagType(f.Type).ToString(), f.Value)
	}
}
