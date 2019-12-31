package unit

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"testing"

	util "github.com/Optum/dce-cli/internal/util"
	"github.com/stretchr/testify/assert"
)

func TestMainTFTemplate_MissingRequiredValues(t *testing.T) {

	var actual bytes.Buffer
	var err error

	mockFS := &mockFileSystemer
	mockFS.On("GetConfigDir").Return("/Users/jexmple/.dce")
	tf := util.NewMainTFTemplate(mockFS)
	tf.Version = "v0.23.0"
	tf.TFWorkspaceDir = ""

	err = tf.AddVariable("namespace", "string", "dcecliut")
	assert.Nil(t, err, "should have been able to add valid namespace")
	// Try to write without the required TFWorkspaceDir
	err = tf.Write(&actual)
	assert.NotNil(t, err)
	assert.Equal(t, "non-zero length value required for workspace dir", err.Error())

	tf.TFWorkspaceDir = "/Users/jexmple/.dce/tf-workspace"
	tf.LocalTFStateFilePath = ""
	// Set it, but now try to write without the required LocalTFStateFilePath
	err = tf.Write(&actual)
	assert.NotNil(t, err)
	assert.Equal(t, "non-zero length value required for local tf state file path", err.Error())

}

// TestMainTFTemplate_Write tests the `Write` happy path
func TestMainTFTemplate_Write(t *testing.T) {

	var actual bytes.Buffer

	expected, err := ioutil.ReadFile("examples/maintf-basic.example")

	assert.Nil(t, err, "should have been able to read from file")

	mockFS := &mockFileSystemer
	mockFS.On("GetConfigDir").Return("/Users/jexmple/.dce")
	tf := util.NewMainTFTemplate(mockFS)
	tf.Version = "v0.23.0"
	tf.LocalTFStateFilePath = "/Users/jexmple/.dce/terraform.tfstate"
	tf.TFWorkspaceDir = "/Users/jexmple/.dce/tf-workspace"

	err = tf.AddVariable("namespace", "string", "dcecliut")
	assert.Nil(t, err, "should have been able to add valid namespace")

	err = tf.AddVariable("budget_notification_from_email", "string", "noreply@example.com")
	assert.Nil(t, err, "should have been able to add valid budget_notification_from_email")

	err = tf.Write(&actual)
	assert.Nil(t, err)
	assert.Equal(t, expected, actual.Bytes(), fmt.Sprintf("templates should have same content: \"%s\" vs. \"%s\"", string(expected), actual.String()))

}
