package service

import (
	"github.com/Optum/dce-cli/configs"
	observ "github.com/Optum/dce-cli/internal/observation"
	"github.com/pkg/browser"
)

type AuthService struct {
	Config      *configs.Root
	Observation *observ.ObservationContainer
}

func (s *AuthService) Authenticate(authUrl string) {
	Log.Println("Opening web browser. Please u.Observationin and copy/paste the provided credentials into this terminal.")

	if authUrl == "" {
		authUrl = *s.Config.System.Auth.LoginURL
	}
	browser.OpenURL(authUrl)
}
