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

type FunctionStatement struct {
	StructStatement
	Text   string
	Params []string
}

type Param struct {
}

type WhereClause struct {
	Param []string
	Text  string
}

func (wc *WhereClause) AddIdent() {
}

func (wc *WhereClause) AddArray() {

}
