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
	SqlText        map[string]string     // generated where clauses based on repository
	TableName      string                // table name
	EntityPath     string                // entity path, from import xxxx/xxx/xxx/xxx
	RepositoryPath string                // repository path from import xxx/xx/xx/repostoryxxx
	Status         MappingStatus         // enumerate mapping status
	FieldMap       map[string]*ast.Field // key is field name, value is field
}

func NewMapping() *Mapping {
	return &Mapping{SqlText: make(map[string]string), FieldMap: make(map[string]*ast.Field)}
}

func (m *Mapping) OnEntityReady() {
	if m.Status&EntityReady == EntityReady {
		st := m.Entity
		var selects []string = make([]string, 0, len(st.Fields.List))
		for _, field := range st.Fields.List {
			if len(field.Names) != 1 || !field.Names[0].IsExported() {
				continue
			}
			var theOnlyField *ast.Ident = field.Names[0]
			column, ok := ExtractTagValue(field, "jorm-column")
			// depending on use tag or field name
			if ok {
				selects = append(selects, column)
			} else {
				selects = append(selects, theOnlyField.Name)
			}
			m.FieldMap[theOnlyField.Name] = field
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
				m.buildSelectSqlText(methodName, method)
			} else if strings.HasPrefix(methodName, "Insert") {
				m.buildInsertSqlText(methodName, method)
			} else if strings.HasPrefix(methodName, "Update") {
				m.buildUpdateSqlText(methodName, method)
			} else if strings.HasPrefix(methodName, "Delete") {
				m.buildDeleteSqlText(methodName, method)
			}
		}
		fmt.Println(interfaceType)
	}
}

func (m *Mapping) buildSelectSqlText(methodName string, method *ast.Field) {
	var criteria = make([]string, 0)
	// split function name into column list
	funcName := method.Names[0].Name
	// todo findAutherByName
	// todo findAllByName based on return type
	_, after, ok := strings.Cut(funcName, "FindBy")
	if !ok {
		return
	}
	names := strings.Split(after, "And")
	for _, name := range names {
		if len(names) == 1 && (name == "Id" || name == "ID") {
			// TODO handle id particularly
			return
		}
		if field, op := m.ParseFieldNameAndOperand(name); field != nil {
			if n, ok := ExtractTagValue(field, "jorm-column"); ok {
				criteria = append(criteria, op.BuildOper(n))
			} else {
				criteria = append(criteria, op.BuildOper(n))
			}
		} else {
			return
		}
	}
	// check1: the length should match
	// the param type must match struct field type
	e := method.Type.(*ast.FuncType)
	for i, v := range e.Params.List {
		fieldOfEntity := m.FieldMap[names[i]]
		if v.Type.(*ast.Ident).Name == fieldOfEntity.Type.(*ast.Ident).Name {
			fmt.Println("good ")
		} else {
			return
		}
	}
	m.SqlText[funcName] = m.SelectClause + " " + m.TableName + " where " + strings.Join(criteria, " and ")
}

func (m *Mapping) buildInsertSqlText(methodName string, method *ast.Field) {

}

func (m *Mapping) buildUpdateSqlText(methodName string, method *ast.Field) {

}

func (m *Mapping) buildDeleteSqlText(methodName string, method *ast.Field) {

}

/**
type Book struct {
	Name         string
	NameLessThan string
}
current impl is match BookLessThan such a field 1st, then NameLessThan = "some thing"
*/
func (m *Mapping) ParseFieldNameAndOperand(section string) (field *ast.Field, op Operand) {
	utf8str := utf8string.NewString(section)
	var runeCount = utf8str.RuneCount()
	// try match by less runes
	for i := runeCount; i > 0; i-- {
		var ok bool
		field, ok = m.FieldMap[utf8str.Slice(0, i)]
		if ok {
			if i == runeCount {
				return field, EQ
			}
			op = NewOperand(utf8str.Slice(i, runeCount))
			return field, op
		}
	}
	// find field name prior to
	return nil, NewOperand("")
}
