package jorm

type Operand string

// TODO, maybe it is too long for method name, consider shorter abbr
const (
	EQ  Operand = "Equals"      // A == B
	LT  Operand = "LessThan"    // A < B
	GT  Operand = "GreaterThan" // A > B
	IN  Operand = "In"          // A in (strings.join([B,C,D],","))
	NOT Operand = "Not"         // A != B
	BT  Operand = "Between"     // A between B and C
)
