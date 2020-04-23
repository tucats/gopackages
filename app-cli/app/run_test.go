// Package app provides the top-level framework for CLI execution. This includes
// the Run() method to run the program, plus a number of action routines that can
// be invoked from the grammar or by a user action routine.
package app

import (
	"testing"

	"github.com/tucats/gopackages/app-cli/cli"
)

var testGrammar1 = []cli.Option{
	cli.Option{
		LongName:    "stuff",
		OptionType:  cli.BooleanType,
		Description: "stuff mart",
	},
}

func TestRun(t *testing.T) {
	type args struct {
		grammar []cli.Option
		args    []string
		appName string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Test",
			args: args{
				testGrammar1,
				[]string{"test-driver", "--booboo"},
				"testing: the test app",
			},
			wantErr: true,
		},
		{
			name: "Test",
			args: args{
				testGrammar1,
				[]string{"test-driver", "--stuff"},
				"testing",
			},
			wantErr: false,
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Run(tt.args.grammar, tt.args.args, tt.args.appName); (err != nil) != tt.wantErr {
				t.Errorf("Run() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
