package jorm

import (
	"fmt"
	"go/ast"
	"strings"

	"golang.org/x/exp/utf8string"
)

type MappingStatus int

const (
	EntityReady MappingStatus = 1 << iota
	RepositoryReady
)

const ERMapReady = EntityReady & RepositoryReady

type Mapping struct {
	Entity         *ast.StructType       // entity struct to be parsed
	Repository     *ast.InterfaceType    // repository interface to be parsed
	PrimaryKeys    map[string]any        // the type of pk, can be int, float64, strign and so on... todo how about struct(??)
	SelectClause   string                // generated select clause
	CodeText       map[string]string     // generated where clauses based on repository
	TableName      string                // table name
	EntityPath     string                // entity path, from import xxxx/xxx/xxx/xxx
	RepositoryPath string                // repository path from import xxx/xx/xx/repostoryxxx
	Status         MappingStatus         // enumerate mapping status
	FieldMap       map[string]*ast.Field // key is field name, value is field
}

func NewMapping() *Mapping {
	return &Mapping{CodeText: make(map[string]string), FieldMap: make(map[string]*ast.Field)}
}

func (m *Mapping) OnEntityReady() {
	if m.Status&EntityReady == EntityReady {
		st := m.Entity
		var selects []string = make([]string, 0, len(st.Fields.List))
		for _, field := range st.Fields.List {
			if len(field.Names) != 1 || !field.Names[0].IsExported() {
				continue
			}
			column, ok := ExtractTagValue(field, "jorm-column")
			// depending on use tag or field name
			if ok {
				selects = append(selects, column)
			} else {
				selects = append(selects, field.Names[0].Name)
			}
			m.FieldMap[field.Names[0].Name] = field
		}
		m.SelectClause = "select " + strings.Join(selects, ", ")
	}
}

func (m *Mapping) BuildSqlText() {
	if m.Status&RepositoryReady == RepositoryReady {
		interfaceType := m.Repository
		for _, method := range interfaceType.Methods.List {
			methodName := method.Names[0].Name
			if strings.HasPrefix(methodName, "Find") {
				m.buildSelectCodeText(method)
			} else if strings.HasPrefix(methodName, "Insert") {
				m.buildInsertCodeText(method)
			} else if strings.HasPrefix(methodName, "Update") {
				m.buildUpdateCodeText(method)
			} else if strings.HasPrefix(methodName, "Delete") {
				m.buildDeleteCodeText(method)
			}
		}
		fmt.Println(interfaceType)
	}
}

func (m *Mapping) buildSelectCodeText(method *ast.Field) {
	var criteria = make([]string, 0)
	// split function name into column list
	funcName := method.Names[0].Name
	params := method.Type.(*ast.FuncType).Params.List
	// todo findAutherByName
	// todo findAllByName based on return type
	_, after, ok := strings.Cut(funcName, "FindBy")
	if !ok {
		return
	}
	names := strings.Split(after, "And")
	for i, name := range names {
		field, op := m.ParseFieldNameAndOperator(name)
		if field == nil {
			return
		}
		// check for type match
		if m.FieldMap[name].Type.(*ast.Ident).Name != params[i].Type.(*ast.Ident).Name {
			return
		}
		var columnName string
		var ok bool
		if columnName, ok = ExtractTagValue(field, "jorm-column"); !ok {
			columnName = m.FieldMap[name].Names[0].Name
		}
		exp, err := op.BuildElementExpression(columnName, m.FieldMap[name].Type, params[i].Names[0].Name, params[i].Type)
		if err != nil {
			panic(err)
		}
		//TODO
		criteria = append(criteria, exp)
	}

	m.CodeText[funcName] = m.SelectClause + " " + m.TableName + " where " + strings.Join(criteria, " and ")
}

func (m *Mapping) buildInsertCodeText(method *ast.Field) {

}

func (m *Mapping) buildUpdateCodeText(method *ast.Field) {

}

func (m *Mapping) buildDeleteCodeText(method *ast.Field) {

}

/*
type Book struct {
	Name         string
	NameLessThan string
}

current impl is match BookLessThan such a field 1st, then NameLessThan = "some thing"
*/
func (m *Mapping) ParseFieldNameAndOperator(section string) (field *ast.Field, op RalationalOperator) {
	utf8str := utf8string.NewString(section)
	var runeCount = utf8str.RuneCount()
	// try match by less runes
	for i := runeCount; i > 0; i-- {
		var ok bool
		field, ok = m.FieldMap[utf8str.Slice(0, i)]
		if ok {
			if i == runeCount {
				return field, OP_EQ
			}
			op = NewRalationalOperator(utf8str.Slice(i, runeCount))
			return field, op
		}
	}
	// find field name prior to
	return nil, NewRalationalOperator("")
}
