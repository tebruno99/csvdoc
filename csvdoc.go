// Package csvdoc is a module that provides tools for reading documents and converting csv lines into struct fields via tag names.
package csvdoc

import (
	"reflect"
	"strings"
)

// buildReflectTagIndexCache builds a map[string]int of the field tag name and field index so this does not have to
// be completed on every Read().
func buildReflectTagIndexCache[T any](forWrite bool) (map[string]int, error) {
	// Precalculate the struct field indexes
	rt := new(T)
	ft := reflect.TypeOf(rt).Elem()

	fieldIndexes := make(map[string]int, ft.NumField())

	for i := range ft.NumField() {
		csvTag := ft.Field(i).Tag
		if _, ok := csvTag.Lookup("csv"); !ok {
			continue
		}
		tag := csvTag.Get("csv")
		readTag := strings.Split(tag, ",")
		if forWrite {
			if len(readTag) > 1 {
				readTag = readTag[1:]
			}
		}
		if _, tok := fieldIndexes[readTag[0]]; tok {
			return nil, ErrStructTagDuplicate
		}
		if readTag[0] != "" && readTag[0] != "-" {
			fieldIndexes[readTag[0]] = i
		}
	}

	return fieldIndexes, nil
}
