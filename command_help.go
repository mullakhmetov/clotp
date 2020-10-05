package main

import "fmt"

var help = `usage: clotp <command> [<args>]

There are clotp commands:
list - get available TOTP's
new - create new TOPT
get - get particular TOTP code by it's name

help - show this help
`

func NewHelpCommand(cfg *Config) *CommandHelp {
	return &CommandHelp{cfg}
}

type CommandHelp struct {
	cfg *Config
}

func (c CommandHelp) Execute(_ []string) int {
	fmt.Println(help)

	return 0
}
