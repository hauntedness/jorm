package entity

//jorm-entity:"true"
//jorm-table:"book"
type Book struct {
	Id      int    `jorm-column:"id"`
	Name    string `jorm-column:"name"`
	Author  string `jorm-column:"author"`
	Version int64  `jorm-column:"version"`
}

/*@Entity(table="book")*/
type Book2 struct {
	Id      int    `column:"id"`
	Name    string `column:"name"`
	Author  string `column:"author"`
	Version int64  `column:"version"`
}
