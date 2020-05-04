package tables

import (
	"errors"
	"strings"
)

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

// NewCSV creates a new table using a single string with comma-separated
// heading names. These typically correspond to the first row in a CSV
// data file.
func NewCSV(h string) (Table, error) {

	return New(csvSplit(h))
}

// CsvSplit takes a line that is comma-separated and splits it into
// an array of strings. Quoted commas are ignored as separators. The
// values are trimmed of extra spaces.
func CsvSplit(data string) []string {
	var headings []string
	var inQuote = false
	var currentHeading strings.Builder

	for _, c := range data {
		if c == '"' {
			inQuote = !inQuote
			continue
		}
		if !inQuote && c == ',' {
			headings = append(headings, strings.TrimSpace(currentHeading.String()))
			currentHeading.Reset()
			continue
		}
		currentHeading.WriteRune(rune(c))
	}

	if currentHeading.Len() > 0 {
		headings = append(headings, strings.TrimSpace(currentHeading.String()))
	}
	return headings
}
