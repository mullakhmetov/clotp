package main

import (
	"fmt"
	"os"

	"github.com/mullakhmetov/clotp/config"
)

type Command interface {
	Execute([]string) int
}

func main() {
	cfg, err := config.NewConfig(config.Opts{})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	var cmd Command
	if len(os.Args) < 2 {
		cmd = config.NewCommandList(cfg)
		os.Exit(cmd.Execute(os.Args[2:]))
	}

	switch os.Args[1] {
	case "new":
		cmd = config.NewCommandNewItem(cfg)
	case "list", "search":
		cmd = config.NewCommandList(cfg)
	case "get":
		cmd = config.NewCommandGet(cfg)
	}

	os.Exit(cmd.Execute(os.Args[2:]))
}
