package pagination

import "testing"

// paging query option collection
func TestGetOptionCollection(t *testing.T) {
	option := DefaultPagingOption()

	got, err := GetOptionCollection(option)
	if err != nil {
		t.Errorf("\n testing : GetOptionCollection error : %v \n", err)
	} else {
		t.Logf("\n GetOptionCollection result : %v\n", got)
	}
}

// init paging option
func TestInitPagingOption(t *testing.T) {
	option := DefaultPagingOption()

	if err := InitPagingOption(option); err != nil {
		t.Errorf("\n testing : InitPagingOption error : %v \n", err)
		return
	}

	t.Logf("\n InitPagingOption %v\n", option)
}

// new cursor value
func TestNewCursorValue(t *testing.T) {
	type Model struct {
		ID          string
		ValueColumn int32
	}

	var model *Model
	model = &Model{ID: "1.23456789", ValueColumn: 1234}

	collection := &PagingOptionCollection{}
	collection.Option = DefaultPagingOption()
	collection.Option.CursorColumn = "ID"
	//collection.Option.CursorColumn = "value_column"

	got, err := DefaultCursorValueHandler(collection, model)
	if err != nil {
		t.Errorf("\n testing : DefaultCursorValueHandler error : %v \n", err)
	} else {
		t.Logf("\n DefaultCursorValueHandler result : %v\n", got)
	}
}

// cursor column exist in model
func TestCursorColumnExistInModel(t *testing.T) {
	type Model struct {
		ID int32
	}

	option := DefaultPagingOption()
	option.CursorColumn = "ID"
	model := new(Model)

	got, err := DefaultCursorColumnHandler(option, model)
	if err != nil {
		t.Errorf("\n testing : DefaultCursorColumnHandler error : %v \n", err)
	} else {
		t.Logf("\n DefaultCursorColumnHandler result : %v\n", got)
	}
}
