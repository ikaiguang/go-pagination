package pagination

import "testing"

// paging query option collection
func TestGetPagingOptionCollection(t *testing.T) {
	option := DefaultPagingOption()

	got, err := GetPagingOptionCollection(option)
	if err != nil {
		t.Errorf("\n testing : GetPagingOptionCollection error : %v \n", err)
	} else {
		t.Logf("\n GetPagingOptionCollection result : %v\n", got)
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

// is secure sql column
func TestIsSecureSqlColumn(t *testing.T) {
	sqlColumnSlice := []string{
		"abc",
		"a-b-c",
		"a_b_c",
		"a b c",
	}
	for _, sqlColumn := range sqlColumnSlice {
		got := IsSecureSQLColumn(sqlColumn)
		t.Logf("\n secure(%v) : %s\n", got, sqlColumn)
	}
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

// set paging result
func TestSetPagingResult(t *testing.T) {
	//option := DefaultPagingOption()
	//collection := &PagingResultCollection{}
	//
	//got, err := SetPagingResult(option, collection)
	//if err != nil {
	//	t.Errorf("\n testing : SetPagingResult error : %v \n", err)
	//	return
	//} else {
	//	t.Logf("\n SetPagingResult result : %v\n", got)
	//}
}
