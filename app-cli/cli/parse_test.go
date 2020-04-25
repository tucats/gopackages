package cli

import (
	"errors"
	"testing"
)

func dummyAction(c *Context) error {
	return nil
}

func integerAction(c *Context) error {
	v, found := c.GetInteger("integer")
	if !found {
		return errors.New("No integer option found")
	}
	if v != 42 {
		return errors.New("Integer value not 42")
	}
	return nil
}

func stringAction(c *Context) error {
	v, found := c.GetString("string")
	if !found {
		return errors.New("No string option found")
	}
	if v != "foobar" {
		return errors.New("String value not foobar")
	}
	return nil
}

func booleanValueAction(c *Context) error {
	v := c.GetBool("boolean")
	if v != true {
		return errors.New("Boolean value not true")
	}
	return nil
}

func booleanAction(c *Context) error {
	v := c.GetBool("flag")
	if v != true {
		return errors.New("Boolean not present")
	}
	return nil
}

func TestContext_ParseGrammar(t *testing.T) {
	type fields struct {
		AppName                string
		MainProgram            string
		Description            string
		Command                string
		Grammar                []Option
		Args                   []string
		Parent                 *Context
		Parameters             []string
		ParameterCount         int
		ExpectedParameterCount int
		ParameterDescription   string
	}
	type args struct {
		args []string
	}

	var grammar1 = []Option{
		Option{
			LongName:   "sub1",
			OptionType: Subcommand,
			Action:     dummyAction,
			Value: []Option{
				Option{
					LongName:   "subopt1",
					OptionType: BooleanType,
				},
			},
		},
		Option{
			LongName:           "sub2",
			Aliases:            []string{"s2"},
			OptionType:         Subcommand,
			Action:             dummyAction,
			ParametersExpected: -3,
			Value: []Option{
				Option{
					LongName:   "subopt2",
					OptionType: BooleanType,
				},
			},
		},
		Option{
			LongName:           "sub3",
			Aliases:            []string{"s3"},
			OptionType:         Subcommand,
			Action:             dummyAction,
			ParametersExpected: 1,
			Value: []Option{
				Option{
					LongName:   "subopt2",
					OptionType: BooleanType,
				},
			},
		},
		Option{
			ShortName:   "a",
			LongName:    "alpha",
			OptionType:  BooleanType,
			Description: "alpha option",
			Action:      dummyAction,
		},
		Option{
			ShortName:   "i",
			LongName:    "integer",
			OptionType:  IntType,
			Description: "integer option",
			Action:      integerAction,
		},
		Option{
			ShortName:   "b",
			LongName:    "boolean",
			OptionType:  BooleanValueType,
			Description: "boolean value option",
			Action:      booleanValueAction,
		},
		Option{
			ShortName:   "f",
			LongName:    "flag",
			OptionType:  BooleanType,
			Description: "boolean option",
			Action:      booleanAction,
		},
		Option{
			ShortName:   "s",
			LongName:    "string",
			OptionType:  StringType,
			Description: "string option",
			Action:      stringAction,
		},
	}
	var fields1 = fields{
		AppName:     "unit test",
		MainProgram: "unit-test",
		Description: "Unit test stream",
		Command:     "unit-test",
		Grammar:     grammar1,
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:   "Short option name test",
			fields: fields1,
			args: args{
				args: []string{"-a"},
			},
			wantErr: false,
		},
		{
			name:   "Long option name test",
			fields: fields1,
			args: args{
				args: []string{"--alpha"},
			},
			wantErr: false,
		},
		{
			name:   "Unknown option test",
			fields: fields1,
			args: args{
				args: []string{"-x"},
			},
			wantErr: true,
		},
		{
			name:   "Option with = value test",
			fields: fields1,
			args: args{
				args: []string{"-i=42"},
			},
			wantErr: false,
		},

		// Integer options

		{
			name:   "Integer option test",
			fields: fields1,
			args: args{
				args: []string{"-i", "42"},
			},
			wantErr: false,
		},
		{
			name:   "Integer option missing value test",
			fields: fields1,
			args: args{
				args: []string{"-i"},
			},
			wantErr: true,
		},
		{
			name:   "Integer option illegal value test",
			fields: fields1,
			args: args{
				args: []string{"-i", "42f"},
			},
			wantErr: true,
		},

		// Boolean options

		{
			name:   "Boolean flag option test",
			fields: fields1,
			args: args{
				args: []string{"-f"},
			},
			wantErr: false,
		},
		{
			name:   "Boolean flag not present test",
			fields: fields1,
			args: args{
				args: []string{"-i", "3"},
			},
			wantErr: true,
		},
		{
			name:   "BooleanValue option test",
			fields: fields1,
			args: args{
				args: []string{"-b", "1"},
			},
			wantErr: false,
		},
		{
			name:   "BooleanValue option missing value test",
			fields: fields1,
			args: args{
				args: []string{"-b"},
			},
			wantErr: true,
		},
		{
			name:   "BooleanValue option illegal value test",
			fields: fields1,
			args: args{
				args: []string{"-i", "G"},
			},
			wantErr: true,
		},

		// String options

		{
			name:   "String option test",
			fields: fields1,
			args: args{
				args: []string{"-s", "foobar"},
			},
			wantErr: false,
		},
		{
			name:   "String option missing value test",
			fields: fields1,
			args: args{
				args: []string{"-s"},
			},
			wantErr: true,
		},
		{
			name:   "String option illegal value test",
			fields: fields1,
			args: args{
				args: []string{"-s", "foo"},
			},
			wantErr: true,
		},

		// Subcommands
		{
			name:   "Subcommand not found test",
			fields: fields1,
			args: args{
				args: []string{"sub99"},
			},
			wantErr: true,
		},
		{
			name:   "Subcommand found test",
			fields: fields1,
			args: args{
				args: []string{"sub1"},
			},
			wantErr: false,
		},
		{
			name:   "Subcommand with valid option test",
			fields: fields1,
			args: args{
				args: []string{"sub1", "--subopt1"},
			},
			wantErr: false,
		},
		{
			name:   "Subcommand with invalid option test",
			fields: fields1,
			args: args{
				args: []string{"sub1", "--subopt199"},
			},
			wantErr: true,
		},
		{
			name:   "Subcommand with alias name  test",
			fields: fields1,
			args: args{
				args: []string{"s2"},
			},
			wantErr: false,
		},
		{
			name:   "Subcommand with parameters test",
			fields: fields1,
			args: args{
				args: []string{"sub2", "parm1", "parm2"},
			},
			wantErr: false,
		},
		{
			name:   "Subcommand with too many parameters test",
			fields: fields1,
			args: args{
				args: []string{"sub2", "parm1", "parm2", "parm3", "parm4"},
			},
			wantErr: true,
		},
		{
			name:   "Subcommand with break before parameters test",
			fields: fields1,
			args: args{
				args: []string{"sub2", "--", "--parm1", "parm2", "parm3"},
			},
			wantErr: false,
		},

		// TODO: Add test cases.
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Context{
				AppName:                tt.fields.AppName,
				MainProgram:            tt.fields.MainProgram,
				Description:            tt.fields.Description,
				Command:                tt.fields.Command,
				Grammar:                tt.fields.Grammar,
				Args:                   tt.fields.Args,
				Parent:                 tt.fields.Parent,
				Parameters:             tt.fields.Parameters,
				ParameterCount:         tt.fields.ParameterCount,
				ExpectedParameterCount: tt.fields.ExpectedParameterCount,
				ParameterDescription:   tt.fields.ParameterDescription,
			}
			if err := c.ParseGrammar(tt.args.args); (err != nil) != tt.wantErr {
				t.Errorf("Context.ParseGrammar() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
