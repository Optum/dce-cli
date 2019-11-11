package functional

import (
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

var destinationDir string
var originDir string

func setUp() {
	log.Println("Building dce binary in temp dir")
	var err error
	destinationDir, err = ioutil.TempDir("", "")
	if err != nil {
		log.Fatalln(err)
	}

	out, _ := exec.Command("go", "build", "-o", filepath.Join(destinationDir, testBinary), "../..").CombinedOutput()
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
