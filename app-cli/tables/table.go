package tables

import "errors"

const (

	// AlignmentLeft aligns the column to the left
	AlignmentLeft = -1

	// AlignmentRight aligns the column to the right
	AlignmentRight = 1

	// AlignmentCenter aligns the column to the center
	AlignmentCenter = 0
)

// Table is the wrapper object around a table to be printed
type Table struct {
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

// New creates a new table object, given a list of headings
func New(headings []string) (Table, error) {

	var t Table

	if len(headings) == 0 {
		return t, errors.New("cannot create table with zero columns")
	}
	t.rowLimit = -1
	t.columnCount = len(headings)
	t.columns = headings
	t.maxWidth = make([]int, len(headings))
	t.alignment = make([]int, len(headings))
	t.spacing = "    "
	t.indent = ""
	t.rows = make([][]string, 0)
	t.orderBy = -1
	t.ascending = true
	t.showUnderlines = true
	t.showHeadings = true

	for n, h := range headings {
		t.maxWidth[n] = len(h)
		t.columns[n] = h
		t.alignment[n] = AlignmentLeft
	}
	return t, nil
}

// GetHeadings returns an array of the headings already stored
// in the table. This can be used to validate a name against
// the list of headings, for example
func (t *Table) GetHeadings() []string {
	return t.columns
}
