package tables

import (
	"reflect"
	"testing"
)

func TestTable_SortRows(t *testing.T) {

	tests := []struct {
		name       string
		headers    []string
		rows       [][]string
		sortColumn string
		result     []string
	}{
		{
			name:       "Simple table with one column, two rows",
			headers:    []string{"first"},
			rows:       [][]string{[]string{"v2"}, []string{"v1"}},
			sortColumn: "first",
			result:     []string{"first    ", "=====    ", "v1       ", "v2       "},
		},
		{
			name:       "Simple table with two columns, two rows",
			headers:    []string{"first", "second"},
			rows:       [][]string{[]string{"v2", "d1"}, []string{"v1", "d2"}},
			sortColumn: "first",
			result:     []string{"first    second    ", "=====    ======    ", "v1       d2        ", "v2       d1        "},
		},
		{
			name:       "Simple table with two columns, two rows, alternate sort",
			headers:    []string{"first", "second"},
			rows:       [][]string{[]string{"v2", "d1"}, []string{"v1", "d2"}},
			sortColumn: "second",
			result:     []string{"first    second    ", "=====    ======    ", "v2       d1        ", "v1       d2        "},
		},
		// TODO add tests
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			table, _ := New(tt.headers)
			for _, r := range tt.rows {
				table.AddRow(r)
			}
			table.SetOrderBy(tt.sortColumn)
			table.SortRows(table.orderBy, table.ascending)
			x := table.FormatText()

			if !reflect.DeepEqual(x, tt.result) {
				t.Errorf("Sorted row results wrong. Got %v want %v", x, tt.result)
			}
		})
	}
}

func TestTable_AddRow(t *testing.T) {

	type args struct {
		row []string
	}
	tests := []struct {
		name    string
		table   Table
		args    args
		want    Table
		wantErr bool
	}{
		{
			name: "Add one row to a single column",
			table: Table{
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
			args: args{
				row: []string{"first"},
			},
			want: Table{
				rowLimit:       -1,
				columnCount:    1,
				columns:        []string{"simple"},
				maxWidth:       []int{6},
				alignment:      []int{AlignmentLeft},
				spacing:        "    ",
				indent:         "",
				rows:           [][]string{[]string{"first"}},
				orderBy:        -1,
				ascending:      true,
				showUnderlines: true,
				showHeadings:   true,
			},
			wantErr: false,
		},
		{
			name: "Add one row to a three-column table",
			table: Table{
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
			args: args{
				row: []string{"first", "second", "third"},
			},
			want: Table{
				rowLimit:       -1,
				columnCount:    3,
				columns:        []string{"simple", "test", "table"},
				maxWidth:       []int{6, 6, 5},
				alignment:      []int{AlignmentLeft, AlignmentLeft, AlignmentLeft},
				spacing:        "    ",
				indent:         "",
				rows:           [][]string{[]string{"first", "second", "third"}},
				orderBy:        -1,
				ascending:      true,
				showUnderlines: true,
				showHeadings:   true,
			},
			wantErr: false,
		},
		{
			name: "Add two-column row to three-column table",
			table: Table{
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
			args: args{
				row: []string{"first", "second"},
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
			wantErr: true,
		},

		// TODO: Add test cases.
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ttable := &tt.table
			err := ttable.AddRow(tt.args.row)
			if (err != nil) != tt.wantErr {
				t.Errorf("Table.AddRow() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(*ttable, tt.want) {
				t.Errorf("Table.AddRow() got %v, want %v", ttable, tt.want)
			}
		})
	}
}
