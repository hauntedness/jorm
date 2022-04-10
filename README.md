# jorm
an orm library for go, based on go generate, zero document needed, inspired by jpa


### Agenda 

- below functions list what is going to be implented or not
```go
type BookRepository[T entity.Book, K int] interface {
    // done
    func FindByNameAndAuthor(name string, author string) (T, error)
    // needed? this signature is not good for readness
    func FindByNameAndAuthor(name, author string) (T, error) 
    // in progress
    func FindById(k K) (T,error)
    //in progress
    func FindAuthorByName(name string) (string , error)
    //in progress
    func FindAllByName(name string) ([]T , error)
    //in progress
    func FindAllAuthorByName(name string) ([]string, error)
    //in progress
    func FindByCreatedDateLessThan(ptime time.Time) ([]T, error) 
        //in progress
    func FindByCreatedDateBetween(start time.Time, end time.Time) ([]T, error) 
    //in progress 
    // batch selection is not worth here
    func FindByNameIn(names []string) ([]T, error) 
    //in progress 
    // batch selection is not worth here
    func FindByNameNotIn(names []string) ([]T, error)
    //in progress will not provide auto generated id
    func Insert(book T) error
      //in progress will not provide auto generated id
    func InsertAll(books []T) (int, error)
    // in progress
    // same to update by id?
    func Update(book T) error
        // in progress
    func UpdateAll(books []T) (int, error)
    // to keep the similar order to update author = ? where name = ?
    // needn't to return the number of updated rows?
    // in progress
    func UpdateAuthorByName(author string, name string) (int, error)
    // in progress
    // very risky operation, all column same as each other only except id?
    // func UpdateByName(book T, name string) error
}
```
### Remove type parameter for support other go version?
- generic looks not necessary enough in this library, while we use comment to determine the interface is repository or not, while not all user like generics, it does imporve the readability.
without type paramter, we need to go through method list to determine the entity, but unfortunately, it is likely that all method doesn't bind entity, so that we can't get column name from entity's tag.
