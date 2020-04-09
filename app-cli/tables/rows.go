package tables

import (
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"
)

// AddRow adds a row to an existing table using an array of string objects,
// where each object represents a column of the data.
func (t *Table) AddRow(row []string) error {

	if len(row) != t.columnCount {
		return errors.New("Invalid column count in added row")
	}

	for n, h := range row {
		if len(h) > t.maxWidth[n] {
			t.maxWidth[n] = len(h)
		}
	}

	t.rows = append(t.rows, row)
	return nil
}

// AddRowItems adds a row to an existing table using individual parameters.
// Each parameter is converted to a string representation, and the set of all
// formatted values are added to the table as a row.
func (t *Table) AddRowItems(items ...interface{}) error {

	if len(items) != t.columnCount {
		return errors.New("Invalid column count in added row")
	}

	row := make([]string, t.columnCount)
	buffer := ""

	for n, item := range items {

		switch item.(type) {
		case int:
			buffer = strconv.Itoa(item.(int))

		case string:
			buffer = item.(string)

		case bool:
			if item.(bool) {
				buffer = "true"
			} else {
				buffer = "false"
			}
		default:
			buffer = fmt.Sprintf("%v", item)
		}
		row[n] = buffer
	}

	return t.AddRow(row)
}

// SortRows sorts the existing table rows. The column to sort by is specified by
// ordinal position (zero-based). The ascending flag is true if the sort is to be
// in ascending order, and false if a descending sort is required.
func (t *Table) SortRows(column int, ascending bool) error {
	if column < 0 || column >= t.columnCount {
		return errors.New("Invalid column number for sort")
	}

	sort.SliceStable(t.rows, func(i, j int) bool {
		if ascending {
			return t.rows[i][column] < t.rows[j][column]
		}
		return t.rows[i][column] > t.rows[j][column]
	})

	return nil
}

// SetOrderBy sets the name of the column that should be used for
// sorting the output data.
func (t *Table) SetOrderBy(name string) error {
	ascending := true
	if name[0] == '~' {
		name = name[1:]
		ascending = false
	}

	for n, v := range t.GetHeadings() {
		if strings.ToLower(name) == strings.ToLower(v) {
			t.orderBy = n
			t.ascending = ascending
			return nil
		}
	}
	return errors.New("Invalid order-by column name: " + name)
}
