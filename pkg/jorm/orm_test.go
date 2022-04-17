package jorm

import (
	"flag"
	"fmt"
	"os"
	"testing"

	"golang.org/x/tools/go/packages"
)

func TestOrmMap(t *testing.T) {
	flag.Parse()

	// Many tools pass their command-line arguments (after any flags)
	// uninterpreted to packages.Load so that it can interpret them
	// according to the conventions of the underlying build system.
	cfg := &packages.Config{Mode: packages.NeedFiles | packages.NeedSyntax | packages.NeedName}
	//pkgs, err := packages.Load(cfg, flag.Args()...)
	pkgs, err := packages.Load(cfg, "github.com/hauntedness/jorm/pkg/entity", "github.com/hauntedness/jorm/pkg/repository")
	//github.com/hauntedness/testast/pkg/entity
	if err != nil {
		fmt.Fprintf(os.Stderr, "load: %v\n", err)
		t.Error(err)
	}
	if packages.PrintErrors(pkgs) > 0 {
		t.Error(err)
	}

	// Print the names of the source files
	// for each package listed on the command line.
	var orm = NewORM()
	orm.iterate(pkgs)
	for _, v := range orm.MappingStore {
		fmt.Printf("%#v", v)
	}
}

//
func TestGetColumnName(t *testing.T) {
	text := `column:   " ve   rsion"`
	c, ok := extract(text, "column")
	if !ok {
		t.Error("test fail")
	}
	if c != "version" {
		t.Error("test fail")
	}
	fmt.Println(c)
}

func TestGetTableName(t *testing.T) {
	text := `table:"book"`
	c, ok := extract(text, "table")
	if !ok {
		t.Error("test fail")
	}
	if c != "book" {
		t.Error("test fail")
	}
	fmt.Println(c)
}
