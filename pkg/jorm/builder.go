package jorm

import (
	"strings"
)

type Builder interface {
	build() string
}

/*
A File represent implement of the repository interface
package xxx
import (
	"xxxx/entity"
	"error"
)

type bookRepository struct{}

func (b *bookRepository) FindByNameInAndAuthorIn(names []string, authors []string) (books []entity.Book, err error) {
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
}
*/
type File struct {
	PackageStmt   PackageStatement
	ImportStmts   ImportStatement
	StructStmts   []StructStatement
	FunctionStmts []FunctionStatement
}

type PackageStatement struct{}

func (ps *PackageStatement) Build() string {
	return "package jormgen"
}

func NewPackageStatement() *PackageStatement {
	return &PackageStatement{}
}

type ImportStatement struct {
	Import string
	LBrace string
	Pkgs   []string
	RBrace string
}

func (imprts *ImportStatement) AddImport(pkg string) {
	imprts.Pkgs = append(imprts.Pkgs, pkg)

}

func (imprts *ImportStatement) Build() string {
	if len(imprts.Pkgs) > 0 {
		return imprts.Import + " " + imprts.LBrace + "\n" + strings.Join(imprts.Pkgs, "\n") + "\n" + imprts.RBrace + "\n"
	} else {
		return "\n"
	}
}

func NewImportStatement() *ImportStatement {
	return &ImportStatement{
		Import: "import",
		LBrace: "(",
		Pkgs:   make([]string, 0),
		RBrace: ")",
	}
}

type StructStatement struct {
	Type   string
	Name   string
	Struct string
	LBrace string
	Fields []Field
	RBrace string
}

type Field struct {
	Name string
	Type string
}

func NewField(name string, typ string) Field {
	return Field{
		Name: name,
		Type: typ,
	}
}

func NewStructStatement(name string) *StructStatement {
	return &StructStatement{
		Type:   "type",
		Name:   name,
		Struct: "struct",
		LBrace: "{",
		Fields: make([]Field, 0),
		RBrace: "}",
	}
}

func (ss *StructStatement) Build() string {
	var fields = make([]string, 0, len(ss.Fields))
	for _, field := range ss.Fields {
		fields = append(fields, field.Name+" "+field.Type)
	}
	return ss.Type + " " + ss.Name + " " + ss.Struct + " " + ss.LBrace + "\n" + strings.Join(fields, "\n") + "\n" + ss.RBrace + "\n"
}

/*

func (b *bookRepository) FindByNameInAndAuthorIn(names []string, authors []string) (books []entity.Book, err error) {
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
}
*/
type FunctionStatement struct {
	Func   string
	Recv   ReceiverStatement
	FuncP  FuncName
	Return FuncReturn
	LBrace string
	Body   FunctionBody
	RBrace string
}

type FuncName struct {
	Name   string
	LBrace string
	Fields []Field
	RBrace string
}

func (fn *FuncName) Build() {
	var fields = make([]string, 0, len(fn.Fields))
	for _, field := range fn.Fields {
		fields = append(fields, field.Name+" "+field.Type)
	}
	if len(fn.Fields) > 0 {
		fn.Name = fn.Name + " " + fn.LBrace + "\n" + strings.Join(fields, "\n") + "\n" + fn.RBrace + "\n"
	} else {
		fn.Name = fn.Name + " " + fn.LBrace + "\n" + fn.RBrace + "\n"
	}
}

type FuncReturn struct {
	LBrace string
	Params []Field
	RBrace string
}

type ReceiverStatement struct {
	LBrace string
	Alias  string
	Type   string
	RBrace string
}

func NewReceiverStatement(alias string, typ string) *ReceiverStatement {
	return &ReceiverStatement{
		LBrace: "(",
		Alias:  alias,
		Type:   typ,
		RBrace: ")",
	}
}

func (rs *ReceiverStatement) Build() string {
	return rs.LBrace + rs.Alias + " " + rs.Type + " " + rs.RBrace
}

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

func NewFunctionStatement(funcName string, recv ReceiverStatement, returnType string) *FunctionStatement {
	return &FunctionStatement{
		Func:   "func",
		Recv:   recv,
		FuncP:  FuncName{Name: funcName, LBrace: "(", RBrace: ")"},
		Return: FuncReturn{LBrace: "(", RBrace: ")", Params: make([]Field, 0)},
	}
}
