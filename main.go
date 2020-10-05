package main

import (
	"fmt"
	"os"
)

type Command interface {
	Execute([]string) int
}

func main() {
	cfg, err := NewConfig(Opts{})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	var cmd Command
	if len(os.Args) < 2 {
		cmd = NewCommandList(cfg)
		os.Exit(cmd.Execute([]string{}))
	}

	switch os.Args[1] {
	case CommandNewName:
		cmd = NewCommandNewItem(cfg)
	case CommandListName:
		cmd = NewCommandList(cfg)
	case CommandGetName:
		cmd = NewCommandGet(cfg)
	default:
		cmd = NewHelpCommand(cfg)
	}

	os.Exit(cmd.Execute(os.Args[2:]))
}
