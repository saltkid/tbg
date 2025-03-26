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
	err = command.Execute()
	if err != nil {
		fmt.Println(err)
		return
	}
}
