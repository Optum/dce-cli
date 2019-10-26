package util

import (
	"github.com/Optum/dce-cli/configs"
	observ "github.com/Optum/dce-cli/internal/observation"
	"github.com/manifoldco/promptui"
)

type PromptUtil struct {
	Config      *configs.Root
	Observation *observ.ObservationContainer
}

func (u *PromptUtil) PromptBasic(label string, validator func(input string) error) *string {
	prompt := promptui.Prompt{
		Label:    label,
		Validate: validator,
	}
	input, err := prompt.Run()
	if err != nil {
		Log.Fatalf("Prompt failed %v\n", err)
	}

	return &input
}

func (u *PromptUtil) PromptSelect(label string, items []string) *string {
	prompt := promptui.Select{
		Label: label,
		Items: items,
	}
	_, input, err := prompt.Run()
	if err != nil {
		Log.Fatalf("Prompt failed %v\n", err)
	}

	return &input
}
