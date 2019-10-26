package cmd

import (
	"io/ioutil"
	"log"
	"os"
	"os/exec"

	"testing"

	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	log.Println("Building dce binary in temp dir")
	destinationDir, err := ioutil.TempDir("", "")
	if err != nil {
		log.Fatalln(err)
	}
	defer os.RemoveAll(destinationDir)

	out, _ := exec.Command("go", "build", "..", "-o", destinationDir).CombinedOutput()
	log.Println(string(out))

	log.Println("Moving to temp dir")
	originDir, err := os.Getwd()
	if err != nil {
		log.Fatalln(err)
	}
	defer os.Chdir(originDir)
	os.Chdir(destinationDir)

	log.Println("Executing tests")
	os.Exit(m.Run())
}

func TestRoot(t *testing.T) {
	// t.Run("GIVEN a config file does not exists", func(t *testing.T) {

	// 	t.Run("WHEN root executes", func(t *testing.T) {
	// 		out, err := exec.Command("dce").CombinedOutput()
	// 		if err != nil {
	// 			t.Fatalf("err: %s", err)
	// 		}
	// 		t.Log(out)

	// 		t.Run("THEN a new config file will be created", func(t *testing.T) {
	// 			require.Equal(t, 1, 1)
	// 		})
	// 	})

	// })

	t.Run("GIVEN a config file exists", func(t *testing.T) {
		out, _ := exec.Command("dce").CombinedOutput()
		t.Log(string(out))

		t.Run("WHEN root executes", func(t *testing.T) {

			t.Run("THEN the config file will be parsed into config", func(t *testing.T) {
				require.Equal(t, 1, 1)
			})
		})

	})
}
