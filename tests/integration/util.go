package integration

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"os"
	"testing"
)

type file struct{ Name, Body string }

// Adapted from https://golang.org/pkg/archive/zip/#example_Writer
func zipFiles(t *testing.T, files []file) []byte {
	// Create a buffer to write our archive to.
	buf := new(bytes.Buffer)

	// Create a new zip archive.
	w := zip.NewWriter(buf)

	// Add some files to the archive.
	for _, file := range files {
		f, err := w.Create(file.Name)
		require.Nil(t, err)
		_, err = f.Write([]byte(file.Body))
		require.Nil(t, err)
	}

	// Make sure to check the error on Close.
	err := w.Close()
	require.Nil(t, err)

	zipBytes, err := ioutil.ReadAll(buf)
	require.Nil(t, err)

	return zipBytes
}


// mockEnvVar Sets an env var
// returns a function to revert the env var back to its previous value
func mockEnvVar(key, val string) func() {
	prevVal, ok := os.LookupEnv(key)
	_ = os.Setenv(key, val)

	return func() {
		if ok {
			_ = os.Setenv(key, prevVal)
		} else {
			_ = os.Unsetenv(key)
		}
	}
}


func stringp(str string) *string {
	return &str
}

func boolp(b bool) *bool {
	return &b
}

// copyStructVals copies all values from the source stuct
// into the target struct
// Structs must be JSON encode-able
func copyStructVals(t *testing.T, source interface{}, target interface{}) {
	// We're using JSON as a means to copy values,
	// which has some limitations.
	// This could probably be done more reliably, but
	// with a lot more complex code, using reflection.
	sourceJSON, err := json.Marshal(&source)
	require.Nil(t, err)
	err = json.Unmarshal(sourceJSON, target)
	require.Nil(t, err)
}
