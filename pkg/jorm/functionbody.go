package jorm

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

func NewFunctionBody() *FunctionBody {
	return &FunctionBody{
		VarQueryClause:  "",
		VarSelectClause: "",
		VarWhereClause:  "",
		VarExpression:   `var exp = selectClause + where + whereClause`,
		StmtQuery:       `rows, err := db.Query(exp, queryParams...)`,
		ForRowsNext:     `for rows.Next()`,
		LBrace:          "{",
		ForVarEntity:    "",
		ForStmtScan:     "",
		ForAppend:       "",
		RBrace:          "}",
		StmtReturn:      "return",
	}
}
