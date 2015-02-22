package clitable

import (
	"bytes"
	"database/sql"
	"fmt"
	"strings"
)

func fetch_data(rows *sql.Rows, columns []string) (result [][]string, err error) {

	// get all data
	result = make([][]string, 1)
	result[0] = columns

	ptr := make([]interface{}, len(columns))
	data_row := make([]interface{}, len(columns))

	for i := range data_row {
		ptr[i] = &data_row[i]
	}

	for rows.Next() {
		if err = rows.Scan(ptr...); err != nil {
			return nil, err
		}

		row := make([]string, len(data_row))
		for i := range row {
			row[i] = fmt.Sprint(data_row[i])
		}

		result = append(result, row)
	}

	return result, nil
}

func get_max_len(data [][]string) []int {
	// draw table
	// need to now max data length in column
	p := make([]int, len(data[0]))
	for _, row := range data {
		for i, cell := range row {
			if p[i] < len(cell) {
				p[i] = len(cell)
			}
		}
	}
	return p
}

func Print(rows *sql.Rows) error {
	columns, err := rows.Columns()
	if err != nil {
		return err
	}

	data, err := fetch_data(rows, columns)
	if err != nil {
		return err
	}

	max_len := get_max_len(data)

	var buf bytes.Buffer
	var lbuf bytes.Buffer
	lbuf.WriteString("+")
	first := true
	line := ""

	for _, row := range data {
		buf.Reset()
		buf.WriteString("|")

		for i, cell := range row {
			n := buf.Len()
			if max_len[i] > len(cell) {
				pad := max_len[i] - len(cell)
				buf.WriteString(strings.Repeat(" ", pad/2+2))
				buf.WriteString(cell)
				buf.WriteString(strings.Repeat(" ", pad-(pad/2)+2))
				buf.WriteString("|")

			} else {
				buf.WriteString("  ")
				buf.WriteString(cell)
				buf.WriteString("  |")
			}

			if first {
				lbuf.WriteString(strings.Repeat("-", buf.Len()-n-1))
				lbuf.WriteString("+")
			}
		}

		if first {
			first = false
			line = lbuf.String()
			fmt.Println(line)
		}

		fmt.Println(buf.String())
		fmt.Println(line)
	}

	return rows.Close()
}
