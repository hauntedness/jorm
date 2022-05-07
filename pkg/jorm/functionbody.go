package jorm

import (
	"errors"
)

/*
var queryParams = make([]any, 0, len(names))
var selectClause = "SELECT id,name,author,version FROM book"
var where = "where"
var whereClause = jormgen.AddArray("name", names, queryParams) + " and " + jormgen.AddArray("author", authors, queryParams)
var exp = selectClause + " " + where + " " + whereClause
rows, err := db.Query(exp, queryParams...)
for rows.Next() {
	var book entity.Book
	rows.Scan(&book.Id, &book.Name, &book.Author, &book.Version)
	books = append(books, book)
}
return
*/
type FunctionBody struct {
	VarQueryClause  string
	VarSelectClause string
	VarWhereClause  string
	VarExpression   string
	StmtQuery       string
	ForRowsNext     string
	LBrace          string
	ForVarEntity    string
	ForStmtScan     string
	ForAppend       string
	RBrace          string
	StmtReturn      string
}

func (fb *FunctionBody) Build() string {
	return fb.VarQueryClause + "\n" + fb.VarSelectClause + "\n" + fb.VarWhereClause + "\n" + fb.VarExpression + "\n" + fb.StmtQuery + "\n" + fb.ForRowsNext + " " + fb.LBrace + "\n" + fb.ForVarEntity + "\n" + fb.ForStmtScan + "\n" + fb.ForAppend + "\n" + fb.RBrace + "\n" + fb.StmtReturn + "\n"
}

// var queryParams = make([]any, 0, len(names))
func NewFunctionBody() *FunctionBody {
	return &FunctionBody{
		VarQueryClause:  "var queryParams = make([]any, 0)",
		VarSelectClause: "",
		VarWhereClause:  "",
		VarExpression:   `var exp = selectClause + " where " + whereClause`,
		StmtQuery:       "rows, err := db.Query(exp, queryParams...)\nif err != nil {\nreturn nil, err\n}\ndefer rows.Close()",
		ForRowsNext:     `for rows.Next()`,
		LBrace:          "{",
		ForVarEntity:    "",
		ForStmtScan:     "",
		ForAppend:       "",
		RBrace:          "}",
		StmtReturn:      "return",
	}
}

//TODO below code is to optimize the performance
//when: there are only simple paramters
// there is no need to concate the sql dynamicly. it is enough that use one simple line to do all the things
// eg:  var whereClause = "name = ? and id = ?" rather than var whereClause = "name = ?" + " and " + "id = ?"
type WhereClause interface {
	Build() string
	AddVar(criteria Criteria)
}

//	var whereClause = jormgen.AddArray("name", names, queryParams) + " and " + jormgen.AddArray("author", authors, queryParams) + " and " + "id = ?"
type whereClause struct {
	noDynamic bool
	criterias []Criteria
}

var _ WhereClause = (*whereClause)(nil)
var _ Builder = (*whereClause)(nil)

func (w *whereClause) AddVar(criteria Criteria) {
	if _, ok := criteria.(DynamicCriteria); ok {
		w.noDynamic = false
	}
}

//TODO in function body ,we should based on noDynamic generate different source
func (w *whereClause) Build() (clause string) {
	for index, element := range w.criterias {
		if index == 0 {
			clause = element.Build()
		} else {
			clause = clause + `" and "` + element.Build()
		}
	}
	clause = "var whereClause = " + clause
	return
}

func NewWhereClause() WhereClause {
	return &whereClause{
		noDynamic: true,
		criterias: make([]Criteria, 0),
	}
}

type Criteria interface {
	Builder
}

// the whole can be quoted
type LiteralCriteria struct {
	value string
}

func NewLiteralCriteria(op RalationalOperator, column string) LiteralCriteria {
	var value string
	switch op {
	case OP_EQ:
		value = column + " = ?"
	case OP_LT:
		value = column + " < ?"
	case OP_GT:
		value = column + " > ?"
	case OP_NOTEQ:
		value = column + " <> ?"
	case OP_LE:
		value = column + " <= ?"
	case OP_GE:
		value = column + " >= ?"
	default:
		panic(errors.New("invalid operator for " + column))
	}
	return LiteralCriteria{value: `"` + value + `"`}
}

func (l LiteralCriteria) Build() string {
	return l.value
}

// only the literal parameter can be quoted
type DynamicCriteria struct{ value string }

func NewDynamicCriteria(op RalationalOperator, column string, paramName string) DynamicCriteria {
	var value string
	switch op {
	case OP_IN:
		value = `jormgen.AddIn("` + column + `", ` + paramName + `, queryParams)`
	case OP_NOTIN:
		value = `jormgen.AddNotIn("` + column + `", ` + paramName + `, queryParams)`
	default:
		panic("invalid operator: for " + column)
	}
	return DynamicCriteria{value: value}
}

func (l DynamicCriteria) Build() string {
	return l.value
}
