package main

import (
	"fmt"
	"strings"

	"github.com/AlecAivazis/survey/v2"
)

const CommandListName = "list"

func NewCommandList(cfg *Config) *CommandList {
	return &CommandList{cfg}
}

type CommandList struct {
	cfg *Config
}

func (c CommandList) Help() {
	fmt.Println("")
}

func (c CommandList) Execute(_ []string) int {
	m := make(map[string]*Item)
	options := make([]string, 0, len(c.cfg.Items))
	for _, i := range c.cfg.Items {
		m[i.Name] = i
		options = append(options, i.Name)
	}

	if len(options) == 0 {
		fmt.Println("You have no configured TOTP entities. Run `new` command to create one")
		return 1
	}

	q := &survey.Select{
		Message: "Choose a TOTP name:",
		Options: options,
		Filter:  myFilter,
	}

	var name string
	if err := survey.AskOne(q, &name, survey.WithValidator(survey.Required)); err != nil {
		fmt.Println(err)
		return 1
	}

	t := NewFromConfigItem(m[name])
	fmt.Println(t.Now())

	return 0
}

func myFilter(filterValue string, optValue string, optIndex int) bool {
	return strings.Contains(optValue, filterValue) && len(optValue) >= 3
}
