package tables

import "testing"

func TestTable_FormatJSON(t *testing.T) {
	type fields struct {
		showUnderlines bool
		showHeadings   bool
		showRowNumbers bool
		rowLimit       int
		startingRow    int
		columnCount    int
		rowCount       int
		orderBy        int
		ascending      bool
		rows           [][]string
		columns        []string
		alignment      []int
		maxWidth       []int
		spacing        string
		indent         string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "Single simple column",
			fields: fields{
				columnCount: 1,
				columns:     []string{"one"},
				rows:        [][]string{[]string{"1"}},
			},
			want: "[{\"one\":1}]",
		},
		{
			name: "Three columns of int, bool, string types",
			fields: fields{
				columnCount: 3,
				columns:     []string{"one", "two", "three"},
				rows:        [][]string{[]string{"1", "true", "Tom"}},
			},
			want: "[{\"one\":1,\"two\":true,\"three\":\"Tom\"}]",
		},
		{
			name: "Two rows of two columns",
			fields: fields{
				columnCount: 3,
				columns:     []string{"one", "two"},
				rows: [][]string{
					[]string{"60", "Tom"},
					[]string{"59", "Mary"},
				},
			},
			want: "[{\"one\":60,\"two\":\"Tom\"},{\"one\":59,\"two\":\"Mary\"}]",
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tx := &Table{
				showUnderlines: tt.fields.showUnderlines,
				showHeadings:   tt.fields.showHeadings,
				showRowNumbers: tt.fields.showRowNumbers,
				rowLimit:       tt.fields.rowLimit,
				startingRow:    tt.fields.startingRow,
				columnCount:    tt.fields.columnCount,
				rowCount:       tt.fields.rowCount,
				orderBy:        tt.fields.orderBy,
				ascending:      tt.fields.ascending,
				rows:           tt.fields.rows,
				columns:        tt.fields.columns,
				alignment:      tt.fields.alignment,
				maxWidth:       tt.fields.maxWidth,
				spacing:        tt.fields.spacing,
				indent:         tt.fields.indent,
			}
			if got := tx.FormatJSON(); got != tt.want {
				t.Errorf("Table.FormatJSON() = %v, want %v", got, tt.want)
			}
		})
	}
}
