package config

import "fmt"

func NewCommandGet(cfg *Config) *CommandGet {
	return &CommandGet{cfg}
}

type CommandGet struct {
	cfg *Config
}

func (c CommandGet) Execute(args []string) int {
	if len(args) != 1 {
		fmt.Printf("invalid TOTP name input: %s\n", args)
		return 1
	}

	name := args[0]

	var item *Item
	for _, i := range c.cfg.Items {
		if i.Name == name {
			item = i
			break
		}
	}

	if item != nil {
		fmt.Printf("%+v", item)
		return 0
	}

	fmt.Printf("unknown TOTP name: %s\n", name)
	return 1
}
