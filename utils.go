package pagination

import (
	"fmt"
	"reflect"
)

// ToCamelString : camel string, xx_yy to XxYy
func ToCamelString(s string) string {
	data := make([]byte, 0, len(s))
	j := false
	k := false
	num := len(s) - 1
	for i := 0; i <= num; i++ {
		d := s[i]
		if k == false && d >= 'A' && d <= 'Z' {
			k = true
		}
		if d >= 'a' && d <= 'z' && (j || k == false) {
			d = d - 32
			j = false
			k = true
		}
		if k && d == '_' && num > i && s[i+1] >= 'a' && s[i+1] <= 'z' {
			j = true
			continue
		}
		data = append(data, d)
	}
	return string(data[:])
}

// FieldInStruct : field in struct（structPointer : must be a valid struct pointer）
func FieldInStruct(structPointer interface{}, field string) (bool, error) {
	// not pointer
	reflectValue := reflect.ValueOf(structPointer)
	if reflectValue.Kind() != reflect.Ptr {
		err := fmt.Errorf("structPointer is not a valid struct pointer")
		return false, err
	}

	// not struct
	structElem := reflectValue.Elem()
	if structElem.Kind() != reflect.Struct {
		err := fmt.Errorf("structPointer is not a valid struct pointer")
		return false, err
	}
	return structElem.FieldByName(field).IsValid(), nil
}

// ReverseSlice : reverse slice（slicePointer : must be a valid slice pointer）
func ReverseSlice(slicePointer interface{}) error {
	// not pointer
	reflectValue := reflect.ValueOf(slicePointer)
	if reflectValue.Kind() != reflect.Ptr {
		err := fmt.Errorf("slicePointer is not a valid slice pointer")
		return err
	}

	// not slice
	sliceElem := reflectValue.Elem()
	if sliceElem.Kind() != reflect.Slice {
		err := fmt.Errorf("slicePointer is not a valid slice pointer")
		return err
	}

	// slice size
	size := sliceElem.Len()
	if size <= 1 {
		return nil
	}

	// swap
	swap := reflect.Swapper(sliceElem.Interface())
	for i, j := 0, size-1; i < j; i, j = i+1, j-1 {
		swap(i, j)
	}
	return nil
}
