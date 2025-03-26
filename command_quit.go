package main

import (
	"fmt"
	"net/http"
)

type QuitCommand struct{}

func (cmd *QuitCommand) Type() CommandType { return QuitCommandType }

func (r *QuitCommand) Debug() {
	fmt.Println("Quit Command Command:", r.Type())
}

func (cmd *QuitCommand) ValidateValue(val *string) error {
	if val == nil || *val == "" {
		return nil
	}
	return fmt.Errorf("'quit' takes no args. got: '%s'", *val)
}

func (cmd *QuitCommand) ValidateFlag(f Flag) error {
	return fmt.Errorf("'quit' takes no flags. got: '%s'", f.Type)
}

func (cmd *QuitCommand) ValidateSubCommand(sc Command) error {
	switch sc.Type() {
	case NoCommandType:
		return nil
	default:
		return fmt.Errorf("'quit' takes no sub commands. got: '%s'", sc.Type())
	}
}

func (cmd *QuitCommand) Execute() error {
	url := fmt.Sprintf("http://127.0.0.1%s/quit", TbgPort)
	resp, err := http.Post(url, "application/json", nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	fmt.Println("resp:", resp.Status)
	return nil
}
