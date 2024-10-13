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
	// LogTokens(tokens)
	command, err := ParseArgs(tokens)
	if err != nil {
		fmt.Println(err)
		return
	}
	// command.Debug()
	ClearScreen()
	err = command.Execute()
	if err != nil {
		fmt.Println(err)
		return
	}
}
