package main

import (
	"fmt"
	"net/http"
)

type NextImageCommand struct{}

func (cmd *NextImageCommand) Type() CommandType { return NextImageCommandType }

func (r *NextImageCommand) Debug() {
	fmt.Println("Next Image Command:", r.Type())
}

func (cmd *NextImageCommand) ValidateValue(val *string) error {
	if val == nil || *val == "" {
		return nil
	}
	return fmt.Errorf("'next-image' takes no args. got: '%s'", *val)
}

func (cmd *NextImageCommand) ValidateFlag(f Flag) error {
	return fmt.Errorf("'next-image' takes no flags. got: '%s'", f.Type)
}

func (cmd *NextImageCommand) ValidateSubCommand(sc Command) error {
	switch sc.Type() {
	case NoCommandType:
		return nil
	default:
		return fmt.Errorf("'next-image' takes no sub commands. got: '%s'", sc.Type())
	}
}

func (cmd *NextImageCommand) Execute() error {
	url := fmt.Sprintf("http://127.0.0.1%s/next-image", TbgPort)
	resp, err := http.Post(url, "application/json", nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}
