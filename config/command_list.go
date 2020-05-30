package config

import (
	"fmt"
	"strings"

	"github.com/AlecAivazis/survey"
)

func NewCommandList(cfg *Config) *CommandList {
	return &CommandList{cfg}
}

type CommandList struct {
	cfg *Config
}

func (c CommandList) Execute(_ []string) int {
	m := make(map[string]*Item)
	options := make([]string, 0, len(c.cfg.Items))
	for _, i := range c.cfg.Items {
		m[i.Name] = i
		options = append(options, i.Name)
	}

	q := &survey.Select{
		Message: "Choose a color:",
		Options: options,
		Filter:  myFilter,
	}

	var name string
	if err := survey.AskOne(q, &name, survey.WithValidator(survey.Required)); err != nil {
		fmt.Println(err)
		return 1
	}

	fmt.Printf("%+v\n", m[name])

	return 0
}

func myFilter(filterValue string, optValue string, optIndex int) bool {
	return strings.Contains(optValue, filterValue) && len(optValue) >= 3
}
