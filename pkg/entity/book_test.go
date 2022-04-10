package entity

import (
	"testing"
)

func BenchmarkNewBook(b *testing.B) {
	for i := 0; i < 1000000; i++ {
		var b Book
		b.Author = "a"
		b.Name = "N"
		b.Id = 1
		b.Version = 1
	}
}

//1000000000               0.0000223 ns/op               0 B/op          0 allocs/op
func BenchmarkNewBook2(b *testing.B) {
	for i := 0; i < 1000000; i++ {
		var b Book2
		b.Author = "a"
		b.Name = "N"
		b.Id = 1
		b.Version = 1
	}
}
