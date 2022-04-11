# jorm
an orm library for go, based on go generate, zero document needed, inspired by jpa


### Agenda 

- below functions list what is going to be implented or not
```go
type BookRepository[T entity.Book, K int] interface {
    //done
    FindByNameAndAuthor(name string, author string) (T, error)
    //needed? this signature is not good for readness
    FindByNameAndAuthor(name, author string) (T, error) 
    //in progress
    FindById(k K) (T,error)
    //in progress
    FindAuthorByName(name string) (string , error)
    //in progress
    FindByName(name string) ([]T , error)
    //in progress
    FindAuthorByName(name string) ([]string, error)
    //in progress where created_date < ptime
    FindByCreatedDateLt(ptime time.Time) ([]T, error) 
    //in progress 
    //batch selection is not worth here
    FindByNameIn(names []string) ([]T, error) 
    //in progress 
    // batch selection is not worth here
    FindByNameNotIn(names []string) ([]T, error)
    //in progress will not provide auto generated id
    Insert(book T) error
    //in progress will not provide auto generated id
    InsertAll(books []T) (int, error)
    //in progress
    //same to update by id?
    Update(book T) error
    //in progress
    UpdateAll(books []T) (int, error)
    //to keep the similar order to update author = ? where name = ?
    //needn't to return the number of updated rows?
    //in progress
    UpdateAuthorByName(author string, name string) (int, error)
}
```
### Remove type parameter for support other go version?
- generic looks not necessary enough in this library, while we use comment to determine the interface is repository or not, while not all user like generics, it does imporve the readability.
without type paramter, we need to go through method list to determine the entity, but unfortunately, it is likely that all method doesn't bind entity, so that we can't get column name from entity's tag.
