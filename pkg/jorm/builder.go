package jorm

import (
	"go/ast"
	"strings"
	"unicode"
)

type Builder interface {
	Build() string
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
	Recv   *ReceiverStatement
	FuncP  *FuncName
	Return *FuncReturn
	LBrace string
	Body   *FunctionBody
	RBrace string
}

type FuncName struct {
	Name   string
	LBrace string
	Fields []Field
	RBrace string
}

func NewFuncName(name string) *FuncName {
	return &FuncName{
		Name:   name,
		LBrace: "(",
		Fields: make([]Field, 0),
		RBrace: ")",
	}
}

//TODO add parameters consider slice
func (fn *FuncName) Build() string {
	var fields = make([]string, 0, len(fn.Fields))
	for _, field := range fn.Fields {
		fields = append(fields, field.Name+" "+field.Type)
	}
	if len(fn.Fields) > 0 {
		return fn.Name + fn.LBrace + " " + strings.Join(fields, ",") + " " + fn.RBrace
	} else {
		return fn.Name + fn.LBrace + fn.RBrace
	}
}

type FuncReturn struct {
	singleResult bool
	LBrace       string
	Fields       []Field
	RBrace       string
}

func NewFuncReturn(field *ast.Field, entity *ast.TypeSpec, pkg string) *FuncReturn {
	var singleResult = true
	var fields = make([]Field, 2)
	switch field.Type.(type) {
	case *ast.ArrayType:
		singleResult = false
		fields[0] = NewField(CaseTitleToCamel(entity.Name.Name)+"List", "[]"+pkg+`.`+entity.Name.Name)
	case *ast.Ident:
		fields[0] = NewField(CaseTitleToCamel(entity.Name.Name), pkg+`.`+entity.Name.Name)
	}
	fields[1] = NewField("err", "error")
	return &FuncReturn{
		singleResult: singleResult,
		LBrace:       "(",
		Fields:       fields,
		RBrace:       ")",
	}
}

func (fr *FuncReturn) SingleResult() bool {
	return fr.singleResult
}

func (fr *FuncReturn) Build() string {
	var params = make([]string, 0, len(fr.Fields))
	for _, param := range fr.Fields {
		params = append(params, param.Name+" "+param.Type)
	}
	if len(fr.Fields) > 0 {
		return fr.LBrace + strings.Join(params, ",") + fr.RBrace
	} else {
		return fr.LBrace + fr.RBrace
	}
}

type ReceiverStatement struct {
	LBrace string
	Alias  string
	Type   string
	RBrace string
}

func NewReceiverStatement(typ string) *ReceiverStatement {
	return &ReceiverStatement{
		LBrace: "(",
		Alias:  string(unicode.ToLower([]rune(typ)[0])),
		Type:   typ,
		RBrace: ")",
	}
}

func (rs *ReceiverStatement) Build() string {
	return rs.LBrace + rs.Alias + " *" + rs.Type + " " + rs.RBrace
}

func NewFunctionStatement(funcName string, recv string, returnType string) *FunctionStatement {
	return &FunctionStatement{
		Func:   "func",
		Recv:   NewReceiverStatement(recv),
		FuncP:  NewFuncName(funcName),
		Return: nil,
		LBrace: "{",
		Body:   NewFunctionBody(),
		RBrace: "}",
	}
}

func (fs *FunctionStatement) Build() string {
	return fs.Func + " " + fs.Recv.Build() + " " + fs.FuncP.Build() + " " + fs.Return.Build() + " " + fs.LBrace + "\n" + fs.Body.Build() + "\n" + fs.RBrace
}
