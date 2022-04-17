package jorm

import (
	"fmt"
	"strconv"
	"testing"
)

func TestArray(t *testing.T) {
	var clause []any
	// literal A = ?
	var literal string = "hahaha"
	// list
	var list []any = []any{"aaa", 111, true, []string{"iii", "jjj"}}
	for i := 0; i < 10000; i++ {
		list = append(list, literal+strconv.Itoa(i))
	}
	clause = append(clause, list...)
	for i, v := range clause {
		fmt.Println("i:", i, "v:", v)
	}
}
