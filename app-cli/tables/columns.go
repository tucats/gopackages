package tables

import (
	"errors"
	"strconv"
	"strings"
)

// FindColumn returns the column number for a named column. The boolean return
// value indicates if the value was found, if true then the integer result is a
// zero-based column number.
func (t *Table) FindColumn(name string) (int, bool) {
	for n, v := range t.columns {
		if strings.ToLower(v) == strings.ToLower(name) {
			return n, true
		}
	}
	return -1, false
}

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

// SelectAllColumns sets the column selection attribute of all
// columns in the table.
func (t *Table) SelectAllColumns(active bool) {
	for n := range t.active {
		t.active[n] = active
	}
}

// SetColumnOrder accepts a list of column positions and uses it
// to set the order in which columns of output are printed.
func (t *Table) SetColumnOrder(order []int) error {
	if len(order) == 0 {
		return errors.New("invalid empty column order specification")
	}

	newOrder := make([]int, len(order))
	for n, v := range order {
		if v < 1 || v > t.columnCount {
			return errors.New("invalid column order specification: " + strconv.Itoa(v))
		}
		newOrder[n] = v - 1
	}
	t.columnOrder = newOrder
	return nil
}

// SetColumnOrderByName accepts a list of column positions and uses it
// to set the order in which columns of output are printed.
func (t *Table) SetColumnOrderByName(order []string) error {
	if len(order) == 0 {
		return errors.New("invalid empty column order specification")
	}

	newOrder := make([]int, len(order))
	for n, name := range order {
		v, found := t.FindColumn(name)
		if !found {
			return errors.New("invalid column order specification: " + strconv.Itoa(v))
		}
		newOrder[n] = v
	}
	t.columnOrder = newOrder
	return nil
}
