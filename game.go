package game

const FieldHeight, FieldWidth = 16, 16

type Field [FieldHeight][FieldWidth]bool
type RowCounts [FieldHeight]int
type ColCounts [FieldWidth]int

type Shot struct {
	R, C int
	Hit  bool
}

var ShipLengths = [10]int{5, 4, 4, 3, 3, 3, 2, 2, 2, 2}

// CountShips computes the per row and column counts of ships in a field.
func CountShips(field *Field) (rows RowCounts, cols ColCounts) {
	for r, row := range (*field) {
		for c, cell := range (row) {
			if cell {
				rows[r]++
				cols[c]++
			}
		}
	}
	return
}

// formats a field as a string (useful for debug printing)
func (field *Field) String() string {
	result := ""
	for _, row := range (*field) {
		line := ""
		for _, cell := range (row) {
			if cell {
				line += "#"
			} else {
				line += "."
			}
		}
		result += line + "\n"
	}
	return result
}
