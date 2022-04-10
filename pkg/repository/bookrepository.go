package repository

import "github.com/hauntedness/jorm/pkg/entity"

//jorm-repository:"true"
type BookRepository[T entity.Book, K int] interface {
	FindById(k K) (book T, err error)
	FindByNameAndAuthor(name string, author string) (book T, err error)
	FindAllByName(name string) (books []T, err error)
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

var _ Repository[entity.Book, int] = &bookRepository{}
