package config

import (
	"flag"
	"fmt"
)

type Command struct {
	fs     *flag.FlagSet
	config *Config
}

func NewConfigCommand(fs *flag.FlagSet) (*Command, error) {
	cfg, err := NewDefaultConfig()
	if err != nil {
		return nil, err
	}

	return &Command{fs, cfg}, nil
}

func (c *Command) Execute(args []string) int {
	if len(args) == 0 {

	}

	switch args[0] {
	case "list":
		fmt.Println(c.config.List())
		return 0
	}

	return 1
}
