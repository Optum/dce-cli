package functional

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"testing"
	"time"

	"path/filepath"

	"github.com/Optum/dce-cli/configs"
	"github.com/manifoldco/promptui"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
)

var configFileName string

// Reused test steps
var dceInitWritesToConfig = func(t *testing.T) {
	cmd := exec.Command(testBinary, "init", "--config", configFileName)

	stdin, err := cmd.StdinPipe()
	if err != nil {
		t.Log(err)
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err = cmd.Start(); err != nil {
		t.Log(err)
	}

	input := user{
		stdin: stdin,
		test:  t,
	}
	doesntMatter := "doesnt matter"
	usEast1 := "us-east-1"
	config := configs.Root{}
	config.System.Auth.LoginURL = &doesntMatter
	config.Region = &usEast1
	config.API.Host = &doesntMatter
	config.API.BasePath = &doesntMatter
	config.GithubToken = &doesntMatter

	input.types(*config.System.Auth.LoginURL)
	input.pressesEnter()
	input.types(*config.API.Host)
	input.types(*config.API.BasePath)
	input.types(*config.GithubToken)
	input.types("yes")

	stdin.Close()

	err = cmd.Wait()
	if err != nil {
		t.Log(err)
	}
	t.Run("THEN config is written to file at specified path", func(t *testing.T) {
		fullConfigPath := filepath.Join(destinationDir, configFileName)
		assert.FileExists(t, fullConfigPath, "Config file not found")
		t.Run("AND it matches the configuration provided by the user", func(t *testing.T) {
			actualStruct := configs.Root{}
			actualConfigFile, _ := ioutil.ReadFile(fullConfigPath)
			yaml.Unmarshal(actualConfigFile, &actualStruct)
			var actualBytes []byte
			actualBytes, err = json.Marshal(actualStruct)
			if err != nil {
				t.Log(err)
			}
			actualJSON := string(actualBytes)

			var exptectedBytes []byte
			exptectedBytes, err = json.Marshal(&config)
			expectedJSON := string(exptectedBytes)
			if err != nil {
				t.Log(err)
			}

			assert.Equal(t, expectedJSON, actualJSON)
		})
	})
}

func TestInitializationHappyPath(t *testing.T) {
	t.Run("GIVEN custom config path flag is provided", func(t *testing.T) {
		configFileName = "testConfig.yaml"
		t.Run("AND config file does not exists", func(t *testing.T) {
			setUp()
			t.Run("WHEN user types 'dce init' with --config flag AND enters config values", dceInitWritesToConfig)
			tearDown()
		})
		t.Run("AND config file exists", func(t *testing.T) {
			setUp()
			createFile(t, configFileName)
			t.Run("WHEN user types 'dce init' with --config flag AND enters config values", dceInitWritesToConfig)
			tearDown()
		})
	})
}

type user struct {
	stdin io.WriteCloser
	test  *testing.T
}

func (i *user) types(input string) {
	i.test.Helper()
	_, err := i.stdin.Write([]byte(input + string(promptui.KeyEnter)))
	if err != nil {
		i.test.Log(err)
	}
	time.Sleep(800 * time.Millisecond)
}

func (i *user) pressesEnter() {
	i.test.Helper()
	_, err := i.stdin.Write([]byte(string(promptui.KeyEnter)))
	if err != nil {
		i.test.Log(err)
	}
	time.Sleep(800 * time.Millisecond)
}

func createFile(t *testing.T, path string) {
	t.Helper()
	var file *os.File
	var err error
	file, err = os.Create(path)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	defer file.Close()
}
