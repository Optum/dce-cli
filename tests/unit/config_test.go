package unit

import (
	"os"
	"testing"

	cfg "github.com/Optum/dce-cli/configs"
)

const Empty string = ""

func TestCoalesce(t *testing.T) {

	os.Setenv("TEST_CONFIG_VAL", "envval")

	tests := []struct {
		name   string
		arg    string
		config string
		envvar string
		def    string
		want   string
	}{
		{
			name:   "get default because everything else is empty",
			arg:    Empty,
			config: Empty,
			envvar: Empty,
			def:    "defaultval",
			want:   "defaultval",
		},
		{
			name:   "get default because env doesn't exist",
			arg:    Empty,
			config: Empty,
			envvar: "SOME_RANDOM_VAR_THAT_SHOULD_NOT_EXIST",
			def:    "defaultval",
			want:   "defaultval",
		},
		{
			name:   "get environment val because everything else is empty",
			arg:    Empty,
			config: Empty,
			envvar: "TEST_CONFIG_VAL",
			def:    "defaultval",
			want:   "envval",
		},
		{
			name:   "get environment val over config",
			arg:    Empty,
			config: "configval",
			envvar: "TEST_CONFIG_VAL",
			def:    "defaultval",
			want:   "envval",
		},
		{
			name:   "get config if env is empty",
			arg:    Empty,
			config: "configval",
			envvar: Empty,
			def:    "defaultval",
			want:   "configval",
		},
		{
			name:   "get arg first",
			arg:    "argval",
			config: "configval",
			envvar: "TEST_CONFIG_VAL",
			def:    "defaultval",
			want:   "argval",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := cfg.Coalesce(&tt.arg, &tt.config, &tt.envvar, &tt.def); *got != tt.want {
				t.Errorf("Coalesce() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCoalesce_WithNils(t *testing.T) {
	expected := "defaultval"
	if got := cfg.Coalesce(nil, nil, nil, &expected); *got != expected {
		t.Errorf("Coalesce() = %v, want %v", got, expected)
	}
}
