package config

import (
	"fmt"

	"github.com/AlecAivazis/survey"
)

var qs = []*survey.Question{
	{
		Name:     "name",
		Prompt:   &survey.Input{Message: "Enter service name"},
		Validate: survey.Required,
	},
	{
		Name:   "issuer",
		Prompt: &survey.Input{Message: "Enter issuer name (empty for no issuer)"},
	},
	{
		Name: "algorithm",
		Prompt: &survey.Select{
			Message: "Choose an hash algorithm:",
			Options: supportedAlgorithms,
			Default: defaultAlgorithm,
		},
	},
	{
		Name: "digits",
		Prompt: &survey.Input{
			Message: "How many digits should be in the TOTP code?",
			Default: "6",
		},
	},
	{
		Name:     "key",
		Prompt:   &survey.Password{Message: "Enter secret key"},
		Validate: survey.Required,
	},
}

func NewCommandNewItem(cfg *Config) *CommandNewItem {
	return &CommandNewItem{cfg}
}

type CommandNewItem struct {
	cfg *Config
}

func (c CommandNewItem) Execute(_ []string) int {
	item := &Item{}
	if err := survey.Ask(qs, item); err != nil {
		fmt.Println(err)
		return 1
	}

	d, err := c.cfg.parseAlgorithmFn(item.Algorithm)
	if err != nil {
		fmt.Println(err)
		return 1

	}

	item.digest = d

	if err := c.cfg.Add(item); err != nil {
		fmt.Println(err)
		return 1
	}

	if err := c.cfg.Write(); err != nil {
		fmt.Println(err)
		return 1
	}

	fmt.Printf("TOTP %s entity was successfully created\n", item.Name)

	return 0
}
