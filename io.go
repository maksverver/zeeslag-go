package io

import "./game"

import "os"
import "regexp"
import "strconv"
import "strings"

// ParseShips parses a canonical description of ships into a field array
func ParseShips(desc string) *game.Field {
	var field game.Field
	for _, ship := range (strings.Split(desc, ".", 0)) {
		const pattern = "^[2-5][HV][A-P]([1-9]|1[0-6])$"
		if matched, _ := regexp.MatchString(pattern, ship); !matched {
			return nil
		}
		var r1, c1, len, r2, c2 int
		if i, err := strconv.Atoi(ship[3:]); err == nil {
			r1 = i - 1
		}
		c1 = int(ship[2]) - int('A')
		len = int(ship[0]) - int('0')
		if ship[1] == 'H' {
			r2, c2 = r1, c1+len-1
		} else {
			r2, c2 = r1+len-1, c1
		}
		if r2 >= game.FieldHeight || c2 >= game.FieldWidth {
			return nil
		}
		for r := r1; r <= r2; r++ {
			for c := c1; c <= c2; c++ {
				field[r][c] = true
			}
		}
	}
	return &field
}

// ParseRows parses a canonical description of row counts
func ParseRows(desc string) *game.RowCounts {
	var res game.RowCounts
	parts := strings.Split(desc, ".", 0)
	if len(parts) != len(res) {
		return nil
	}
	for i, part := range (parts) {
		var err os.Error
		res[i], err = strconv.Atoi(part)
		if err != nil || res[i] < 0 || res[i] > game.FieldWidth {
			return nil
		}
	}
	return &res
}

// ParseCols parses a canonical description of column counts
func ParseCols(desc string) *game.ColCounts {
	var res game.ColCounts
	parts := strings.Split(desc, ".", 0)
	if len(parts) != len(res) {
		return nil
	}
	for i, part := range (parts) {
		var err os.Error
		res[i], err = strconv.Atoi(part)
		if err != nil || res[i] < 0 || res[i] > game.FieldHeight {
			return nil
		}
	}
	return &res
}
