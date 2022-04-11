package jorm

type Operand interface {
	BuildOper(string) string
}

type operand string

var (
	OP_EQ    Operand = operand("Eq")    // A == B
	OP_LT    Operand = operand("Lt")    // A < B
	OP_GT    Operand = operand("Gt")    // A > B
	OP_IN    Operand = operand("In")    // A in (strings.join([B,C,D],","))
	OP_NOTEQ Operand = operand("NotEq") // A != B
	OP_LE    Operand = operand("Le")    // A <= B
	OP_GE    Operand = operand("Ge")    // A >= B
	OP_NOTIN Operand = operand("NotIn") // A not in B
)

func NewOperand(str string) Operand {
	return operand(str)
}

func (op operand) BuildOper(column string) string {
	switch op {
	case OP_EQ:
		return column + " = ?"
	case OP_LT:
		return column + " < ?"
	case OP_GT:
		return column + " > ?"
	case OP_IN:
		return column + " in (" //todo mark as
	case OP_NOTEQ:
		return column + " <> ?"
	case OP_LE:
		return column + " <= ?"
	case OP_GE:
		return column + " >= ?"
	case OP_NOTIN:
		return column + "not in (" //todo mark as
	default:
		return ""
	}
}
