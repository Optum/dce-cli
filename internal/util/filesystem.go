package util

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/Optum/dce-cli/configs"
	"github.com/mholt/archiver"
	"github.com/mitchellh/go-homedir"
	"gopkg.in/yaml.v2"
)

type FileSystemUtil struct {
	Config     *configs.Root
	ConfigFile string
}

func (u *FileSystemUtil) writeToYAMLFile(path string, _struct interface{}) error {

	_yaml, err := yaml.Marshal(_struct)
	if err != nil {
		return err
	}

	if !u.IsExistingFile(path) {
		var file *os.File
		file, err = os.Create(path)
		if err != nil {
			return err
		}
		defer file.Close()
	}

	err = ioutil.WriteFile(path, _yaml, 0644)
	if err != nil {
		return err
	}
	return nil
}

// WriteConfig writes the Config objects as YAML
// to the config file location (dce.yml)
func (u *FileSystemUtil) WriteConfig() error {
	return u.writeToYAMLFile(u.ConfigFile, u.Config)
}

// ReadInConfig loads the configuration from `dce.yml`
// and unmarshals it into the config object
func (u *FileSystemUtil) ReadInConfig() error {
	yamlStr, err := ioutil.ReadFile(u.ConfigFile)
	if err != nil {
		return err
	}

	return yaml.Unmarshal(yamlStr, u.Config)
}

func (u *FileSystemUtil) GetConfigFile() string {
	return u.ConfigFile
}

func (u *FileSystemUtil) GetHomeDir() string {
	homeDir, err := homedir.Dir()
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	return homeDir
}

func (u *FileSystemUtil) GetConfigDir() string {
	return filepath.Join(u.GetHomeDir(), ".dce")
}

func (u *FileSystemUtil) IsExistingFile(path string) bool {
	isExists := true
	if _, err := os.Stat(path); os.IsNotExist(err) {
		isExists = false
	}
	return isExists
}

func (u *FileSystemUtil) ReadFromFile(path string) string {
	contents, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	return string(contents)
}

func (u *FileSystemUtil) Unarchive(source string, destination string) {
	err := archiver.Unarchive(source, destination)
	if err != nil {
		log.Fatalln(err)
	}
}

func (u *FileSystemUtil) ChToConfigDir() (string, string) {
	destinationDir := u.GetConfigDir()

	mode := int(0700)
	if _, err := os.Stat(destinationDir); os.IsNotExist(err) {
		os.Mkdir(destinationDir, os.ModeDir|os.FileMode(mode))
	}

	err := os.Chdir(destinationDir)

	if err != nil {
		log.Fatalln(err)
	}
	originDir, err := os.Getwd()
	if err != nil {
		log.Fatalln(err)
	}
	os.Chdir(destinationDir)
	return destinationDir, originDir
}

func (u *FileSystemUtil) RemoveAll(path string) {
	err := os.RemoveAll(path)
	if err != nil {
		log.Fatalln(err)
	}
}

func (u *FileSystemUtil) Chdir(path string) {
	err := os.Chdir(path)
	if err != nil {
		log.Fatalln(err)
	}
}

func (u *FileSystemUtil) WriteFile(fileName string, data string) {
	err := ioutil.WriteFile(fileName, []byte(data), 0644)
	if err != nil {
		log.Fatalln(err)
	}
}

func (u *FileSystemUtil) ReadDir(path string) []os.FileInfo {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatalln(err)
	}
	return files
}

func (u *FileSystemUtil) OpenFileWriter(path string) (*os.File, error) {
	return os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
}
