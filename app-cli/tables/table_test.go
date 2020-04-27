package tables

import (
	"reflect"
	"testing"
)

func TestNew(t *testing.T) {
	type args struct {
		headings []string
	}
	tests := []struct {
		name      string
		args      args
		want      Table
		wantError bool
	}{
		{
			name: "Simple table with one column",
			args: args{
				headings: []string{"simple"},
			},
			want: Table{
				rowLimit:       -1,
				columnCount:    1,
				columns:        []string{"simple"},
				maxWidth:       []int{6},
				alignment:      []int{AlignmentLeft},
				spacing:        "    ",
				indent:         "",
				rows:           make([][]string, 0),
				orderBy:        -1,
				ascending:      true,
				showUnderlines: true,
				showHeadings:   true,
			},
			wantError: false,
		},
		{
			name: "Simple table with three columns",
			args: args{
				headings: []string{"simple", "test", "table"},
			},
			want: Table{
				rowLimit:       -1,
				columnCount:    3,
				columns:        []string{"simple", "test", "table"},
				maxWidth:       []int{6, 4, 5},
				alignment:      []int{AlignmentLeft, AlignmentLeft, AlignmentLeft},
				spacing:        "    ",
				indent:         "",
				rows:           make([][]string, 0),
				orderBy:        -1,
				ascending:      true,
				showUnderlines: true,
				showHeadings:   true,
			},
			wantError: false,
		},
		{
			name: "Invalid table with no columns",
			args: args{
				headings: []string{},
			},
			want:      Table{},
			wantError: true,
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := New(tt.args.headings)
			if err != nil && !tt.wantError {
				t.Errorf("New() resulted in unexpected error %v", err)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}
