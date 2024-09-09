package util

import (
	"fmt"
	"github.com/Optum/dce-cli/configs"
	observ "github.com/Optum/dce-cli/internal/observation"
	"github.com/chzyer/readline"
)

type PromptUtil struct {
	Config      *configs.Root
	Observation *observ.ObservationContainer
}

func (u *PromptUtil) PromptBasic(label string, validator func(input string) error) *string {
	rl, err := readline.NewEx(&readline.Config{
		Prompt:                 fmt.Sprint(label, " "),
		DisableAutoSaveHistory: true,
	})
	if err != nil {
		log.Fatalln(err)
	}
	defer rl.Close() //nolint:errcheck

	input, err := rl.Readline()
	if err != nil {
		log.Fatalln(err)
	}
	if validator != nil {
		err = validator(input)
		if err != nil {
			log.Fatalln(err)
		}
	}

	return &input
}
