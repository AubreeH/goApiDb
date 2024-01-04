package structParsing

import (
	"reflect"
)

type Tag string

func GetTag(field reflect.StructField, tag Tag) string {
	return field.Tag.Get(string(tag))
}

func (t Tag) Get(field reflect.StructField) string {
	return GetTag(field, t)
}

func (t Tag) Lookup(field reflect.StructField) (string, bool) {
	return field.Tag.Lookup(string(t))
}
