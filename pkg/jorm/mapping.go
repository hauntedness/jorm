package jorm

import (
	"fmt"
	"go/ast"
	"strings"
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
	// only build once on entity ready
	// currently doesn't build after once, eg Repository Ready will not build again
	if m.Status == EntityReady {
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
		var criteria = make([]string, 0)

	methodLoop:
		for _, method := range interfaceType.Methods.List {
			// split function name into column list
			funcName := method.Names[0].Name
			// todo findAutherByName
			// todo findAllByName based on return type
			_, after, ok := strings.Cut(funcName, "FindBy")
			if !ok {
				continue
			}
			names := strings.Split(after, "And")
			for _, name := range names {
				if len(names) == 1 && (name == "Id" || name == "ID") {
					// TODO handle id particularly
					continue methodLoop
				}
				if field, ok := m.FieldMap[name]; ok {
					//TODO here should use the tag name
					if n, ok := ExtractTagValue(field, "jorm-column"); ok {
						criteria = append(criteria, n+" = ?")
					} else {
						criteria = append(criteria, name+" = ?")
					}
				} else {
					continue methodLoop
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
					continue methodLoop
				}
			}
			m.SqlText[funcName] = m.SelectClause + " " + m.TableName + " where " + strings.Join(criteria, " and ")
		}
		fmt.Println(interfaceType)
	}
}

func (m *Mapping) buildSelectSqlText(methodName string, map1 map[string]*ast.Field) {

}
