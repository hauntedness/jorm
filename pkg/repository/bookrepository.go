package repository

import (
	"strings"

	"github.com/hauntedness/jorm/pkg/entity"
)

//jorm-repository:"true"
type BookRepository[T entity.Book] interface {
	FindById(id int) (book T, err error)
	FindByNameAndAuthor(name string, author string) (book T, err error)
	FindByNameIn(name []string) (books []T, err error)
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

func (b *bookRepository) FindByNameIn(names []string) (books []entity.Book, err error) {
	var q = make([]string, 0, len(names))
	for range names {
		q = append(q, "?")
	}
	rows, err := db.Query("SELECT id,name,author,version FROM book where id in ("+strings.Join(q, ",")+")", names)
	for rows.Next() {
		var book entity.Book
		rows.Scan(&book.Id, &book.Name, &book.Author, &book.Version)
		books = append(books, book)
	}
	return
}

var _ BookRepository[entity.Book] = &bookRepository{}
