package bytecode

import (
	"reflect"
	"testing"

	"github.com/tucats/gopackages/errors"
	"github.com/tucats/gopackages/expressions/data"
	"github.com/tucats/gopackages/expressions/symbols"
)

func Test_stopByteCode(t *testing.T) {
	ctx := &Context{running: true}

	if e := stopByteCode(ctx, nil); !e.(*errors.Error).Equal(errors.ErrStop) {
		t.Errorf("stopByteCode unexpected error %v", e)
	}

	if ctx.running {
		t.Errorf("stopByteCode did not turn off running flag")
	}
}

func Test_typeCast(t *testing.T) {
	name := "call to a type"
	tests := []struct {
		name string
		t    *data.Type
		v    interface{}
		want interface{}
		err  error
	}{
		{
			name: "cast int to string",
			t:    data.StringType,
			v:    55,
			want: "55",
		},
		{
			name: "cast bool to string",
			t:    data.StringType,
			v:    true,
			want: "true",
		},
	}

	for _, tt := range tests {
		ctx := &Context{
			stack:          make([]interface{}, 5),
			stackPointer:   0,
			running:        true,
			symbols:        symbols.NewSymbolTable("cast test"),
			programCounter: 1,
			bc: &ByteCode{
				instructions: make([]instruction, 5),
				nextAddress:  5,
			},
		}

		// Push the type on the stack that is to be used as the function pointer,
		// then the value to convert.
		_ = ctx.push(tt.t)
		_ = ctx.push(tt.v)

		err := callByteCode(ctx, 1)
		if err != nil {
			e1 := nilError
			e2 := nilError

			if tt.err != nil {
				e1 = tt.err.Error()
			}

			if err != nil {
				e2 = err.Error()
			}

			if e1 == e2 {
				return
			}

			t.Errorf("%s() error %v", name, err)
		} else if tt.err != nil {
			t.Errorf("%s() expected error not reported: %v", name, tt.err)
		}

		v, err := ctx.Pop()
		if err != nil {
			t.Errorf("%s() pop error: %v", name, err)
		}

		if !reflect.DeepEqual(v, tt.want) {
			t.Errorf("%s() got: %#v, want %#v", name, v, tt.want)
		}
	}
}

func Test_branchFalseByteCode(t *testing.T) {
	ctx := &Context{
		stack:          make([]interface{}, 5),
		stackPointer:   0,
		running:        true,
		programCounter: 1,
		bc: &ByteCode{
			instructions: make([]instruction, 5),
			nextAddress:  5,
		},
	}

	// Test if TOS is false
	_ = ctx.push(false)

	e := branchFalseByteCode(ctx, 2)
	if !errors.Nil(e) {
		t.Errorf("branchFalseByteCode unexpected error %v", e)
	}

	if ctx.programCounter != 2 {
		t.Errorf("branchFalseByteCode wrong program counter %v", ctx.programCounter)
	}

	// Test if TOS is true
	_ = ctx.push(true)

	e = branchFalseByteCode(ctx, 1)
	if !errors.Nil(e) {
		t.Errorf("branchFalseByteCode unexpected error %v", e)
	}

	if ctx.programCounter != 2 {
		t.Errorf("branchFalseByteCode wrong program counter %v", ctx.programCounter)
	}

	// Test if target is invalid
	_ = ctx.push(true)

	e = branchTrueByteCode(ctx, 20)
	if !e.(*errors.Error).Equal(errors.ErrInvalidBytecodeAddress) {
		t.Errorf("branchFalseByteCode unexpected error %v", e)
	}

	if ctx.programCounter != 2 {
		t.Errorf("branchFalseByteCode wrong program counter %v", ctx.programCounter)
	}
}

func Test_branchTrueByteCode(t *testing.T) {
	ctx := &Context{
		stack:          make([]interface{}, 5),
		stackPointer:   0,
		running:        true,
		programCounter: 1,
		bc: &ByteCode{
			instructions: make([]instruction, 5),
			nextAddress:  5,
		},
	}

	// Test if TOS is false
	_ = ctx.push(false)

	e := branchTrueByteCode(ctx, 2)
	if !errors.Nil(e) {
		t.Errorf("branchTrueByteCode unexpected error %v", e)
	}

	if ctx.programCounter != 1 {
		t.Errorf("branchTrueByteCode wrong program counter %v", ctx.programCounter)
	}

	// Test if TOS is true
	_ = ctx.push(true)

	e = branchTrueByteCode(ctx, 2)
	if !errors.Nil(e) {
		t.Errorf("branchTrueByteCode unexpected error %v", e)
	}

	if ctx.programCounter != 2 {
		t.Errorf("branchTrueByteCode wrong program counter %v", ctx.programCounter)
	}

	// Test if target is invalid
	_ = ctx.push(true)

	e = branchTrueByteCode(ctx, 20)
	if !e.(*errors.Error).Equal(errors.ErrInvalidBytecodeAddress) {
		t.Errorf("branchTrueByteCode unexpected error %v", e)
	}

	if ctx.programCounter != 2 {
		t.Errorf("branchTrueByteCode wrong program counter %v", ctx.programCounter)
	}
}
