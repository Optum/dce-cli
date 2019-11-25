package service

import (
	"errors"
	"fmt"
	"github.com/Optum/dce-cli/configs"
	observ "github.com/Optum/dce-cli/internal/observation"
	utl "github.com/Optum/dce-cli/internal/util"
	"net/url"
	"path"
	"time"
)

type AuthService struct {
	Config      *configs.Root
	Observation *observ.ObservationContainer
	Util        *utl.UtilContainer
}

func (s *AuthService) Authenticate() error {
	// Check that our API is configured properly
	if s.Config.API.Host == nil || s.Config.API.BasePath == nil {
		return errors.New("Unable to authenticate against DCE API: missing API configuration")
	}

	log.Println("Opening web browser. Please Login and copy/paste the provided token into this terminal.")
	// Wait a moment, so the user can see our message, and know what's going on
	time.Sleep(1 * time.Second)

	// Open the DCE API's /auth URL
	// this will use Cognito to redirect the user to
	// their configured IDP, then back to the /auth page,
	// which will display a "auth code" to the end-user.
	// The user will then need to copy the auth code
	// into their CLI prompt.
	authUrl := url.URL{
		Scheme:     "https",
		Host:       *s.Config.API.Host,
		Path:       path.Join(*s.Config.API.BasePath, "/auth"),
	}
	s.Util.OpenURL(authUrl.String())

	// Prompt for the auth code
	authCode := s.Util.PromptBasic(
		"Enter API Token: ", nil,
	)


	// Update the dce.yml config, with the token
	log.Printf("Saving API Token to %s", s.Util.GetConfigFile())
	s.Config.API.Token = authCode
	err := s.Util.WriteConfig()
	if err != nil {
		return fmt.Errorf("Failed to write to %s: %s",
			s.Util.GetConfigFile(), err)
	}

	return nil
}
