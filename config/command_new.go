package config

import (
	"flag"
	"fmt"

	"github.com/AlecAivazis/survey/v2"
)

const CommandNewName string = "new"

var (
	shortQs = []*survey.Question{
		{
			Name:     "name",
			Prompt:   &survey.Input{Message: "Enter service name"},
			Validate: survey.Required,
		},
		{
			Name:     "key",
			Prompt:   &survey.Password{Message: "Enter secret key"},
			Validate: survey.Required,
		},
	}

	verboseQs = []*survey.Question{
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
)

func NewCommandNewItem(cfg *Config) *CommandNewItem {
	return &CommandNewItem{cfg}
}

type CommandNewItem struct {
	cfg *Config
}

func (c CommandNewItem) Execute(args []string) int {
	newCommand := flag.NewFlagSet(CommandNewName, flag.ExitOnError)
	helpFlag := newCommand.Bool("help", false, "Get this help")
	verboseFlag := newCommand.Bool("verbose", false, "Show verbose TOTP create input form")

	if err := newCommand.Parse(args); err != nil {
		fmt.Print(err)
		return 1
	}

	if *helpFlag {
		newCommand.PrintDefaults()
		return 0
	}

	qs := shortQs
	if *verboseFlag {
		qs = verboseQs
	}

	return c.ask(qs)
}

func (c CommandNewItem) ask(qs []*survey.Question) int {
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
