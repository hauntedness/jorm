package jorm

import (
	"go/ast"
	"go/token"
	"strings"

	"golang.org/x/exp/utf8string"
	"golang.org/x/tools/go/packages"
)

var (
	JORM_ENTITY     = "jorm-entity"
	JORM_TABLE      = "jorm-table"
	JORM_COLUMN     = "jorm-column"
	JORM_REPOSITORY = "jorm-repository"
)

type ORM struct {
	RootPath     string
	MappingStore map[string]*Mapping
}

func NewORM() *ORM {
	return &ORM{MappingStore: make(map[string]*Mapping)}
}

func (o *ORM) Parse(pkgs []*packages.Package) {
	// find all repositories and entities
	// one to one map
	// entity -> select column1, column2, column3
	// repository -> findbycolumn1andcolumn2. -> where column1 = param1 abd column2 = param2
	// code generate to file
}

func (o *ORM) iterate(pkgs []*packages.Package) {
	// iterate all packages
	for _, pkg := range pkgs {
		// iterate all files of each package
		for _, file := range pkg.Syntax {
			//iterate all declaritions of each file
			o.iterateFile(file, pkg)
		}
	}
}

// one ast file can hold many entities or repositories
// one file can't ensure a mapping would be done
// currently just treat the enity as key, so one entity can only have one mapping, one repository
func (o *ORM) iterateFile(file *ast.File, pkg *packages.Package) {
	//iterate all declaritions of each file
	for _, decl := range file.Decls {
		// find import, struct type and interface type
		if genDecl, ok := decl.(*ast.GenDecl); ok {
			o.iterateGenDecls(genDecl, file, pkg)
		}
	}
}

/**
 * iterate to refine usefull general declares,
 * the import -> which entity are references
 * eg: import "github.com/hauntedness/testast/pkg/entity"
 * struct -> entity
 * //jorm-entity:"true"
 * //jorm-table:"book"
	type Book struct {
		//jorm:entity:true1
		//jorm:table:book1
		_       struct{} `table:"book"`
		ID      int64    `column:"id"`
		Name    string   `column:"name"`
		Author  string   `column:"author"`
		Version int64    `column:"version"`
	}
 *
 * interface -> repository
 * 	type BookRepository[T entity.Book, K int] interface {
		FindById(k K) (book T, err error)
		FindByNameAndAuthor(name string, author string) (book T, err error)
		FindAllByName(name string) (books []T, err error)
	}
*/
func (o *ORM) iterateGenDecls(genDecl *ast.GenDecl, file *ast.File, pkg *packages.Package) {

	if genDecl.Tok == token.TYPE {
		var genDocList, tsDocList []*ast.Comment
		if genDecl.Doc != nil {
			genDocList = genDecl.Doc.List
		}
		for _, spec := range genDecl.Specs {
			typeSpec, ok := spec.(*ast.TypeSpec)
			if typeSpec.Doc != nil {
				tsDocList = typeSpec.Doc.List
			}

			if !ok || !typeSpec.Name.IsExported() {
				continue
			}

			switch typeSpec.Type.(type) {
			case *ast.StructType:
				o.onStructType(genDocList, tsDocList, typeSpec, file, pkg)
			case *ast.InterfaceType:
				o.onInterfaceType(genDocList, tsDocList, typeSpec, file, pkg)
			default:
				// do nothing
			}
		}
	}
}

func (o *ORM) onInterfaceType(genDocList, tsDocList []*ast.Comment, typeSpec *ast.TypeSpec, file *ast.File, pkg *packages.Package) {
	var key string
	// get doc to see if it is jorm repostiroy
	for _, comment := range genDocList {
		_, ok := extract(comment.Text, JORM_REPOSITORY)
		if !ok {
			return
		}
	}

	// here to find the entity
	// since here the x is short name, we need to go back to the file to find the long import
	// to find the import, we need to get the Imports (type []*ImportSpec) of the file
	// search the slice, there is 2 cases
	// case 1: import entity "github.com/hauntedness/testast/pkg/entity"
	// case 2: import "github.com/hauntedness/testast/pkg/entity"
	var entField = typeSpec.TypeParams.List[0]
	//TODO currently only allow entity and repostiory in different package
	se, ok := entField.Type.(*ast.SelectorExpr)
	if ok {
		var entityPath = serachImportPaths(se.X.(*ast.Ident).Name, file.Imports)
		if entityPath == "" {
			return
		}
		key = entityPath + "::" + se.Sel.Name
	}
	if key == "" {
		return
	}

	if _, ok := o.MappingStore[key]; !ok {
		o.MappingStore[key] = NewMapping()
	}
	o.MappingStore[key].Repository = typeSpec
	o.MappingStore[key].Status = o.MappingStore[key].Status | RepositoryReady
	o.MappingStore[key].BuildSqlText()
}

func (o *ORM) onStructType(genDocList, tsDocList []*ast.Comment, typeSpec *ast.TypeSpec, file *ast.File, pkg *packages.Package) {
	jormEntity := ""
	jormTable := ""
	for _, v := range genDocList {
		if strings.Contains(v.Text, JORM_ENTITY) {
			jormEntity = v.Text
		} else if strings.Contains(v.Text, JORM_TABLE) {
			jormTable = v.Text
		}
	}

	for _, v := range tsDocList {
		if strings.Contains(v.Text, JORM_ENTITY) {
			jormEntity = v.Text
		} else if strings.Contains(v.Text, JORM_TABLE) {
			jormTable = v.Text
		}
	}
	if jormEntity == "" {
		return
	}
	key := pkg.PkgPath + "::" + typeSpec.Name.Name
	if _, ok := o.MappingStore[key]; !ok {
		o.MappingStore[key] = NewMapping()
	}
	if value, ok := extract(jormTable, JORM_TABLE); ok {
		o.MappingStore[key].TableName = value
	}
	o.MappingStore[key].Entity = typeSpec
	o.MappingStore[key].EntityPath = key
	o.MappingStore[key].Status = o.MappingStore[key].Status | EntityReady
	o.MappingStore[key].OnEntityReady()
}

func serachImportPaths(x string, importPaths []*ast.ImportSpec) string {
	for _, is := range importPaths {
		runes := utf8string.NewString(is.Path.Value)
		trimed := runes.Slice(1, runes.RuneCount()-1)
		if is.Name != nil && x == is.Name.Name {
			return trimed
		} else if strings.HasSuffix(trimed, x) {
			return trimed
		}
	}
	return ""
}
