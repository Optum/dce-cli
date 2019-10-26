package util

import (
	"io/ioutil"
	"os"

	observ "github.com/Optum/dce-cli/internal/observation"
	"github.com/Optum/dce-cli/configs"
	"github.com/mitchellh/go-homedir"
	"gopkg.in/yaml.v2"
)

type FileSystemUtil struct {
	Config                *configs.Root
	Observation *observ.ObservationContainer
	DefaultConfigFileName string
}

func (u *FileSystemUtil) WriteToYAMLFile(path string, _struct interface{}) {
	_yaml, err := yaml.Marshal(_struct)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	err = ioutil.WriteFile(path, _yaml, 0644)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
}

func (u *FileSystemUtil) GetDefaultConfigFile() string {
	return u.GetHomeDir() + "/" + u.DefaultConfigFileName
}

func (u *FileSystemUtil) GetHomeDir() string {
	homeDir, err := homedir.Dir()
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	return homeDir
}

func (u *FileSystemUtil) IsExistingFile(path string) bool {
	isExists := true
	if _, err := os.Stat(path); os.IsNotExist(err) {
		isExists = false
	}
	return isExists
}
