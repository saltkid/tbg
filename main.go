package main

import (
	"fmt"
	"os"
)

func main() {
	tokens, err := TokenizeArgs(os.Args[1:])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	command, err := ParseArgs(tokens)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	err = command.Execute()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
