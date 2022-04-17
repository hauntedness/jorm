package jorm

import (
	"go/format"
	"testing"
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

func TestFunctionBodyBuild(t *testing.T) {
	var functionBody = &FunctionBody{
		VarQueryClause:  "var queryParams = make([]any, 0, len(names))",
		VarSelectClause: `var selectClause = "SELECT id,name,author,version FROM book"`,
		VarWhereClause:  `var whereClause = "where " + jormgen.AddArray("name", names, queryParams) + " and " + jormgen.AddArray("author", authors, queryParams)`,
		VarExpression:   `var exp = selectClause + " " + where + " " + whereClause`,
		StmtQuery:       `rows, err := db.Query(exp, queryParams...)`,
		ForRowsNext:     `for rows.Next()`,
		LBrace:          "{",
		ForVarEntity:    "var book entity.Book",
		ForStmtScan:     "rows.Scan(&book.Id, &book.Name, &book.Author, &book.Version)",
		ForAppend:       "books = append(books, book)",
		RBrace:          "}",
		StmtReturn:      "return",
	}
	s := functionBody.Build()
	t.Log("\n")
	t.Log(s)
	f := []byte(s)
	b, err := format.Source(f)
	if err != nil {
		t.Error(err)
	}
	t.Log(string(b))
}
