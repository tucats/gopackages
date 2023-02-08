package bytecode

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/tucats/gopackages/errors"
	"github.com/tucats/gopackages/expressions/data"
	"github.com/tucats/gopackages/expressions/symbols"
)

func Test_loadByteCode(t *testing.T) {
	target := loadByteCode
	name := "loadByteCode"

	tests := []struct {
		name         string
		arg          interface{}
		initialValue interface{}
		stack        []interface{}
		want         interface{}
		err          error
		static       int
		debug        bool
	}{
		{
			name:         "simple integer load",
			arg:          "a",
			initialValue: int32(55),
			stack:        []interface{}{},
			static:       2,
			err:          nil,
			want:         int32(55),
		},
		{
			name:   "variable not found",
			arg:    "a",
			stack:  []interface{}{},
			static: 2,
			err:    errors.ErrUnknownIdentifier.Context("a"),
			want:   int32(55),
		},
		{
			name:   "variable name invalid",
			arg:    "",
			stack:  []interface{}{},
			static: 2,
			err:    errors.ErrInvalidIdentifier.Context(""),
			want:   int32(55),
		},
	}

	for _, tt := range tests {
		syms := symbols.NewSymbolTable("testing")
		bc := ByteCode{}
		varname := data.String(tt.arg)

		c := NewContext(syms, &bc)
		c.typeStrictness = tt.static

		for _, item := range tt.stack {
			_ = c.push(item)
		}

		if tt.initialValue != nil {
			_ = c.create(varname)
			_ = c.set(varname, tt.initialValue)
		}

		t.Run(tt.name, func(t *testing.T) {
			if tt.debug {
				fmt.Println("DEBUG")
			}

			err := target(c, tt.arg)

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

			v, err := c.Pop()

			if err != nil {
				t.Errorf("%s() error popping value from stack: %v", name, tt.arg)
			}

			if !reflect.DeepEqual(v, tt.want) {
				t.Errorf("%s() got %v, want %v", name, v, tt.want)
			}
		})
	}
}
