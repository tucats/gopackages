// Package app provides the top-level framework for CLI execution. This includes
// the Run() method to run the program, plus a number of action routines that can
// be invoked from the grammar or by a user action routine.
package app

import (
	"errors"
	"fmt"
	"strconv"
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
		Action:      testAction0,
	},
	cli.Option{
		LongName:    "sub2",
		OptionType:  cli.Subcommand,
		Description: "sub2 subcommand has options",
		Action:      testAction0,
		Value: []cli.Option{
			cli.Option{
				ShortName:   "x",
				LongName:    "explode",
				Description: "Make something blow up",
				OptionType:  cli.StringType,
				Action:      testAction1,
			},
			cli.Option{
				LongName:    "count",
				Description: "Count of things to blow up",
				OptionType:  cli.IntType,
				Action:      testAction2,
			},
		},
	},
}

func testAction0(c *cli.Context) error {
	return nil
}

func testAction1(c *cli.Context) error {
	v, _ := c.GetString("explode")
	fmt.Printf("Found the option value %s\n", v)
	if v != "bob" {
		return errors.New("Invalid explode name: " + v)
	}
	return nil
}

func testAction2(c *cli.Context) error {
	v, _ := c.GetInteger("count")
	fmt.Printf("Found the option value %v\n", v)
	if v != 42 {
		return errors.New("Invalid count: " + strconv.Itoa(v))
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
		{
			name: "Valid subcommand with invalid int valid option ",
			args: args{
				testGrammar1,
				[]string{"test-driver", "-d", "sub2", "--count", "42F"},
				"testing: the test app",
			},
			wantErr: true,
		},
		{
			name: "Valid subcommand with valid option ",
			args: args{
				testGrammar1,
				[]string{"test-driver", "-d", "sub2", "--count", "42"},
				"testing: the test app",
			},
			wantErr: false,
		},
		{
			name: "Valid subcommand with valid option with wrong value",
			args: args{
				testGrammar1,
				[]string{"test-driver", "-d", "sub2", "--count", "43"},
				"testing: the test app",
			},
			wantErr: true,
		},
		// TODO: Add test cases.
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			app := New("tt.args.appName")
			app.SetCopyright("(c) 2020 Tom Cole. All rights reserved.")
			app.SetVersion(1, 1, 0)

			if err := app.Run(tt.args.grammar, tt.args.args); (err != nil) != tt.wantErr {
				t.Errorf("Run() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
