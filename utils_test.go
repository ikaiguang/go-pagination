package pagination

import (
	"testing"
)

// snake string to camel string
func TestToCamelString(t *testing.T) {
	camelString := "XxYy"
	snakeString := "xx_yy"

	got := StringToCamel(snakeString)
	if got != camelString {
		t.Errorf("\n testing : StringToCamel error : StringToCamel(snakeString) != camelString\n")
	} else {
		t.Logf("\n snakeString(%s) => camelString(%s)\n", snakeString, got)
	}
}
