package repository

import (
	"github.com/hauntedness/jorm/pkg/entity"
	"github.com/hauntedness/jorm/pkg/repository/jormgen"
)

//jorm-repository:"true"
type BookRepository[T entity.Book] interface {
	FindById(id int) (book T, err error)
	FindByNameAndAuthor(name string, author string) (book T, err error)
	FindByNameInAndAuthorIn(name []string, author []string) (books []T, err error)
}

type bookRepository struct {
}

func (b *bookRepository) FindById(id int) (book entity.Book, err error) {
	row := db.QueryRow("SELECT id,name,author,version FROM book where id = ?", id)
	err = row.Scan(&book.Id, &book.Name, &book.Author, &book.Version)
	return
}

func (b *bookRepository) FindByNameAndAuthor(name string, author string) (book entity.Book, err error) {
	row := db.QueryRow("SELECT id,name,author,version FROM book where name = ? and author = ?", name, author)
	err = row.Scan(&book.Id, &book.Name, &book.Author, &book.Version)
	return
}

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

var _ BookRepository[entity.Book] = &bookRepository{}
