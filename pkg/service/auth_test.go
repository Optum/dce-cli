package service

import (
	"testing"

	"github.com/Optum/dce-cli/configs"
	"github.com/stretchr/testify/mock"
)

func TestAuthenticate(t *testing.T) {

	t.Run("GIVEN login url is specified in config", func(t *testing.T) {
		config := configs.Root{}
		defaultLoginURL := "default login url"
		config.System.Auth.LoginURL = &defaultLoginURL

		t.Run("WHEN Authenticate is called without a url override THEN open browser at default url", func(t *testing.T) {
			initMocks(config)
			mockWeber.On("OpenURL", mock.Anything)
			service.Authenticate("")
			mockWeber.AssertExpectations(t)
		})

		t.Run("WHEN Authenticate is called with a url override THEN open browser at override url", func(t *testing.T) {
			overrideLoginURL := "override login url"

			initMocks(config)
			mockWeber.On("OpenURL", overrideLoginURL)
			service.Authenticate(overrideLoginURL)
			mockWeber.AssertExpectations(t)
		})
	})
}
