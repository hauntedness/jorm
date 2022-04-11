package jorm

type Operand interface {
	BuildOper(string) string
}

type operand string

var (
	EQ    Operand = operand("Eq")      // A == B
	LT    Operand = operand("Lt")      // A < B
	GT    Operand = operand("Gt")      // A > B
	IN    Operand = operand("In")      // A in (strings.join([B,C,D],","))
	BT    Operand = operand("Between") // A between B and C
	NOTEQ Operand = operand("NotEq")   // A != B
	LE    Operand = operand("Le")      // A <= B
	GE    Operand = operand("Ge")      // A >= B
	NOTIN Operand = operand("NotIn")   // A not in B
)

func ConvertToOperand(str string) Operand {
	return operand(str)
}

func (op operand) BuildOper(column string) string {
	switch op {
	case EQ:
		return " = ?"
	case LT:
		return " < ?"
	case GT:
		return " > ?"
	case IN:
		return " in (" //todo mark as
	case BT:
		return " between ? and ?"
	case NOTEQ:
		return " <> ?"
	case LE:
		return " <= ?"
	case GE:
		return " >= ?"
	case NOTIN:
		return "not in (" //todo mark as
	default:
		return ""
	}
}
