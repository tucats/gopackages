// Package i18n provides localization and internationalization
// functionality for Ego itself.
package i18n

import (
	"testing"
)

func TestO(t *testing.T) {
	type args struct {
		key      string
		valueMap []map[string]interface{}
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "key exists",
			args: args{
				key: "table.read.row.ids",
			},
			want: "Include the row UUID column in the output",
		},
		{
			name: "key does not exist",
			args: args{
				key: "Do the things",
			},
			want: "Do the things",
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := O(tt.args.key, tt.args.valueMap...); got != tt.want {
				t.Errorf("O() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRegister(t *testing.T) {
	t.Run("Register", func(t *testing.T) {
		localizations := map[string]map[string]string{
			"test.key": {
				"en": "test key",
			},
			"another.key": {
				"en": "english",
				"fr": "french",
			},
		}

		Register(localizations)

		want := "test key"
		if m := T("test.key"); m != want {
			t.Errorf("got %s, want %s", m, want)
		}

		want = "not a key"
		if m := T("not a key"); m != want {
			t.Errorf("got %s, want %s", m, want)
		}

		want = "incorrect function argument type"
		if m := T("error.arg.type"); m != want {
			t.Errorf("got %s, want %s", m, want)
		}

	})

}
