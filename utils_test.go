package pagination

import (
	"fmt"
	"testing"
)

// snake string to camel string
func TestToCamelString(t *testing.T) {
	camelString := "XxYy"
	snakeString := "xx_yy"

	got := ToCamelString(snakeString)
	if got != camelString {
		t.Errorf("\n testing : ToCamelString error : ToCamelString(snakeString) != camelString\n")
	} else {
		t.Logf("\n snakeString(%s) => camelString(%s)\n", snakeString, got)
	}
}

// field in struct
func TestFieldInStruct(t *testing.T) {

	type Pointer struct {
		Field string
	}

	var p Pointer

	field := "Field"
	got, err := FieldInStruct(&p, field)
	if err != nil {
		t.Errorf("\n testing : FieldInStruct error : %v \n", err)
	} else {
		t.Logf("\n exist(%v) : %s in %#v\n", got, field, &p)
	}
}

// reverse slice
func TestReverseSlice(t *testing.T) {
	type User struct {
		Name string
		Age  int
	}

	var users []User

	for i := 1; i <= 6; i++ {
		users = append(users, User{Name: fmt.Sprintf("name_%d", i), Age: i})
	}

	t.Logf("\n before ReverseSlice : %v \n", users)

	err := ReverseSlice(&users)
	if err != nil {
		t.Errorf("\n testing : ReverseSlice error : %v \n", err)
	} else {
		t.Logf("\n behind ReverseSlice : %v \n", users)
	}
}
