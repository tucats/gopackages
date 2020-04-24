// Package app provides the top-level framework for CLI execution. This includes
// the Run() method to run the program, plus a number of action routines that can
// be invoked from the grammar or by a user action routine.
package app

import (
	"errors"
	"fmt"
	"testing"

	"github.com/tucats/gopackages/app-cli/cli"
)

var testGrammar1 = []cli.Option{
	cli.Option{
		LongName:    "stuff",
		OptionType:  cli.BooleanType,
		Description: "stuff mart",
	},
	cli.Option{
		LongName:    "sub1",
		OptionType:  cli.Subcommand,
		Description: "sub1 subcommand",
	},
	cli.Option{
		LongName:    "sub2",
		OptionType:  cli.Subcommand,
		Description: "sub2 subcommand has options",
		Value: []cli.Option{
			cli.Option{
				ShortName:   "x",
				LongName:    "explode",
				Description: "Make something blow up",
				OptionType:  cli.StringType,
				Action:      testAction1,
			},
		},
	},
}

func testAction1(c *cli.Context) error {
	v, _ := c.GetString("explode")
	fmt.Printf("Found the option value %s\n", v)
	if v != "bob" {
		return errors.New("Invalid explode name: " + v)
	}
	return nil
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
			name: "Invalid global option",
			args: args{
				testGrammar1,
				[]string{"test-driver", "-d", "--booboo"},
				"testing: the test app",
			},
			wantErr: true,
		},
		{
			name: "Valid global option",
			args: args{
				testGrammar1,
				[]string{"test-driver", "-d", "--stuff"},
				"testing: the test app",
			},
			wantErr: false,
		},
		{
			name: "Valid subcommand",
			args: args{
				testGrammar1,
				[]string{"test-driver", "-d", "sub1"},
				"testing: the test app",
			},
			wantErr: false,
		},
		{
			name: "Invalid subcommand",
			args: args{
				testGrammar1,
				[]string{"test-driver", "-d", "subzero"},
				"testing: the test app",
			},
			wantErr: true,
		},
		{
			name: "Valid subcommand with valid short option name",
			args: args{
				testGrammar1,
				[]string{"test-driver", "-d", "sub2", "-x", "bob"},
				"testing: the test app",
			},
			wantErr: false,
		},
		{
			name: "Valid subcommand with valid long option",
			args: args{
				testGrammar1,
				[]string{"test-driver", "-d", "sub2", "--explode", "bob"},
				"testing: the test app",
			},
			wantErr: false,
		},
		{
			name: "Valid subcommand with valid option but invalid value",
			args: args{
				testGrammar1,
				[]string{"test-driver", "-d", "sub2", "-x", "bob2"},
				"testing: the test app",
			},
			wantErr: true,
		},
		{
			name: "Valid subcommand with valid option but missing value",
			args: args{
				testGrammar1,
				[]string{"test-driver", "-d", "sub2", "-x"},
				"testing: the test app",
			},
			wantErr: true,
		},
		// TODO: Add test cases.
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			SetCopyright("(c) 2020 Tom Cole. All rights reserved.")
			SetVersion([]int{1, 1, 1})

			if err := Run(tt.args.grammar, tt.args.args, tt.args.appName); (err != nil) != tt.wantErr {
				t.Errorf("Run() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
