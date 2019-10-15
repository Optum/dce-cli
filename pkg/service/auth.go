package service

import (
	"log"

	"github.com/Optum/dce-cli/configs"
	"github.com/pkg/browser"
)

type AuthService struct {
	Config *configs.Root
}

func (s *AuthService) Authenticate(authUrl string) {
	log.Println("Opening web browser. Please login and copy/paste the provided credentials into this terminal.")

	if authUrl == "" {
		authUrl = *s.Config.System.Auth.LoginURL
	}
	browser.OpenURL(authUrl)
}
