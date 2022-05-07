# jorm
an orm library for go, based on go generate, zero document needed, inspired by jpa


### Agenda 

- below functions list what is going to be implented or not
```go
type BookRepository[T entity.Book] interface {
    //done
    FindByNameAndAuthor(name string, author string) (T, error)
    //needed? this signature is not good for readness
    FindByNameAndAuthor(name, author string) (T, error) 
    //in progress
    FindById(k int) (T,error)
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
## Remove type parameter for support other go version?
generic looks not necessary enough in this library, while we use comment to determine the interface is repository or not, while not all user like generics, it does imporve the readability.
without type paramter, we need to go through method list to determine the entity, but unfortunately, it is likely that all method doesn't bind entity, so that we can't get column name from entity's tag.

## Do not limit the method name
instead of using method name to extract the sql criteria, looks for comment to extract the criteria is more graceful way
long method name is ugly if there are many paramters, the long name is a pain to IDE autocompletion as the popup windows is small.
we want to display the whole method signature

## Reduce the times of concating strings
FindByNameInAndAuthorAndCreatedLessThan
select * from book where name in (?,?,?,?) and author = ? and Created < ?
author = ? and Created < ? can be in one string literal rather than concating auther then concating created

## as many as possible output 
for the generated code 
print to console
save to one file
save to directory for many file

## code generation
able to generate for one method
able to generate for one source file
able to generate for whole package
## as many as possible init config
able to customize the annotation in order to compat other orm libs like use gorm instead of jorm-column 
able to generate without repostiory (maybe console only as file is better seperated by entity or repostiroy)