package util

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"text/template"
)

// tfMainTemplate is the template for the main.tf file that is generated
// and written to the DCE home directory (`~/.dce`). The file contains
// variables with default values. When building out the template this
// way, I had considered just directly setting the variables for the
// module but decided instead to include them as referenced variables
// with defaults so that `terraform apply` could be easier used again on
// the file directly.
const tfMainTemplate string = `terraform {
{{if .LocalBackend }}
	backend "local" {
		path="{{.LocalTFStateFilePath}}"
		workspace_dir="{{.TFWorkspaceDir}}"
	}
{{- /* Put other backend configurations here, like when they come from the YAML file... */ -}}
{{end}}
}
{{range .TFVars}}
variable "{{.Name}}" {
	type = {{.Type}}
	default = "{{.Value}}"
}
{{end}}
module "dce" {
	source="{{.ModuleSource}}"
{{range .TFVars}}
	{{.Name}} = var.{{.Name}}
{{- end}}
}
{{/* This is just hard-coded because the code depends on this */}}
output "artifacts_bucket_name" {
	description = "S3 bucket for artifacts like AWS Lambda code"
	value = module.dce.artifacts_bucket_name
}

output "api_url" {
	description = "URL of DCE API"
	value = module.dce.api_url
}
`

// TFVar represents a variable that is in the
type TFVar struct {
	Name  string
	Type  string
	Value string
}

// MainTFTemplate is the template for writing the main.tf file
type MainTFTemplate struct {
	TFVars               []TFVar
	LocalBackend         bool
	LocalTFStateFilePath string
	TFWorkspaceDir       string
	// DCE terraform module source location
	ModuleSource string
}

// AddVariable adds a variable with the given `name`, variable type (`vartype`),
// and default value (`vardefault`) to the template
func (t *MainTFTemplate) AddVariable(name string, vartype string, val string) error {
	if len(name) == 0 {
		return fmt.Errorf("non-zero length value required for name")
	}

	if len(vartype) == 0 {
		return fmt.Errorf("non-zero length value required for vartype")
	}

	if len(val) == 0 {
		return fmt.Errorf("non-zero length value required for val")
	}

	t.TFVars = append(t.TFVars, TFVar{
		Name:  name,
		Type:  vartype,
		Value: val,
	})

	return nil
}

func (t *MainTFTemplate) SetModuleSource(source string) {
	t.ModuleSource = source
}

// Write writes the template to the given writer
func (t *MainTFTemplate) Write(w io.Writer) error {

	if t.LocalBackend {
		if len(t.TFWorkspaceDir) == 0 {
			return fmt.Errorf("non-zero length value required for workspace dir")
		}
		if len(t.LocalTFStateFilePath) == 0 {
			return fmt.Errorf("non-zero length value required for local tf state file path")
		}
	}

	tplate := template.Must(template.New("tfmain").Parse(tfMainTemplate))
	err := tplate.Execute(w, t)
	return err
}

// NewMainTFTemplate creates a new instance of the MainTFTemplate
func NewMainTFTemplate(fs FileSystemer) *MainTFTemplate {

	tfWorkDir := filepath.Join(fs.GetCacheDir(), "tf-workspace")
	if _, err := os.Stat(tfWorkDir); os.IsNotExist(err) {
		/*
			#nosec CWE-703: failures accounted for in error returned by Write function. TODO: refactor this
		 */
		_ = os.Mkdir(tfWorkDir, os.ModeDir|os.FileMode(int(0700)))
	}

	tfStateFilePath := fs.GetTerraformStateFile()

	tf := &MainTFTemplate{
		LocalBackend:         true,
		LocalTFStateFilePath: tfStateFilePath,
		TFWorkspaceDir:       tfWorkDir,
	}
	return tf
}
