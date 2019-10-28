package service

import (
	"github.com/Optum/dce-cli/configs"
	observ "github.com/Optum/dce-cli/internal/observation"
	utl "github.com/Optum/dce-cli/internal/util"
)

type AuthService struct {
	Config      *configs.Root
	Observation *observ.ObservationContainer
	Util        *utl.UtilContainer
}

func (s *AuthService) Authenticate(authUrl string) {
	log.Println("Opening web browser. Please Login and copy/paste the provided credentials into this terminal.")

	if authUrl == "" {
		authUrl = *s.Config.System.Auth.LoginURL
	}
	s.Util.OpenURL(authUrl)
}
