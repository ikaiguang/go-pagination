package pagination

// PagingWhere : paging where (example : where id = ? => where id = 1)
type PagingWhere struct {
	Column      string      // where column  (default : id)
	Symbol      string      // where symbol ( (default : =)
	Placeholder string      // where placeholder  (default : ?)
	Data        interface{} // where data (default : interface{})
}

// OptionCollection : paging option collection
type OptionCollection struct {
	Option    *PagingOption  // option
	Limit     int64          // limit
	Offset    int64          // offset
	Where     []*PagingWhere // where
	Order     []*PagingOrder // order
	IsReverse bool           // cursor mode order by reverse
}

// GetOptionCollection : get paging option collection
func GetOptionCollection(pagingOption *PagingOption, models ...interface{}) (*OptionCollection, error) {

	// init paging option
	if pagingOption == nil {
		pagingOption = DefaultPagingOption()
	} else {
		initPagingOption(pagingOption)
	}

	switch pagingOption.PagingMode {

	case PagingModeCursor:
		return getCursorOptionCollection(pagingOption, models...)

	default:
		return getNumberOptionCollection(pagingOption), nil
	}
}

// PagingResultCollection : paging result collection
// @Param ListPointer must be a valid slice pointer, and
// slice element must be a struct type
//
//
//
// example :
//			type User struct {
//				Name string
//				Age  int
//			}
//
//			var collection1 = &PagingResultCollection{
//				TotalRecords: 15,
//				ListPointer:  []User{},
//			}
//
//			var collection2 = &PagingResultCollection{
//				TotalRecords: 15,
//				ListPointer:  []*User{},
//			}
//
//			var collection2 = &PagingResultCollection{
//				TotalRecords: 15,
//				ListPointer:  &[]*User{},
//			}
//
// PagingResultCollection : paging result collection
type PagingResultCollection struct {
	TotalRecords int64       // total records
	ResultSlice  interface{} // slice(example:[]struct{} or []*struct{})
}

// SetPagingResult : set paging result
func SetPagingResult(optionCollection *OptionCollection, resultCollection *PagingResultCollection) (*PagingResult, error) {

	// paging option
	pagingOption := optionCollection.Option

	// paging result
	pagingResult := &PagingResult{
		PagingMode:      pagingOption.PagingMode,       // paging mode
		TotalSize:       resultCollection.TotalRecords, // total size
		PageSize:        pagingOption.PageSize,         // page size
		CurrentPage:     pagingOption.GotoPageNumber,   // current page
		ShowFrom:        0,                             // current page show from - to record
		ShowTo:          0,                             // current page show from - to record
		LastPage:        0,                             // last page
		OrderBy:         pagingOption.OrderBy,          // order
		CursorColumn:    pagingOption.CursorColumn,     // cursor column
		CursorDirection: pagingOption.CursorDirection,  // cursor direction
		CursorValue:     0,                             // cursor value
		Option:          pagingOption,                  // paging option
	}

	// empty records
	if resultCollection.TotalRecords <= 0 {
		return pagingResult, nil
	}

	// last page
	if pagingResult.TotalSize%pagingOption.PageSize == 0 {
		pagingResult.LastPage = resultCollection.TotalRecords / pagingOption.PageSize
	} else {
		pagingResult.LastPage = resultCollection.TotalRecords/pagingOption.PageSize + 1
	}

	// calc ResultSlice
	sliceInfo, err := DefaultCalcResultSliceHandler(optionCollection, resultCollection)
	if err != nil {
		return pagingResult, err
	}

	// CursorValue
	pagingResult.CursorValue = sliceInfo.CursorValue

	// empty slice
	if sliceInfo.SliceLen == 0 {
		return pagingResult, nil
	}

	// show from - to
	pagingResult.ShowFrom = (pagingResult.CurrentPage-1)*pagingOption.PageSize + 1
	pagingResult.ShowTo = pagingResult.ShowFrom + int64(sliceInfo.SliceLen) - 1

	return pagingResult, nil
}
