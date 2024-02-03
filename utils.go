package main

import (
	"fmt"

	"github.com/saltkid/tbg/cmd"
	"github.com/saltkid/tbg/flag"
)

func LogTokens(tokens []Token) {
	fmt.Println("Tokens:")
	for _, token := range tokens {
		if token.isCmd {
			fmt.Println("|", cmd.CmdType(token.id).ToString(), token.value)
		} else if token.isFlag {
			fmt.Println("|", flag.FlagType(token.id).ToString(), token.value)
		}
	}
}

func LogArgs(mainCmd *cmd.Cmd) {
	fmt.Println("Main Command:", cmd.CmdType(mainCmd.Type).ToString())
	fmt.Println("       Value:", mainCmd.Value)
	fmt.Println("Sub Commands:")
	for _, c := range mainCmd.SubCmds {
		fmt.Println("|", cmd.CmdType(c.Type).ToString(), c.Value)
	}
	fmt.Println("Flags:")
	for _, f := range mainCmd.Flags {
		fmt.Println("|", flag.FlagType(f.Type).ToString(), f.Value)
	}
}
