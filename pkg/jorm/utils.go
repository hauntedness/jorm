package jorm

import (
	"go/ast"
	"strings"
	"unicode"
)

func ExtractTagValue(field *ast.Field, tagkey string) (column string, ok bool) {
	return extract(field.Tag.Value, tagkey)
}

//`column:   " ve   rsion"`
func extract(text string, tagkey string) (value string, ok bool) {
	ok = false
	value = ""
	c := []rune{}
	sub := []rune(tagkey)
	sub = append(sub, ':', '"')
	index := 0
	for _, r := range text {
		if r == ' ' {
			continue
		}
		if index == len(sub) {
			if r == '"' {
				value = string(c)
				ok = true
				return
			}
			c = append(c, r)
		} else if r == sub[index] {
			index++
		} else {
			index = 0
		}
	}
	return
}

// id -> ID
// BookName -> BOOK_NAME
func CaseTitleToSnake(text string) string {
	runes := []rune(text)
	length := len(runes)
	words := make([]rune, 0, 2*length)
	for i := 0; i < length; i++ {
		if unicode.IsUpper(runes[i]) && i > 0 {
			words = append(words, rune('_'), runes[i])
		} else {
			words = append(words, unicode.ToUpper(runes[i]))
		}
	}
	return string(words)
}

// Id -> id
// BookName -> bookName
func CaseTitleToCamel(text string) string {
	runes := []rune(text)
	length := len(runes)
	if length == 0 {
		return text
	} else if length == 1 && unicode.IsUpper(runes[0]) {
		return strings.ToLower(text)
	} else if length > 1 && unicode.IsUpper(runes[0]) {
		runes[0] = unicode.ToLower(runes[0])
		return string(runes)
	} else {
		return text
	}
}

func ExtractPackageNameFromPath(path string) string {
	before, _, found := strings.Cut(path, "::")
	if found {
		splits := strings.Split(before, `/`)
		if len(splits) > 0 {
			return splits[len(splits)-1]
		} else {
			panic("can't extract name from " + path)
		}
	} else {
		panic("can't extract name from " + path)
	}
}
