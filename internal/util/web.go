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
		log.Errorln("Error opening web browser", err)
		log.Infoln("Please copy and visit the following URL in your browser:", url)
	}
}

// Download will download the file at the given `url` and save it to the
// `localpath`
// Care should be taken to mitigate CWE-88 (https://cwe.mitre.org/data/definitions/88.html)
// by ensuring inputs comes from a trusted source.
func (w *WebUtil) Download(url string, localpath string) error {

	/*
		#nosec CWE-88: added disclaimer to function docs
	*/
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
	// #nosec
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}
