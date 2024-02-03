package main

import (
	"fmt"
	"os"

	"github.com/saltkid/tbg/cmd"
	"github.com/saltkid/tbg/flag"
)

func main() {
	tokens, err := TokenizeArgs(os.Args[1:])
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Tokens:")
	for _, token := range tokens {
		if token.isCmd {
			fmt.Println(cmd.CmdType(token.id).String(), token.value)
		} else if token.isFlag {
			fmt.Println(flag.FlagType(token.id).String(), token.value)
		}
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
			if tmpCmd, err := cmd.ToCommand(arg); err == nil {
				// is command
				tmpTok = Token{
					id:     uint8(tmpCmd.Type),
					isCmd:  true,
					isFlag: false,
				}
			} else if tmpFlag, err := flag.ToFlag(arg); err == nil {
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
			if tmpCmd, err := cmd.ToCommand(arg); err == nil {
				// encountered command instead of value
				tokens = append(tokens, tmpTok)
				tmpTok = Token{
					id:     uint8(tmpCmd.Type),
					isCmd:  true,
					isFlag: false,
				}
			} else if tmpFlag, err := flag.ToFlag(arg); err == nil {
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
