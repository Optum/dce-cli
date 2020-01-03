package util

import (
	observ "github.com/Optum/dce-cli/internal/observation"
	"github.com/pkg/browser"
)

type WebUtil struct {
	Observation *observ.ObservationContainer
}

func (w *WebUtil) OpenURL(url string) {
	if err := browser.OpenURL(url); err != nil {
		log.Fatalln("Error opening web browser", err)
	}
}

func (w *WebUtil) Download(url string, localpath string) error {
	return nil
}
