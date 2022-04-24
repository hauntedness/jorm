package jorm

import (
	"errors"
	"go/ast"
)

type RalationalOperator interface {
	BuildElementExpression(columnName string, elemType ast.Expr, paramName string, paramType ast.Expr) string
}

type ralationalOperator string

var (
	OP_EQ    RalationalOperator = ralationalOperator("Eq")    // A == B
	OP_LT    RalationalOperator = ralationalOperator("Lt")    // A < B
	OP_GT    RalationalOperator = ralationalOperator("Gt")    // A > B
	OP_IN    RalationalOperator = ralationalOperator("In")    // A in (strings.join([B,C,D],","))
	OP_NOTEQ RalationalOperator = ralationalOperator("NotEq") // A <> B
	OP_LE    RalationalOperator = ralationalOperator("Le")    // A <= B
	OP_GE    RalationalOperator = ralationalOperator("Ge")    // A >= B
	OP_NOTIN RalationalOperator = ralationalOperator("NotIn") // A not in B
)

func NewRalationalOperator(str string) RalationalOperator {
	var op = ralationalOperator(str)
	switch op {
	case OP_EQ, OP_LT, OP_GT, OP_IN, OP_NOTEQ, OP_LE, OP_GE, OP_NOTIN:
		return op
	default:
		panic("invalid operator: " + str)
	}
}

/*
 #check1:
 - if a slice is here, the pattern is something like: name in (?,?,?,?,?,?,?,?,?,?,?) so the sql is dynamic,
 - the count of ? should be the length of the slice (or say can be inferred from names string[])
 - so the code would be something like:
 func (b *bookRepository) FindByNameIn(names []string) (books []entity.Book, err error){
 		var q = make([]string, 0, len(names))
 		for range names {
	 		q = append(q, "?")
 		}
 		...
 }
 rows, err := db.Query("SELECT name,author FROM book where id in ("+strings.Join(q, ",")+")", names)
 #check2:
 - get the corresponding type in entity of the names[i]
 - get the param type or get underlying type if the param is slice
 - jorm require that the two types must match
 #check3:
 - if the method name is like FindByNameIn or FindByNameNotIn, the param type must be Slice
 - and vice versa
*/
//TODO
func (op ralationalOperator) BuildElementExpression(columnName string, elemType ast.Expr, paramName string, paramType ast.Expr) string {
	if columnName == "" {
		panic("columnName is empty")
	}
	if op == OP_IN || op == OP_NOTIN {
		var paramArrayType *ast.ArrayType
		var paramIdent *ast.Ident
		var elemIdent *ast.Ident
		var ok bool
		paramArrayType, ok = paramType.(*ast.ArrayType)
		if !ok {
			panic("wrong param type")
		}
		paramIdent, ok = paramArrayType.Elt.(*ast.Ident)
		if !ok {
			panic("wrong param type")
		}
		elemIdent, ok = elemType.(*ast.Ident)
		if !ok {
			panic("wrong element type")
		}
		if paramIdent.Name != elemIdent.Name {
			panic("element type doesn't match param type")
		}
		return op.buildMultipleExp(columnName, paramName)
	} else {
		var paramIdent *ast.Ident
		var elemIdent *ast.Ident
		var ok bool
		paramIdent, ok = paramType.(*ast.Ident)
		if !ok {
			panic("wrong param type")
		}
		elemIdent, ok = elemType.(*ast.Ident)
		if !ok {
			panic("wrong element type")
		}
		if paramIdent.Name != elemIdent.Name {
			panic("type doesn't match")
		}
		return op.buildSingleValueExp(columnName)
	}
}

/**

func (b *bookRepository) FindByNameIn(names []string) (books []entity.Book, err error) {
	var querys = make([]string, 0, len(names))
	var params = make([]any, 0, len(names))
	for _, name := range names {
		querys = append(querys, "?")
		params = append(params, name)
	}
	var selectClause = "SELECT id,name,author,version FROM book"
	var whereClause = "where"
	var qtext = "id in (" + strings.Join(querys, ",") + ")"
	var exp = selectClause + " " + whereClause + " " + qtext
	rows, err := db.Query(exp, params...)
	for rows.Next() {
		var book entity.Book
		rows.Scan(&book.Id, &book.Name, &book.Author, &book.Version)
		books = append(books, book)
	}
	return
}
*/
func (op ralationalOperator) Build() []any {
	var clause []any
	// literal A = ?
	clause = append(clause, "A = ?")
	// list
	var list []any
	clause = append(clause, list...)
	return clause
}

func (op ralationalOperator) buildSingleValueExp(column string) string {
	switch op {
	case OP_EQ:
		return column + " = ?"
	case OP_LT:
		return column + " < ?"
	case OP_GT:
		return column + " > ?"
	case OP_NOTEQ:
		return column + " <> ?"
	case OP_LE:
		return column + " <= ?"
	case OP_GE:
		return column + " >= ?"
	default:
		panic(errors.New("invalid operator:" + string(op) + " for " + column))
	}
}

func (op ralationalOperator) buildMultipleExp(column string, paramName string) string {
	switch op {
	case OP_IN:
		return `jormgen.AddIn("` + column + `", ` + paramName + `, queryParams)`
	case OP_NOTIN:
		return `jormgen.AddNotIn("` + column + `", ` + paramName + `, queryParams)`
	default:
		panic("invalid operator:" + string(op) + " for " + column)
	}
}
