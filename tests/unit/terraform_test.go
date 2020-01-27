package unit

import (
	"reflect"
	"testing"

	util "github.com/Optum/dce-cli/internal/util"
)

func TestTerraformBinUtil_ParseOptions(t *testing.T) {

	tests := []struct {
		name    string
		str     string
		want    []string
		wantErr bool
	}{
		{
			name: "should return empty array with empty string",
			str:  "",
			want: []string{},
		},
		{
			name: "should parse normal space-delimited string",
			str:  "-backend-config=\"address=demo.consul.io\" -backend-config=\"path=example_app/terraform_state\"",
			want: []string{
				"-backend-config=\"address=demo.consul.io\"",
				"-backend-config=\"path=example_app/terraform_state\"",
			},
		},
		{
			name: "should parse normal multi-space-delimited string",
			str:  "-backend-config=\"address=demo.consul.io\"   -backend-config=\"path=example_app/terraform_state\"",
			want: []string{
				"-backend-config=\"address=demo.consul.io\"",
				"-backend-config=\"path=example_app/terraform_state\"",
			},
		},
		{
			name: "should parse normal tab-delimited string",
			str: `-backend-config="address=demo.consul.io"	-backend-config="path=example_app/terraform_state"`,
			want: []string{
				"-backend-config=\"address=demo.consul.io\"",
				"-backend-config=\"path=example_app/terraform_state\"",
			},
		},
		{
			name: "should parse normal newline-delimited string",
			str: `-backend-config="address=demo.consul.io"
			-backend-config="path=example_app/terraform_state"`,
			want: []string{
				"-backend-config=\"address=demo.consul.io\"",
				"-backend-config=\"path=example_app/terraform_state\"",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := util.ParseOptions(&tt.str)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseOptions() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseOptions() = %v, want %v", got, tt.want)
			}
		})
	}
}
