package game

import (
	"container/vector"
	"os"
	"regexp"
	"strconv"
	"strings"
)

// ParseCoords parses a pair of field coordinates
func ParseCoords(desc string) (int, int, bool) {
	c := int(desc[0]) - int('A')
	r, err := strconv.Atoi(desc[1:])
	r--
	if err != nil || r < 0 || r >= FieldHeight || c < 0 || c >= FieldWidth {
		return 0, 0, false
	}
	return r, c, true
}

// FormatCoords converts a pair of field coordinates into a string description
func FormatCoords(r int, c int) string { return string('A'+c) + strconv.Itoa(r+1) }

// ParseShips parses a canonical description of ships into a field array
func ParseShips(desc string) *Field {
	var field Field
	for _, ship := range (strings.Split(desc, ".", 0)) {
		const pattern = "^[2-5][HV][A-P]([1-9]|1[0-6])$"
		if matched, _ := regexp.MatchString(pattern, ship); !matched {
			return nil
		}
		len := int(ship[0]) - int('0')
		r1, c1, ok := ParseCoords(ship[2:])
		if !ok {
			return nil
		}
		r2, c2 := r1, c1
		if ship[1] == 'H' {
			c2 += len - 1
		} else {
			r2 += len - 1
		}
		if r2 >= FieldHeight || c2 >= FieldWidth {
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

// FormatShips encodes a field in a string, as a series of ship placements
func FormatShips(field *Field) string {
	var parts vector.StringVector
	for r1 := 0; r1 < FieldHeight; r1++ {
		for c1 := 0; c1 < FieldWidth; c1++ {
			if field[r1][c1] &&
				(c1 == 0 || !field[r1][c1-1]) &&
				(r1 == 0 || !field[r1-1][c1]) {
				c2 := c1
				for c2 < FieldWidth && field[r1][c2] {
					c2++
				}
				r2 := r1
				for r2 < FieldHeight && field[r2][c1] {
					r2++
				}
				if c2-c1 > 1 {
					parts.Push(strconv.Itoa(c2-c1) + "H" + FormatCoords(r1, c1))
				}
				if r2-r1 > 1 {
					parts.Push(strconv.Itoa(r2-r1) + "V" + FormatCoords(r1, c1))
				}
			}
		}
	}
	return strings.Join(parts.Data(), ".")
}

// ParseRows parses a canonical description of row counts
func ParseRows(desc string) *RowCounts {
	var res RowCounts
	parts := strings.Split(desc, ".", 0)
	if len(parts) != len(res) {
		return nil
	}
	for i, part := range (parts) {
		var err os.Error
		res[i], err = strconv.Atoi(part)
		if err != nil || res[i] < 0 || res[i] > FieldWidth {
			return nil
		}
	}
	return &res
}

// ParseCols parses a canonical description of column counts
func ParseCols(desc string) *ColCounts {
	var res ColCounts
	parts := strings.Split(desc, ".", 0)
	if len(parts) != len(res) {
		return nil
	}
	for i, part := range (parts) {
		var err os.Error
		res[i], err = strconv.Atoi(part)
		if err != nil || res[i] < 0 || res[i] > FieldHeight {
			return nil
		}
	}
	return &res
}

// FormatCounts formats row or column counts into a canonical string format
func FormatCounts(counts []int) string {
	parts := make([]string, len(counts))
	for i, v := range (counts) {
		parts[i] = strconv.Itoa(v)
	}
	return strings.Join(parts, ".")
}

// ParseShots parses a canonical description of shots
func ParseShots(desc string) []Shot {
	if desc == "" {
		return make([]Shot, 0)
	}
	parts := strings.Split(desc, ".", 0)
	shots := make([]Shot, len(parts))
	for i, part := range (parts) {
		if len(part) < 3 {
			return nil
		}
		switch part[0] {
		case 'S':
			shots[i].Hit = true
		case 'W':
			shots[i].Hit = false
		default:
			return nil
		}
		if r, c, ok := ParseCoords(part[1:]); !ok {
			return nil
		} else {
			shots[i].R, shots[i].C = r, c
		}
	}
	return shots
}
