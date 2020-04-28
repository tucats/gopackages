package tables

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/tucats/gopackages/app-cli/ui"
)

// Print will output a table using current rows and format specifications.
func (t *Table) Print(format string) error {

	// If there is an orderBy set for the table, do the sort now
	if t.orderBy >= 0 {
		t.SortRows(t.orderBy, t.ascending)
	}

	if format == "" {
		format = ui.TextTableFormat
	}

	// Based on the selected format, generate the output
	switch format {
	case ui.TextTableFormat:
		s := t.FormatText()
		for _, line := range s {
			fmt.Printf("%s\n", line)
		}

	case ui.JSONTableFormat:
		fmt.Printf("%s\n", t.FormatJSON())

	default:
		return errors.New("Invalid table format value")
	}
	return nil
}

// FormatJSON will produce the text of the table as JSON
func (t *Table) FormatJSON() string {

	var buffer strings.Builder

	buffer.WriteRune('[')
	for n, row := range t.rows {
		if n > 0 {
			buffer.WriteRune(',')
		}
		buffer.WriteRune('{')
		for i, header := range t.columns {
			if i > 0 {
				buffer.WriteRune(',')
			}
			buffer.WriteRune('"')
			buffer.WriteString(header)
			buffer.WriteString("\":")

			if _, valid := strconv.Atoi(row[i]); valid == nil {
				buffer.WriteString(row[i])
			} else if row[i] == "true" || row[i] == "false" {
				buffer.WriteString(row[i])
			} else {
				buffer.WriteString("\"" + row[i] + "\"")
			}
		}
		buffer.WriteRune('}')

	}
	buffer.WriteRune(']')
	return buffer.String()
}

// FormatText will output a table using current rows and format specifications.
func (t *Table) FormatText() []string {

	output := make([]string, 0)

	var buffer strings.Builder
	var rowLimit = t.rowLimit
	if rowLimit < 0 {
		rowLimit = len(t.rows)
	}

	if t.showHeadings {
		buffer.WriteString(t.indent)
		if t.showRowNumbers {
			buffer.WriteString("Row")
			buffer.WriteString(t.spacing)
		}
		for n, h := range t.columns {

			switch t.alignment[n] {
			case AlignmentLeft:
				buffer.WriteString(h)
				for pad := len(h); pad < t.maxWidth[n]; pad++ {
					buffer.WriteRune(' ')
				}
			case AlignmentRight:
				for pad := len(h); pad < t.maxWidth[n]; pad++ {
					buffer.WriteRune(' ')
				}
				buffer.WriteString(h)
			}
			buffer.WriteString(t.spacing)
		}
		output = append(output, buffer.String())

		if t.showUnderlines {
			buffer.Reset()
			buffer.WriteString(t.indent)
			if t.showRowNumbers {
				buffer.WriteString("===")
				buffer.WriteString(t.spacing)
			}
			for n := range t.columns {
				for pad := 0; pad < t.maxWidth[n]; pad++ {
					buffer.WriteRune('=')
				}
				buffer.WriteString(t.spacing)
			}
			output = append(output, buffer.String())
		}
	}

	for i, r := range t.rows {

		if i < t.startingRow {
			continue
		}
		if i >= t.startingRow+rowLimit {
			break
		}
		buffer.Reset()
		buffer.WriteString(t.indent)
		if t.showRowNumbers {
			buffer.WriteString(fmt.Sprintf("%3d", i+1))
			buffer.WriteString(t.spacing)
		}

		// Loop over the elements of the row. Generate pre- or post-spacing as
		// appropriate for the requested alignment, and any intra-column spacing.
		for n, c := range r {
			if t.alignment[n] == AlignmentRight {
				for pad := len(c); pad < t.maxWidth[n]; pad++ {
					buffer.WriteRune(' ')
				}
			}
			buffer.WriteString(c)
			if t.alignment[n] == AlignmentLeft {
				for pad := len(c); pad < t.maxWidth[n]; pad++ {
					buffer.WriteRune(' ')
				}
			}
			buffer.WriteString(t.spacing)
		}
		output = append(output, buffer.String())
	}

	return output
}
