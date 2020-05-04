package tables

import (
	"errors"
	"strconv"
	"strings"
)

// GetHeadings returns an array of the headings already stored
// in the table. This can be used to validate a name against
// the list of headings, for example
func (t *Table) GetHeadings() []string {
	return t.columns
}

// SelectColumns accepts an array of column numbers that are to be considered
// part of the output. If this is not called, then all columns are assumed.
func (t *Table) SelectColumns(set []int) error {
	t.active = make([]bool, t.columnCount)
	for _, v := range set {
		if v < 1 || v > t.columnCount {
			return errors.New("Invalid column number: " + strconv.Itoa(v))
		}
		t.active[v-1] = true
	}
	return nil
}

// SelectColumn set the status of a specific column to active or inactive
func (t *Table) SelectColumn(n int, active bool) error {
	if n < 1 || n > t.columnCount {
		return errors.New("Invalid column number: " + strconv.Itoa(n))
	}
	t.active[n-1] = active

	return nil
}

// SelectColumnName sets the status of a specific column by name to
// active or inactive.
func (t *Table) SelectColumnName(name string, active bool) error {
	for columnNumber, columnName := range t.GetHeadings() {
		if strings.ToLower(columnName) == strings.ToLower(name) {
			err := t.SelectColumn(columnNumber+1, active)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
