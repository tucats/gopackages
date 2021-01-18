package bytecode

import (
	"reflect"
	"testing"

	"github.com/tucats/gopackages/symbols"
)

func TestStructImpl(t *testing.T) {

	tests := []struct {
		name    string
		stack   []interface{}
		arg     interface{}
		want    interface{}
		wantErr bool
		static  bool
	}{
		{
			name:    "two member invalid static test",
			arg:     3,
			stack:   []interface{}{"usertype", "__type", true, "valid", 123, "test"},
			want:    map[string]interface{}{"active": true, "test": 123, "__static": true},
			wantErr: true,
			static:  true,
		},
		{
			name:    "one member test",
			arg:     1,
			stack:   []interface{}{123, "test"},
			want:    map[string]interface{}{"test": 123, "__static": true},
			wantErr: false,
		},
		{
			name:    "two member test",
			arg:     2,
			stack:   []interface{}{true, "active", 123, "test"},
			want:    map[string]interface{}{"active": true, "test": 123, "__static": true},
			wantErr: false,
		},
		{
			name:    "two member valid static test",
			arg:     2,
			stack:   []interface{}{true, "active", 123, "test"},
			want:    map[string]interface{}{"active": true, "test": 123, "__static": true},
			wantErr: false,
			static:  true,
		},
		{
			name:    "two member invalid static test",
			arg:     3,
			stack:   []interface{}{"usertype", "__type", true, "valid", 123, "test"},
			want:    map[string]interface{}{"active": true, "test": 123, "__static": true},
			wantErr: true,
			static:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := &Context{
				stack:   tt.stack,
				sp:      len(tt.stack),
				Static:  tt.static,
				symbols: symbols.NewSymbolTable("test bench"),
			}
			_ = ctx.symbols.SetAlways("usertype", map[string]interface{}{
				"test":   0,
				"static": false})

			err := StructImpl(ctx, tt.arg)
			if (err != nil) != tt.wantErr {
				t.Errorf("StructImpl() error = %v, wantErr %v", err, tt.wantErr)
			} else if err == nil {
				got, _ := ctx.Pop()
				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("StructImpl() got = %v, want %v", got, tt.want)
				}
			}
		})
	}
}
