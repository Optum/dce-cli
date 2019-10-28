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
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v2"
)

var testBinary string = "./testBinary"
var destinationDir string
var originDir string

func setUp() {
	log.Println("Building dce binary in temp dir")
	var err error
	destinationDir, err = ioutil.TempDir("", "")
	if err != nil {
		log.Fatalln(err)
	}

	out, _ := exec.Command("go", "build", "-o", testBinary, "../..").CombinedOutput()
	log.Println(string(out))

	out, _ = exec.Command("mv", testBinary, destinationDir).CombinedOutput()
	log.Println(string(out))

	log.Println("Moving to temp dir: " + destinationDir)
	originDir, err = os.Getwd()
	if err != nil {
		log.Fatalln(err)
	}
	os.Chdir(destinationDir)
}

func tearDown() {
	os.RemoveAll(destinationDir)
	os.Chdir(originDir)
}

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
	config.System.MasterAccount.Credentials.AwsAccessKeyID = &doesntMatter
	config.System.MasterAccount.Credentials.AwsSecretAccessKey = &doesntMatter
	config.Region = &usEast1
	config.API.BaseURL = &doesntMatter
	config.API.Credentials.AwsAccessKeyID = &doesntMatter
	config.API.Credentials.AwsSecretAccessKey = &doesntMatter
	config.API.Credentials.AwsSessionToken = &doesntMatter
	config.GithubToken = &doesntMatter

	input.types(*config.System.Auth.LoginURL)
	input.types(*config.System.MasterAccount.Credentials.AwsAccessKeyID)
	input.types(*config.System.MasterAccount.Credentials.AwsSecretAccessKey)
	input.pressesEnter()
	input.types(*config.API.BaseURL)
	input.types(*config.API.Credentials.AwsAccessKeyID)
	input.types(*config.API.Credentials.AwsSecretAccessKey)
	input.types(*config.API.Credentials.AwsSessionToken)
	input.types(*config.GithubToken)
	input.types("yes")

	stdin.Close()

	err = cmd.Wait()
	if err != nil {
		t.Log(err)
	}
	t.Run("THEN config is written to file at specified path", func(t *testing.T) {
		fullConfigPath := filepath.Join(destinationDir, configFileName)
		require.FileExists(t, fullConfigPath, "Config file not found")
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

			require.Equal(t, expectedJSON, actualJSON)
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
	time.Sleep(500 * time.Millisecond)
}
func (i *user) pressesEnter() {
	i.test.Helper()
	_, err := i.stdin.Write([]byte(string(promptui.KeyEnter)))
	if err != nil {
		i.test.Log(err)
	}
	time.Sleep(500 * time.Millisecond)
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
