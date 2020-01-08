package util

import (
	"io"
	"net/http"
	"os"

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

// Download will download the file at the given `url` and save it to the
// `localpath`
func (w *WebUtil) Download(url string, localpath string) error {

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(localpath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}
