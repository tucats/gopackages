package bytecode

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/tucats/gopackages/errors"
	"github.com/tucats/gopackages/expressions/symbols"
)

const nilError = "nil error"

func Test_negateByteCode(t *testing.T) {
	target := negateByteCode

	tests := []struct {
		name  string
		arg   any
		stack []any
		want  any
		err   error
	}{
		{
			name:  "stack underflow",
			arg:   nil,
			stack: []any{},
			err:   errors.ErrStackUnderflow,
		},
		{
			name:  "boolean NOT stack underflow",
			arg:   true,
			stack: []any{},
			err:   errors.ErrStackUnderflow,
		},
		{
			name:  "unsupported nil value",
			arg:   nil,
			stack: []any{nil},
			want:  -5,
			err:   errors.ErrInvalidType.Context("nil"),
		},
		{
			name:  "negate nil boolean value",
			arg:   true,
			stack: []any{nil},
			want:  true,
			err:   nil,
		},
		{
			name:  "NOT  integer",
			arg:   true,
			stack: []any{3},
			want:  false,
		},
		{
			name:  "NOT  int64",
			arg:   true,
			stack: []any{int64(3)},
			want:  false,
		},
		{
			name:  "NOT boolean",
			arg:   true,
			stack: []any{false},
			want:  true,
		},
		{
			name:  "NOT float32",
			arg:   true,
			stack: []any{float32(3.5)},
			want:  false,
		},
		{
			name:  "NOT  float64",
			arg:   true,
			stack: []any{float64(0)},
			want:  true,
		},
		{
			name:  "NOT  string",
			arg:   true,
			stack: []any{"invalid"},
			want:  false,
			err:   errors.ErrInvalidType,
		},
		{
			name:  "negate negative integer",
			arg:   nil,
			stack: []any{-3},
			want:  3,
		},
		{
			name:  "negate int64 integer",
			arg:   nil,
			stack: []any{int64(-3)},
			want:  int64(3),
		},
		{
			name:  "negate boolean",
			arg:   nil,
			stack: []any{false},
			want:  true,
		},
		{
			name:  "negate string",
			arg:   nil,
			stack: []any{"test"},
			want:  "tset",
		},
		{
			name:  "negate positive float32",
			arg:   nil,
			stack: []any{float32(6.6)},
			want:  float32(-6.6),
		},
		{
			name:  "negate negative float64",
			arg:   nil,
			stack: []any{float64(-6.6)},
			want:  float64(6.6),
		},
		{
			name:  "negate positive byte",
			arg:   nil,
			stack: []any{byte(2)},
			want:  byte(254),
		},
		{
			name:  "negate negative int32",
			arg:   nil,
			stack: []any{int32(-12)},
			want:  int32(12),
		},
	}

	for _, tt := range tests {
		syms := symbols.NewSymbolTable("testing")
		bc := ByteCode{}

		c := NewContext(syms, &bc)

		for _, item := range tt.stack {
			_ = c.push(item)
		}

		t.Run(tt.name, func(t *testing.T) {
			err := target(c, tt.arg)
			if err != nil {
				e1 := nilError

				if tt.err != nil {
					e1 = tt.err.Error()
				}

				e2 := err.Error()

				if e1 == e2 {
					return
				}

				t.Errorf("negateByteCode() error %v", err)
			} else if tt.err != nil {
				t.Errorf("negateByteCode() expected error not reported: %v", tt.err)
			}

			v, err := c.Pop()

			if err != nil {
				t.Errorf("negateByteCode() stack error %v", err)
			}

			if !reflect.DeepEqual(v, tt.want) {
				t.Errorf("negateByteCode() got %v, want %v", v, tt.want)
			}
		})
	}
}

func Test_addByteCode(t *testing.T) {
	target := addByteCode
	name := "addByteCode"

	tests := []struct {
		name  string
		arg   any
		stack []any
		want  any
		err   error
		debug bool
	}{
		{
			name:  "add string to error",
			arg:   nil,
			stack: []any{errors.ErrAssert, "-thing"},
			want:  "assert-thing",
		},
		{
			name:  "add with first nil",
			arg:   nil,
			stack: []any{2, nil},
			want:  7,
			err:   errors.ErrInvalidType.Context("nil"),
		},
		{
			name:  "add with second nil",
			arg:   nil,
			stack: []any{nil, 5},
			want:  7,
			err:   errors.ErrInvalidType.Context("nil"),
		},
		{
			name:  "add integers",
			arg:   nil,
			stack: []any{2, 5},
			want:  7,
		},
		{
			name:  "AND mixed boolean",
			arg:   nil,
			stack: []any{true, false},
			want:  false,
		},
		{
			name:  "AND same boolean",
			arg:   nil,
			stack: []any{true, true},
			want:  true,
		},
		{
			name:  "add strings",
			arg:   nil,
			stack: []any{"test", "plan"},
			want:  "testplan",
		},
		{
			name:  "add float32",
			arg:   nil,
			stack: []any{float32(1.0), float32(6.6)},
			want:  float32(7.6),
		},
		{
			name:  "add int32 to byte",
			arg:   nil,
			stack: []any{int(5), byte(2)},
			want:  int(7),
		},
		{
			name:  "add byte to int32",
			arg:   nil,
			stack: []any{byte(5), int32(2)},
			want:  int32(7),
		},
		{
			name:  "add int32 to int32",
			arg:   nil,
			stack: []any{int32(5), int32(2)},
			want:  int32(7),
		},
		{
			name:  "add with 1 args on stack",
			arg:   nil,
			stack: []any{int32(-12)},
			want:  int32(12),
			err:   errors.ErrStackUnderflow,
		},
		{
			name:  "add with 0 args on stack",
			arg:   nil,
			stack: []any{},
			want:  int32(12),
			err:   errors.ErrStackUnderflow,
		},
	}

	for _, tt := range tests {
		syms := symbols.NewSymbolTable("testing")
		bc := ByteCode{}

		c := NewContext(syms, &bc)

		for _, item := range tt.stack {
			_ = c.push(item)
		}

		t.Run(tt.name, func(t *testing.T) {
			if tt.debug {
				fmt.Println("DEBUG")
			}

			err := target(c, tt.arg)

			if err != nil {
				e1 := nilError

				if tt.err != nil {
					e1 = tt.err.Error()
				}

				e2 := err.Error()

				if e1 == e2 {
					return
				}

				t.Errorf("%s() error %v", name, err)
			} else if tt.err != nil {
				t.Errorf("%s() expected error not reported: %v", name, tt.err)
			}

			v, err := c.Pop()

			if err != nil {
				t.Errorf("%s() stack error %v", name, err)
			}

			if !reflect.DeepEqual(v, tt.want) {
				t.Errorf("%s() got %v, want %v", name, v, tt.want)
			}
		})
	}
}

func Test_andByteCode(t *testing.T) {
	target := andByteCode
	name := "andByteCode"

	tests := []struct {
		name  string
		arg   any
		stack []any
		want  any
		err   error
		debug bool
	}{
		{
			name:  "AND with no value",
			arg:   nil,
			stack: []any{},
			err:   errors.ErrStackUnderflow,
		},
		{
			name:  "AND with only 1 value",
			arg:   nil,
			stack: []any{true},
			err:   errors.ErrStackUnderflow,
		}, {
			name:  "AND string to error",
			arg:   nil,
			stack: []any{errors.ErrAssert, "-thing"},
			want:  false,
		},
		{
			name:  "AND empty string to error",
			arg:   nil,
			stack: []any{errors.ErrAssert, ""},
			want:  false,
		},
		{
			name:  "AND with first nil",
			arg:   nil,
			stack: []any{2, nil},
			want:  7,
			err:   errors.ErrInvalidType.Context("nil"),
		},
		{
			name:  "AND with second nil",
			arg:   nil,
			stack: []any{nil, 5},
			want:  7,
			err:   errors.ErrInvalidType.Context("nil"),
		},
		{
			name:  "AND integers",
			arg:   nil,
			stack: []any{2, 5},
			want:  true,
		},
		{
			name:  "AND mixed boolean",
			arg:   nil,
			stack: []any{true, false},
			want:  false,
		},
		{
			name:  "AND same boolean",
			arg:   nil,
			stack: []any{true, true},
			want:  true,
		},
		{
			name:  "AND strings",
			arg:   nil,
			stack: []any{"test", "plan"},
			want:  false,
		},
		{
			name:  "AND float32",
			arg:   nil,
			stack: []any{float32(1.0), float32(6.6)},
			want:  true,
		},
	}

	for _, tt := range tests {
		syms := symbols.NewSymbolTable("testing")
		bc := ByteCode{}

		c := NewContext(syms, &bc)

		for _, item := range tt.stack {
			_ = c.push(item)
		}

		t.Run(tt.name, func(t *testing.T) {
			if tt.debug {
				fmt.Println("DEBUG")
			}

			err := target(c, tt.arg)

			if err != nil {
				e1 := nilError

				if tt.err != nil {
					e1 = tt.err.Error()
				}

				e2 := err.Error()

				if e1 == e2 {
					return
				}

				t.Errorf("%s() error %v", name, err)
			} else if tt.err != nil {
				t.Errorf("%s() expected error not reported: %v", name, tt.err)
			}

			v, err := c.Pop()

			if err != nil {
				t.Errorf("%s() stack error %v", name, err)
			}

			if !reflect.DeepEqual(v, tt.want) {
				t.Errorf("%s() got %v, want %v", name, v, tt.want)
			}
		})
	}
}

func Test_orByteCode(t *testing.T) {
	target := orByteCode
	name := "orByteCode"

	tests := []struct {
		name  string
		arg   any
		stack []any
		want  any
		err   error
		debug bool
	}{
		{
			name:  "OR with no value",
			arg:   nil,
			stack: []any{},
			err:   errors.ErrStackUnderflow,
		},
		{
			name:  "OR with only 1 value",
			arg:   nil,
			stack: []any{true},
			err:   errors.ErrStackUnderflow,
		},
		{
			name:  "OR string to error",
			arg:   nil,
			stack: []any{errors.ErrAssert, "-thing"},
			want:  false,
		},
		{
			name:  "OR empty string to error",
			arg:   nil,
			stack: []any{errors.ErrAssert, ""},
			want:  false,
		},
		{
			name:  "OR with first nil",
			arg:   nil,
			stack: []any{2, nil},
			want:  7,
			err:   errors.ErrInvalidType.Context("nil"),
		},
		{
			name:  "OR with second nil",
			arg:   nil,
			stack: []any{nil, 5},
			want:  7,
			err:   errors.ErrInvalidType.Context("nil"),
		},
		{
			name:  "OR integers",
			arg:   nil,
			stack: []any{2, 5},
			want:  true,
		},
		{
			name:  "OR mixed boolean",
			arg:   nil,
			stack: []any{true, false},
			want:  true,
		},
		{
			name:  "OR same boolean",
			arg:   nil,
			stack: []any{true, true},
			want:  true,
		},
		{
			name:  "OR strings",
			arg:   nil,
			stack: []any{"test", "plan"},
			want:  false,
		},
		{
			name:  "OR float32",
			arg:   nil,
			stack: []any{float32(1.0), float32(6.6)},
			want:  true,
		},
	}

	for _, tt := range tests {
		syms := symbols.NewSymbolTable("testing")
		bc := ByteCode{}

		c := NewContext(syms, &bc)

		for _, item := range tt.stack {
			_ = c.push(item)
		}

		t.Run(tt.name, func(t *testing.T) {
			if tt.debug {
				fmt.Println("DEBUG")
			}

			err := target(c, tt.arg)

			if err != nil {
				e1 := nilError

				if tt.err != nil {
					e1 = tt.err.Error()
				}

				e2 := err.Error()
				if e1 == e2 {
					return
				}

				t.Errorf("%s() error %v", name, err)
			} else if tt.err != nil {
				t.Errorf("%s() expected error not reported: %v", name, tt.err)
			}

			v, err := c.Pop()

			if err != nil {
				t.Errorf("%s() stack error %v", name, err)
			}

			if !reflect.DeepEqual(v, tt.want) {
				t.Errorf("%s() got %v, want %v", name, v, tt.want)
			}
		})
	}
}

func Test_subtractByteCode(t *testing.T) {
	target := subtractByteCode
	name := "subByteCode"

	tests := []struct {
		name  string
		arg   any
		stack []any
		want  any
		err   error
		debug bool
	}{
		{
			name:  "sub string from error",
			arg:   nil,
			stack: []any{errors.ErrAssert, "-thing"},
			err:   errors.ErrInvalidType.Context("any"),
		},
		{
			name:  "sub with first nil",
			arg:   nil,
			stack: []any{2, nil},
			want:  7,
			err:   errors.ErrInvalidType.Context("nil"),
		},
		{
			name:  "sub with second nil",
			arg:   nil,
			stack: []any{nil, 5},
			want:  7,
			err:   errors.ErrInvalidType.Context("nil"),
		},
		{
			name:  "sub integers",
			arg:   nil,
			stack: []any{5, 2},
			want:  3,
		},
		{
			name:  "sub booleans",
			arg:   nil,
			stack: []any{true, false},
			want:  false,
			err:   errors.ErrInvalidType.Context("bool"),
		},
		{
			name:  "sub strings",
			arg:   nil,
			stack: []any{"foobar", "oba"},
			want:  "for",
		},
		{
			name:  "sub float32",
			arg:   nil,
			stack: []any{float32(1.0), float32(6.6)},
			want:  float32(-5.6),
		},
		{
			name:  "sub byte from int32 to byte",
			arg:   nil,
			stack: []any{int(5), byte(2)},
			want:  int(3),
		},
		{
			name:  "sub with 0 args on stack",
			arg:   nil,
			stack: []any{},
			want:  int32(12),
			err:   errors.ErrStackUnderflow,
		},
		{
			name:  "sub with 1 args on stack",
			arg:   nil,
			stack: []any{55},
			want:  int32(12),
			err:   errors.ErrStackUnderflow,
		},
	}

	for _, tt := range tests {
		syms := symbols.NewSymbolTable("testing")
		bc := ByteCode{}

		c := NewContext(syms, &bc)

		for _, item := range tt.stack {
			_ = c.push(item)
		}

		t.Run(tt.name, func(t *testing.T) {
			if tt.debug {
				fmt.Println("DEBUG")
			}

			err := target(c, tt.arg)

			if err != nil {
				e1 := nilError

				if tt.err != nil {
					e1 = tt.err.Error()
				}

				e2 := err.Error()
				if e1 == e2 {
					return
				}

				t.Errorf("%s() error %v", name, err)
			} else if tt.err != nil {
				t.Errorf("%s() expected error not reported: %v", name, tt.err)
			}

			v, err := c.Pop()

			if err != nil {
				t.Errorf("%s() stack error %v", name, err)
			}

			if !reflect.DeepEqual(v, tt.want) {
				t.Errorf("%s() got %v, want %v", name, v, tt.want)
			}
		})
	}
}

func Test_multiplyByteCode(t *testing.T) {
	target := multiplyByteCode
	name := "multiplyByteCode"

	tests := []struct {
		name  string
		arg   any
		stack []any
		want  any
		err   error
		debug bool
	}{
		{
			name:  "multiply with first nil",
			arg:   nil,
			stack: []any{2, nil},
			want:  7,
			err:   errors.ErrInvalidType.Context("nil"),
		},
		{
			name:  "multiply with second nil",
			arg:   nil,
			stack: []any{nil, 5},
			want:  7,
			err:   errors.ErrInvalidType.Context("nil"),
		},
		{
			name:  "multiply integers",
			arg:   nil,
			stack: []any{2, 5},
			want:  10,
		},
		{
			name:  "true OR false",
			arg:   nil,
			stack: []any{true, false},
			want:  true,
		},
		{
			name:  "false OR false",
			arg:   nil,
			stack: []any{false, false},
			want:  false,
		},
		{
			name:  "multiply strings",
			arg:   nil,
			stack: []any{"*", 5},
			want:  "*****",
		},
		{
			name:  "multiply float32",
			arg:   nil,
			stack: []any{float32(1.0), float32(6.6)},
			want:  float32(6.6),
		},
		{
			name:  "multiply int32 by byte",
			arg:   nil,
			stack: []any{int(5), byte(2)},
			want:  int(10),
		},
		{
			name:  "multiply with 0 args on stack",
			arg:   nil,
			stack: []any{},
			want:  int32(12),
			err:   errors.ErrStackUnderflow,
		},
		{
			name:  "multiply with 1 args on stack",
			arg:   nil,
			stack: []any{55},
			want:  int32(12),
			err:   errors.ErrStackUnderflow,
		},
	}

	for _, tt := range tests {
		syms := symbols.NewSymbolTable("testing")
		bc := ByteCode{}

		c := NewContext(syms, &bc)

		for _, item := range tt.stack {
			_ = c.push(item)
		}

		t.Run(tt.name, func(t *testing.T) {
			if tt.debug {
				fmt.Println("DEBUG")
			}

			err := target(c, tt.arg)

			if err != nil {
				e1 := nilError

				if tt.err != nil {
					e1 = tt.err.Error()
				}

				e2 := err.Error()

				if e1 == e2 {
					return
				}

				t.Errorf("%s() error %v", name, err)
			} else if tt.err != nil {
				t.Errorf("%s() expected error not reported: %v", name, tt.err)
			}

			v, err := c.Pop()

			if err != nil {
				t.Errorf("%s() stack error %v", name, err)
			}

			if !reflect.DeepEqual(v, tt.want) {
				t.Errorf("%s() got %v, want %v", name, v, tt.want)
			}
		})
	}
}

func Test_exponentyByteCode(t *testing.T) {
	target := exponentByteCode
	name := "exponentByteCode"

	tests := []struct {
		name  string
		arg   any
		stack []any
		want  any
		err   error
		debug bool
	}{
		{
			name:  "exponent with first nil",
			arg:   nil,
			stack: []any{2, nil},
			want:  7,
			err:   errors.ErrInvalidType.Context("nil"),
		},
		{
			name:  "exponent with second nil",
			arg:   nil,
			stack: []any{nil, 5},
			want:  7,
			err:   errors.ErrInvalidType.Context("nil"),
		},
		{
			name:  "exponent integers",
			arg:   nil,
			stack: []any{2, 5},
			want:  int64(32), // 2^5
		},
		{
			name:  "exponent true, false",
			arg:   nil,
			stack: []any{true, false},
			want:  true,
			err:   errors.ErrInvalidType.Context("bool"),
		},
		{
			name:  "exponent false, false",
			arg:   nil,
			stack: []any{false, false},
			want:  false,
			err:   errors.ErrInvalidType.Context("bool"),
		},
		{
			name:  "exponent strings",
			arg:   nil,
			stack: []any{"*", 5},
			err:   errors.ErrInvalidType.Context("string"),
		},
		{
			name:  "exponent float32",
			arg:   nil,
			stack: []any{float32(1.0), float32(6.6)},
			want:  float32(1),
		},
		{
			name:  "exponent int32 by byte",
			arg:   nil,
			stack: []any{int(5), byte(2)},
			want:  int64(25),
		},
		{
			name:  "multiply with 0 args on stack",
			arg:   nil,
			stack: []any{},
			want:  int32(12),
			err:   errors.ErrStackUnderflow,
		},
		{
			name:  "multiply with 1 args on stack",
			arg:   nil,
			stack: []any{55},
			want:  int32(12),
			err:   errors.ErrStackUnderflow,
		}}

	for _, tt := range tests {
		syms := symbols.NewSymbolTable("testing")
		bc := ByteCode{}

		c := NewContext(syms, &bc)

		for _, item := range tt.stack {
			_ = c.push(item)
		}

		t.Run(tt.name, func(t *testing.T) {
			if tt.debug {
				fmt.Println("DEBUG")
			}

			err := target(c, tt.arg)

			if err != nil {
				e1 := nilError

				if tt.err != nil {
					e1 = tt.err.Error()
				}

				e2 := err.Error()

				if e1 == e2 {
					return
				}

				t.Errorf("%s() error %v", name, err)
			} else if tt.err != nil {
				t.Errorf("%s() expected error not reported: %v", name, tt.err)
			}

			v, err := c.Pop()

			if err != nil {
				t.Errorf("%s() stack error %v", name, err)
			}

			if !reflect.DeepEqual(v, tt.want) {
				t.Errorf("%s() got %v, want %v", name, v, tt.want)
			}
		})
	}
}

func Test_divideByteCode(t *testing.T) {
	target := divideByteCode
	name := "divideByteCode"

	tests := []struct {
		name  string
		arg   any
		stack []any
		want  any
		err   error
		debug bool
	}{
		{
			name:  "divide with first nil",
			arg:   nil,
			stack: []any{2, nil},
			want:  7,
			err:   errors.ErrInvalidType.Context("nil"),
		},
		{
			name:  "divide with second nil",
			arg:   nil,
			stack: []any{nil, 5},
			want:  7,
			err:   errors.ErrInvalidType.Context("nil"),
		},

		{
			name:  "divide by integer zero",
			arg:   nil,
			stack: []any{9, 0},
			want:  3,
			err:   errors.ErrDivisionByZero,
		},
		{
			name:  "divide by float64 zero",
			arg:   nil,
			stack: []any{9, float64(0)},
			want:  3,
			err:   errors.ErrDivisionByZero,
		},
		{
			name:  "divide integers",
			arg:   nil,
			stack: []any{9, 3},
			want:  3,
		},
		{
			name:  "divide integers and discard remainder",
			arg:   nil,
			stack: []any{10, 3},
			want:  3,
		},
		{
			name:  "divide booleans",
			arg:   nil,
			stack: []any{true, false},
			want:  true,
			err:   errors.ErrInvalidType.Context("bool"),
		},
		{
			name:  "divide strings",
			arg:   nil,
			stack: []any{"*", 5},
			err:   errors.ErrInvalidType.Context("string"),
		},
		{
			name:  "divide float32 by integer",
			arg:   nil,
			stack: []any{float32(10.0), int(4)},
			want:  float32(2.5),
		},
		{
			name:  "divide int32 by byte",
			arg:   nil,
			stack: []any{int(12), byte(2)},
			want:  int(6),
		},
		{
			name:  "divide with 0 args on stack",
			arg:   nil,
			stack: []any{},
			want:  int32(12),
			err:   errors.ErrStackUnderflow,
		},
		{
			name:  "divide with 1 args on stack",
			arg:   nil,
			stack: []any{55},
			want:  int32(12),
			err:   errors.ErrStackUnderflow,
		},
	}

	for _, tt := range tests {
		syms := symbols.NewSymbolTable("testing")
		bc := ByteCode{}

		c := NewContext(syms, &bc)

		for _, item := range tt.stack {
			_ = c.push(item)
		}

		t.Run(tt.name, func(t *testing.T) {
			if tt.debug {
				fmt.Println("DEBUG")
			}

			err := target(c, tt.arg)

			if err != nil {
				e1 := nilError

				if tt.err != nil {
					e1 = tt.err.Error()
				}

				e2 := err.Error()

				if e1 == e2 {
					return
				}

				t.Errorf("%s() error %v", name, err)
			} else if tt.err != nil {
				t.Errorf("%s() expected error not reported: %v", name, tt.err)
			}

			v, err := c.Pop()

			if err != nil {
				t.Errorf("%s() stack error %v", name, err)
			}

			if !reflect.DeepEqual(v, tt.want) {
				t.Errorf("%s() got %v, want %v", name, v, tt.want)
			}
		})
	}
}

func Test_moduloByteCode(t *testing.T) {
	target := moduloByteCode
	name := "moduloByteCode"

	tests := []struct {
		name  string
		arg   any
		stack []any
		want  any
		err   error
		debug bool
	}{
		{
			name:  "modulo with first nil",
			arg:   nil,
			stack: []any{2, nil},
			want:  7,
			err:   errors.ErrInvalidType.Context("nil"),
		},
		{
			name:  "modulo with second nil",
			arg:   nil,
			stack: []any{nil, 5},
			want:  7,
			err:   errors.ErrInvalidType.Context("nil"),
		},
		{
			name:  "modulo integer zero",
			arg:   nil,
			stack: []any{9, 0},
			want:  3,
			err:   errors.ErrDivisionByZero,
		},
		{
			name:  "modulo integer",
			arg:   nil,
			stack: []any{10, 3},
			want:  1,
		},
		{
			name:  "modulo evenly divisible integers",
			arg:   nil,
			stack: []any{9, 3},
			want:  0,
		},
		{
			name:  "modulo booleans",
			arg:   nil,
			stack: []any{true, false},
			want:  true,
			err:   errors.ErrInvalidType.Context("bool"),
		},
		{
			name:  "modulo strings",
			arg:   nil,
			stack: []any{"*", 5},
			err:   errors.ErrInvalidType.Context("string"),
		},
		{
			name:  "modulo integer by byte",
			arg:   nil,
			stack: []any{15, byte(4)},
			want:  3,
		},
		{
			name:  "modulo with 0 args on stack",
			arg:   nil,
			stack: []any{},
			want:  int32(12),
			err:   errors.ErrStackUnderflow,
		},
		{
			name:  "modulo with 1 args on stack",
			arg:   nil,
			stack: []any{55},
			want:  int32(12),
			err:   errors.ErrStackUnderflow,
		},
	}

	for _, tt := range tests {
		syms := symbols.NewSymbolTable("testing")
		bc := ByteCode{}

		c := NewContext(syms, &bc)

		for _, item := range tt.stack {
			_ = c.push(item)
		}

		t.Run(tt.name, func(t *testing.T) {
			if tt.debug {
				fmt.Println("DEBUG")
			}

			err := target(c, tt.arg)

			if err != nil {
				e1 := nilError

				if tt.err != nil {
					e1 = tt.err.Error()
				}

				e2 := err.Error()
				if e1 == e2 {
					return
				}

				t.Errorf("%s() error %v", name, err)
			} else if tt.err != nil {
				t.Errorf("%s() expected error not reported: %v", name, tt.err)
			}

			v, err := c.Pop()

			if err != nil {
				t.Errorf("%s() stack error %v", name, err)
			}

			if !reflect.DeepEqual(v, tt.want) {
				t.Errorf("%s() got %v, want %v", name, v, tt.want)
			}
		})
	}
}

func Test_bitAndByteCode(t *testing.T) {
	target := bitAndByteCode
	name := "bitAndByteCode"

	tests := []struct {
		name  string
		arg   any
		stack []any
		want  any
		err   error
		debug bool
	}{
		{
			name:  "AND with first nil",
			arg:   nil,
			stack: []any{9, nil},
			err:   errors.ErrInvalidType.Context("nil"),
		},
		{
			name:  "AND with second nil",
			arg:   nil,
			stack: []any{nil, 0},
			want:  0,
			err:   errors.ErrInvalidType.Context("nil"),
		},
		{
			name:  "AND integer values",
			arg:   nil,
			stack: []any{5, 4},
			want:  4,
		},
		{
			name:  "AND int32 values",
			arg:   nil,
			stack: []any{int32(5), int32(4)},
			want:  4,
		},
		{
			name:  "AND boolean values",
			arg:   nil,
			stack: []any{true, true},
			want:  1,
		},
		{
			name:  "AND different boolean values",
			arg:   nil,
			stack: []any{true, false},
			want:  0,
		},
		{
			name:  "AND strings",
			arg:   nil,
			stack: []any{"7", 5},
			want:  5,
		},
		{
			name:  "AND with 0 args on stack",
			arg:   nil,
			stack: []any{},
			want:  int32(12),
			err:   errors.ErrStackUnderflow,
		},
		{
			name:  "AND with 1 args on stack",
			arg:   nil,
			stack: []any{55},
			want:  int32(12),
			err:   errors.ErrStackUnderflow,
		},
	}

	for _, tt := range tests {
		syms := symbols.NewSymbolTable("testing")
		bc := ByteCode{}

		c := NewContext(syms, &bc)

		for _, item := range tt.stack {
			_ = c.push(item)
		}

		t.Run(tt.name, func(t *testing.T) {
			if tt.debug {
				fmt.Println("DEBUG")
			}

			err := target(c, tt.arg)

			if err != nil {
				e1 := nilError

				if tt.err != nil {
					e1 = tt.err.Error()
				}

				e2 := err.Error()

				if e1 == e2 {
					return
				}

				t.Errorf("%s() error %v", name, err)
			} else if tt.err != nil {
				t.Errorf("%s() expected error not reported: %v", name, tt.err)
			}

			v, err := c.Pop()

			if err != nil {
				t.Errorf("%s() stack error %v", name, err)
			}

			if !reflect.DeepEqual(v, tt.want) {
				t.Errorf("%s() got %v, want %v", name, v, tt.want)
			}
		})
	}
}

func Test_bitOrByteCode(t *testing.T) {
	target := bitOrByteCode
	name := "bitOrByteCode"

	tests := []struct {
		name  string
		arg   any
		stack []any
		want  any
		err   error
		debug bool
	}{
		{
			name:  "OR with first nil",
			arg:   nil,
			stack: []any{9, nil},
			err:   errors.ErrInvalidType.Context("nil"),
		},
		{
			name:  "OR with second nil",
			arg:   nil,
			stack: []any{nil, 0},
			want:  0,
			err:   errors.ErrInvalidType.Context("nil"),
		},
		{
			name:  "OR integer zero",
			arg:   nil,
			stack: []any{9, 0},
			want:  9,
		},
		{
			name:  "OR integer values",
			arg:   nil,
			stack: []any{5, 4},
			want:  5,
		},
		{
			name:  "OR int32 values",
			arg:   nil,
			stack: []any{int32(5), int32(4)},
			want:  5,
		},
		{
			name:  "OR boolean values",
			arg:   nil,
			stack: []any{true, true},
			want:  1,
		},
		{
			name:  "OR different boolean values",
			arg:   nil,
			stack: []any{true, false},
			want:  1,
		},
		{
			name:  "OR flase boolean values",
			arg:   nil,
			stack: []any{false, false},
			want:  0,
		},
		{
			name:  "OR strings",
			arg:   nil,
			stack: []any{"7", 5},
			want:  7,
		},
		{
			name:  "OR with 0 args on stack",
			arg:   nil,
			stack: []any{},
			want:  int32(12),
			err:   errors.ErrStackUnderflow,
		},
		{
			name:  "OR with 1 args on stack",
			arg:   nil,
			stack: []any{55},
			want:  int32(12),
			err:   errors.ErrStackUnderflow,
		},
	}

	for _, tt := range tests {
		syms := symbols.NewSymbolTable("testing")
		bc := ByteCode{}

		c := NewContext(syms, &bc)

		for _, item := range tt.stack {
			_ = c.push(item)
		}

		t.Run(tt.name, func(t *testing.T) {
			if tt.debug {
				fmt.Println("DEBUG")
			}

			err := target(c, tt.arg)

			if err != nil {
				e1 := nilError

				if tt.err != nil {
					e1 = tt.err.Error()
				}

				e2 := err.Error()
				if e1 == e2 {
					return
				}

				t.Errorf("%s() error %v", name, err)
			} else if tt.err != nil {
				t.Errorf("%s() expected error not reported: %v", name, tt.err)
			}

			v, err := c.Pop()

			if err != nil {
				t.Errorf("%s() stack error %v", name, err)
			}

			if !reflect.DeepEqual(v, tt.want) {
				t.Errorf("%s() got %v, want %v", name, v, tt.want)
			}
		})
	}
}

func Test_bitShiftByteCode(t *testing.T) {
	target := bitShiftByteCode
	name := "bitShiftByteCode"

	tests := []struct {
		name  string
		arg   any
		stack []any
		want  any
		err   error
		debug bool
	}{
		{
			name:  "bitshift with first nil",
			arg:   nil,
			stack: []any{9, nil},
			err:   errors.ErrInvalidType.Context("nil"),
		},
		{
			name:  "bitshift with second nil",
			arg:   nil,
			stack: []any{nil, 0},
			want:  0,
			err:   errors.ErrInvalidType.Context("nil"),
		},
		{
			name:  "bitshift right 2 bits",
			arg:   nil,
			stack: []any{12, 2},
			want:  3,
		},
		{
			name:  "bitshift left 3 bits",
			arg:   nil,
			stack: []any{5, -3},
			want:  40,
		},
		{
			name:  "bitshift invalid bit count",
			arg:   nil,
			stack: []any{5, -35},
			err:   errors.ErrInvalidBitShift.Context(-35)},
		{
			name:  "shift with 0 args on stack",
			arg:   nil,
			stack: []any{},
			want:  int32(12),
			err:   errors.ErrStackUnderflow,
		},
		{
			name:  "shift with 1 args on stack",
			arg:   nil,
			stack: []any{55},
			want:  int32(12),
			err:   errors.ErrStackUnderflow,
		},
	}

	for _, tt := range tests {
		syms := symbols.NewSymbolTable("testing")
		bc := ByteCode{}

		c := NewContext(syms, &bc)

		for _, item := range tt.stack {
			_ = c.push(item)
		}

		t.Run(tt.name, func(t *testing.T) {
			if tt.debug {
				fmt.Println("DEBUG")
			}

			err := target(c, tt.arg)

			if err != nil {
				e1 := nilError

				if tt.err != nil {
					e1 = tt.err.Error()
				}

				e2 := err.Error()
				if e1 == e2 {
					return
				}

				t.Errorf("%s() error %v", name, err)
			} else if tt.err != nil {
				t.Errorf("%s() expected error not reported: %v", name, tt.err)
			}

			v, err := c.Pop()

			if err != nil {
				t.Errorf("%s() stack error %v", name, err)
			}

			if !reflect.DeepEqual(v, tt.want) {
				t.Errorf("%s() got %v, want %v", name, v, tt.want)
			}
		})
	}
}
