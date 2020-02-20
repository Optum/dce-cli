package util

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/Optum/dce-cli/configs"
	"github.com/Optum/dce-cli/internal/constants"
	"github.com/mholt/archiver"
	"github.com/mitchellh/go-homedir"
	"gopkg.in/yaml.v2"
)

type FileSystemUtil struct {
	Config     *configs.Root
	ConfigFile string
}

func (u *FileSystemUtil) writeToYAMLFile(path string, _struct interface{}) error {
	var deferredErr error = nil
	_yaml, err := yaml.Marshal(_struct)
	if err != nil {
		return err
	}

	if !u.IsExistingFile(path) {
		err := os.MkdirAll(u.GetConfigDir(), os.FileMode(0700))
		if err != nil {
			return err
		}
		var file *os.File
		file, err = os.Create(path)
		if err != nil {
			return err
		}
		// #nosec
		defer file.Close()
	}

	err = ioutil.WriteFile(path, _yaml, 0644)
	if err != nil {
		return err
	}
	return deferredErr
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

// ReadFromFile returns the contents of a file as a string
// Care should be taken when using this function to prevent CWE-22 (https://cwe.mitre.org/data/definitions/22.html)
// i.e. ensure `path` comes from a trusted source.
func (u *FileSystemUtil) ReadFromFile(path string) string {
	/*
		#nosec CWE-22: added disclaimer to function docs
	*/
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
		err := os.Mkdir(destinationDir, os.ModeDir|os.FileMode(mode))
		if err != nil {
			log.Fatalln(err)
		}
	}

	originDir, err := os.Getwd()

	if err != nil {
		log.Fatalln(err)
	}

	err = os.Chdir(destinationDir)

	if err != nil {
		log.Fatalln(err)
	}

	return destinationDir, originDir
}

func (u *FileSystemUtil) ChToTmpDir() (string, string) {
	destinationDir := os.TempDir()
	originDir, err := os.Getwd()

	if err != nil {
		log.Fatalln(err)
	}

	err = os.Chdir(destinationDir)

	if err != nil {
		log.Fatalln(err)
	}

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

// OpenFileWriter opens or creates  a file in write-only mode. Data
// is appended to the file when writing.
// The file permissions are set to 0644, i.e. user-executable and user/group/other-readable.
func (u *FileSystemUtil) OpenFileWriter(path string) (*os.File, error) {
	/*
		#nosec CWE-276: called out the file permissions in function docs
	*/
	return os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
}

func (u *FileSystemUtil) GetCacheDir() string {
	return filepath.Join(u.GetConfigDir(), ".cache")
}

// GetArtifactsDir returns the cached artifacts dir, which by default is
// `~/.dce/.cache/dce/${DCE_VERSION}/`
func (u *FileSystemUtil) GetArtifactsDir() string {
	return filepath.Join(u.GetCacheDir(), constants.CommandShortName, constants.DCEBackendVersion)
}

// GetTerraformBinDir returns the dir in which the `terraform` bin is installed,
// which by default is `~/.dce/.cache/terraform/${TERRAFORM_VERSION}`
func (u *FileSystemUtil) GetTerraformBinDir() string {
	return filepath.Join(u.GetCacheDir(), "terraform", constants.TerraformBinVersion)
}

// GetLocalBackendDir returns the dir for the local terraform backend.
// By default, `~/.dce/.cache/module`
func (u *FileSystemUtil) GetLocalTFModuleDir() string {
	return filepath.Join(u.GetCacheDir(), "module")
}

// CreateConfigDirTree creates all the dirs in the dir specified by GetConfigDir(),
// including the dir itself.
func (u *FileSystemUtil) CreateConfigDirTree() error {
	dirs := []string{
		u.GetArtifactsDir(),
		u.GetTerraformBinDir(),
		u.GetLocalTFModuleDir(),
	}
	for _, dir := range dirs {
		err := os.MkdirAll(dir, os.FileMode(0700))
		if err != nil {
			return err
		}
	}
	return nil
}

// GetLogFile returns the full path of the log file for the deployment messages.
func (u *FileSystemUtil) GetLogFile() string {
	return filepath.Join(u.GetConfigDir(), "deploy.log")
}

// GetLocalBackendFile returns the full path of the local backend file.
func (u *FileSystemUtil) GetLocalMainTFFile() string {
	return filepath.Join(u.GetLocalTFModuleDir(), "main.tf")
}

// GetTerraformBin returns the full path of the terraform binary.
func (u *FileSystemUtil) GetTerraformBin() string {
	return filepath.Join(u.GetTerraformBinDir(), constants.TerraformBinName)
}

// GetTerraformStateFile returns the full path of the terraform state file
func (u *FileSystemUtil) GetTerraformStateFile() string {
	return filepath.Join(u.GetConfigDir(), "terraform.tfstate")
}
