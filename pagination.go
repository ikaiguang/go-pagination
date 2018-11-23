package pagination

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

// pagination mode
const (
	PagingModePageNumber int64 = 1 // pagination mode : page number
	PagingModeCursor     int64 = 2 // pagination mode : cursor
)

const (
	defaultCurrentPageNumber int64   = 0                 // current page number : which page (default : 1)
	defaultGotoPageNumber    int64   = 1                 // goto page number : which page (default : 1)
	defaultPageSize          int64   = 15                // show records number (default : 15)
	defaultCursorValue       float64 = 0                 // cursor value (default : 0)
	defaultWherePlaceholder          = "?"               // where param placeholder
	defaultOrderRegexp               = "^[a-zA-Z0-9_]+$" // regexp
	defaultOrderColumn               = "id"              // default order column
	defaultOrderDirection            = "desc"            // default order direction
	defaultOrderAsc                  = "asc"             // order direction : asc
	defaultOrderDesc                 = "desc"            // order direction : desc
)

// PagingWhere : paging where (example : where id = ? => where id = 1)
type PagingWhere struct {
	Column      string      // where column  (default : id)
	Symbol      string      // where symbol ( (default : =)
	Placeholder string      // where placeholder  (default : ?)
	Data        interface{} // where data (default : interface{})
}

// PagingOptionCollection : paging option collection
type PagingOptionCollection struct {
	Option    *PagingOption  // option
	Limit     int64          // limit
	Offset    int64          // offset
	Where     []*PagingWhere // where
	Order     []*PagingOrder // order
	IsReverse bool           // cursor mode order by reverse
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
//				ListPointer:  &[]User{},
//			}
//
//			var collection2 = &PagingResultCollection{
//				TotalRecords: 15,
//				ListPointer:  &[]*User{},
//			}
//
//
//
// PagingResultCollection : paging result collection
type PagingResultCollection struct {
	TotalRecords int64       // total records
	ListPointer  interface{} // list pointer(example:[]struct{} or []*struct{})
}

// DefaultPagingOption : default paging option
func DefaultPagingOption() *PagingOption {
	return &PagingOption{
		PagingMode:        PagingModePageNumber,
		CurrentPageNumber: defaultCurrentPageNumber,
		GotoPageNumber:    defaultGotoPageNumber,
		PageSize:          defaultPageSize,
		OrderBy:           []*PagingOrder{},
		CursorColumn:      defaultOrderColumn,
		CursorDirection:   defaultOrderDirection,
		CursorValue:       defaultCursorValue,
	}
}

// InitPagingOption : init paging option
func InitPagingOption(pagingOption *PagingOption) error {
	// nil pointer
	if pagingOption == nil {
		return fmt.Errorf("PagingOption cannot be a nil pointer")
	}

	// init
	initPagingOption(pagingOption)

	return nil
}

// init paging option
func initPagingOption(pagingOption *PagingOption) {
	// paging mode
	pagingOption.PagingMode = pagingMode(pagingOption.PagingMode)

	// which page
	pagingOption.CurrentPageNumber = pagingCurrentPageNumber(pagingOption.CurrentPageNumber)

	// which page
	pagingOption.GotoPageNumber = pagingGotoPageNumber(pagingOption.GotoPageNumber)

	// page size
	pagingOption.PageSize = pagingPageSize(pagingOption.PageSize)

	// cursor column
	pagingOption.CursorColumn = pagingOrderColumn(pagingOption.CursorColumn)

	// cursor direction : asc or desc
	pagingOption.CursorDirection = pagingOrderDirection(pagingOption.CursorDirection)

	// order by
	//if pagingOption.OrderBy == nil {
	//	pagingOption.OrderBy = []*PagingOrder{}
	//}
}

// paging mode
func pagingMode(pagingMode int64) int64 {

	switch pagingMode {

	case PagingModeCursor: // paging mode cursor
		pagingMode = PagingModeCursor

	default:
		pagingMode = PagingModePageNumber
	}
	return pagingMode
}

// paging current page number
func pagingCurrentPageNumber(currentPageNumber int64) int64 {

	if currentPageNumber < 0 {
		return 0
	}
	return currentPageNumber
}

// paging goto page number
func pagingGotoPageNumber(gotoPageNumber int64) int64 {

	if gotoPageNumber < 1 {
		return 1
	}
	return gotoPageNumber
}

// paging size
func pagingPageSize(pageSize int64) int64 {

	if pageSize < 1 {
		return defaultPageSize
	}
	return pageSize
}

// paging order column
func pagingOrderColumn(column string) string {

	column = strings.TrimSpace(column)

	// secure column or default id
	if column == "" || !IsSecureSQLColumn(column) {
		return defaultOrderColumn
	}
	return column
}

// paging order direction
func pagingOrderDirection(direction string) string {

	direction = strings.TrimSpace(direction)

	// asc || desc
	if direction != defaultOrderAsc {
		direction = defaultOrderDesc
	}
	return direction
}

// IsSecureSQLColumn : is secure sql column
func IsSecureSQLColumn(value string) bool {
	return DefaultSecureSQLHandler(value)
}

// DefaultSecureSQLHandler : secure sql column
var DefaultSecureSQLHandler = func(value string) bool {
	return regexp.MustCompile(defaultOrderRegexp).MatchString(value)
}

// GetPagingOptionCollection : get paging option collection
func GetPagingOptionCollection(pagingOption *PagingOption, models ...interface{}) (*PagingOptionCollection, error) {

	// init paging option
	if pagingOption == nil {
		pagingOption = DefaultPagingOption()
	} else {
		initPagingOption(pagingOption)
	}

	switch pagingOption.PagingMode {

	case PagingModeCursor:
		return cursorPagingOptionCollection(pagingOption, models...)

	default:
		return pageNumberPagingOptionCollection(pagingOption), nil
	}
}

// cursor mode query option collection
func cursorPagingOptionCollection(pagingOption *PagingOption, models ...interface{}) (*PagingOptionCollection, error) {
	// check cursor column
	if err := DefaultCursorColumnCheckHandler(pagingOption, models...); err != nil {
		return nil, err
	}
	return DefaultCursorPagingOptionCollectionHandler(pagingOption), nil
}

// DefaultCursorColumnCheckHandler : cursor mode : check cursor column
var DefaultCursorColumnCheckHandler = func(pagingOption *PagingOption, models ...interface{}) error {

	if len(models) == 0 {
		return nil
	}

	// check cursor column
	exist, err := DefaultCursorColumnHandler(pagingOption, models[0])
	if err != nil {
		return err
	}

	// not exist
	if !exist {
		return fmt.Errorf("cursorColumn(%s) not exist in model(table)", pagingOption.CursorColumn)
	}
	return nil
}

// DefaultCursorColumnHandler :
// paging cursor column exist in model struct（structPointer : must be a valid struct pointer）
var DefaultCursorColumnHandler = func(pagingOption *PagingOption, structPointer interface{}) (bool, error) {

	fieldName := ToCamelString(pagingOption.CursorColumn)

	return FieldInStruct(structPointer, fieldName)
}

// DefaultCursorPagingOptionCollectionHandler :
// cursor mode option collection
//
// now have a table(tb_goods) : it has 200 records, and the auto_id is 1,2,3,4,5...200
// set each page show 10 records : PageSize = 10
//
//
//
// # example : order by auto_id asc
//
//		* directly jump to tenth page(CursorValue = 100)
// 			SELECT * FROM tb_goods ORDER BY auto_id ASC LIMIT 10 OFFSET 90
//
// 		* tenth page jump to the eleventh page(next page)
// 			SELECT * FROM tb_goods WHERE auto_id > 100 ORDER BY auto_id ASC LIMIT 10 OFFSET 0
//
// 		* tenth page jump to the twelfth page(next page)
// 			SELECT * FROM tb_goods WHERE auto_id > 100 ORDER BY auto_id ASC LIMIT 10 OFFSET 10
//
// 		* tenth page jump to the ninth page(preceding page)
// 			SELECT * FROM tb_goods WHERE auto_id <= 100 ORDER BY auto_id DESC LIMIT 10 OFFSET 10
//
// 		* tenth page jump to the eighth page(preceding page)
// 			SELECT * FROM tb_goods WHERE auto_id <= 100 ORDER BY auto_id DESC LIMIT 10 OFFSET 20
//
//
//
// # example : order by auto_id desc
//
//		* directly jump to tenth page(CursorValue = 101)
// 			SELECT * FROM tb_goods ORDER BY auto_id DESC LIMIT 10 OFFSET 90
//
// 		* tenth page jump to the eleventh page(next page)
// 			SELECT * FROM tb_goods WHERE auto_id < 101 ORDER BY auto_id DESC LIMIT 10 OFFSET 0
//
// 		* tenth page jump to the twelfth page(next page)
// 			SELECT * FROM tb_goods WHERE auto_id < 101 ORDER BY auto_id DESC LIMIT 10 OFFSET 10
//
// 		* tenth page jump to the ninth page(preceding page)
// 			SELECT * FROM tb_goods WHERE auto_id >= 101 ORDER BY auto_id ASC LIMIT 10 OFFSET 10
//
// 		* tenth page jump to the eighth page(preceding page)
// 			SELECT * FROM tb_goods WHERE auto_id >= 101 ORDER BY auto_id ASC LIMIT 10 OFFSET 20
//
// cursor mode option collection
var DefaultCursorPagingOptionCollectionHandler = func(pagingOption *PagingOption) *PagingOptionCollection {

	// init cursor query option collection
	pageSize := pagingOption.PageSize

	collection := &PagingOptionCollection{
		Option: pagingOption,
		Limit:  pageSize,
		Offset: 0,
		Where:  []*PagingWhere{},
		Order:  []*PagingOrder{},
	}

	// jump page
	currentPage := pagingOption.CurrentPageNumber
	gotoPage := pagingOption.GotoPageNumber
	jumpNumber := gotoPage - currentPage

	// order by
	order := &PagingOrder{
		Column:    pagingOrderColumn(pagingOption.CursorColumn),
		Direction: pagingOrderDirection(pagingOption.CursorDirection),
	}

	// where
	where := &PagingWhere{
		Column:      pagingOption.CursorColumn,
		Symbol:      "=",
		Placeholder: defaultWherePlaceholder,
		Data:        pagingOption.CursorValue,
	}

	// offset && where.Symbol
	switch order.Direction {

	case defaultOrderAsc: // asc

		switch {

		case currentPage == 0: // first page
			// offset
			collection.Offset = (gotoPage - 1) * pageSize
			if collection.Offset < 0 {
				collection.Offset = 0
			}
			// order
			order.Direction = defaultOrderAsc
			collection.Order = append(collection.Order, order)

		case jumpNumber < 0: // preceding page
			// where
			where.Symbol = "<="
			collection.Where = append(collection.Where, where)
			// order
			order.Direction = defaultOrderDesc
			collection.Order = append(collection.Order, order)
			collection.IsReverse = true
			// offset
			collection.Offset = (-jumpNumber) * pageSize
			if collection.Offset < 0 {
				collection.Offset = 0
			}

		default: // next page || page not change(page not change will be jump to next page)
			// where
			where.Symbol = ">"
			collection.Where = append(collection.Where, where)
			// order
			order.Direction = defaultOrderAsc
			collection.Order = append(collection.Order, order)
			// offset
			collection.Offset = (jumpNumber - 1) * pageSize
			if collection.Offset < 0 {
				collection.Offset = 0
			}
		}

	default: // desc

		switch {

		case currentPage == 0: // first page
			// offset
			collection.Offset = (gotoPage - 1) * pageSize
			if collection.Offset < 0 {
				collection.Offset = 0
			}
			// order
			order.Direction = defaultOrderDesc
			collection.Order = append(collection.Order, order)

		case jumpNumber < 0: // preceding page
			// where
			where.Symbol = ">="
			collection.Where = append(collection.Where, where)
			// order
			order.Direction = defaultOrderAsc
			collection.Order = append(collection.Order, order)
			collection.IsReverse = true
			// offset
			collection.Offset = (-jumpNumber) * pageSize
			if collection.Offset < 0 {
				collection.Offset = 0
			}

		default: // next page || page not change(page not change will be jump to next page)
			// where
			where.Symbol = "<"
			collection.Where = append(collection.Where, where)
			// order
			order.Direction = defaultOrderDesc
			collection.Order = append(collection.Order, order)
			// offset
			collection.Offset = (jumpNumber - 1) * pageSize
			if collection.Offset < 0 {
				collection.Offset = 0
			}
		}
	}
	return collection
}

// AnotherCursorPagingOptionCollectionHandler :
// cursor mode option collection
//
// now have a table(tb_goods) : it has 200 records, and the auto_id is 1,2,3,4,5...200
// set each page show 10 records : PageSize = 10
//
//
//
// # example : order by auto_id asc
//
//		* directly jump to tenth page(CursorValue = 100)
// 			SELECT * FROM tb_goods ORDER BY auto_id ASC LIMIT 10 OFFSET 90
//
// 		* tenth page jump to the eleventh page(next page)
// 			SELECT * FROM tb_goods WHERE auto_id > 100 ORDER BY auto_id ASC LIMIT 10 OFFSET 0
//
// 		* tenth page jump to the twelfth page(next page)
// 			SELECT * FROM tb_goods WHERE auto_id > 100 ORDER BY auto_id ASC LIMIT 10 OFFSET 10
//
// 		* tenth page jump to the ninth page(preceding page)
// 			SELECT * FROM tb_goods WHERE auto_id <= 100 ORDER BY auto_id ASC LIMIT 10 OFFSET 80
//
// 		* tenth page jump to the eighth page(preceding page)
// 			SELECT * FROM tb_goods WHERE auto_id <= 100 ORDER BY auto_id ASC LIMIT 10 OFFSET 70
//
//
//
// # example : order by auto_id desc
//
//		* directly jump to tenth page(CursorValue = 101)
// 			SELECT * FROM tb_goods ORDER BY auto_id DESC LIMIT 10 OFFSET 90
//
// 		* tenth page jump to the eleventh page(next page)
// 			SELECT * FROM tb_goods WHERE auto_id < 101 ORDER BY auto_id DESC LIMIT 10 OFFSET 0
//
// 		* tenth page jump to the twelfth page(next page)
// 			SELECT * FROM tb_goods WHERE auto_id < 101 ORDER BY auto_id DESC LIMIT 10 OFFSET 10
//
// 		* tenth page jump to the ninth page(preceding page)
// 			SELECT * FROM tb_goods WHERE auto_id >= 101 ORDER BY auto_id DESC LIMIT 10 OFFSET 80
//
// 		* tenth page jump to the eighth page(preceding page)
// 			SELECT * FROM tb_goods WHERE auto_id >= 101 ORDER BY auto_id DESC LIMIT 10 OFFSET 70
//
// cursor mode option collection
var AnotherCursorPagingOptionCollectionHandler = func(pagingOption *PagingOption) *PagingOptionCollection {

	// init cursor query option collection
	pageSize := pagingOption.PageSize

	collection := &PagingOptionCollection{
		Option: pagingOption,
		Limit:  pageSize,
		Offset: 0,
		Where:  []*PagingWhere{},
		Order:  []*PagingOrder{},
	}

	// jump page
	currentPage := pagingOption.CurrentPageNumber
	gotoPage := pagingOption.GotoPageNumber
	jumpNumber := gotoPage - currentPage

	// order by
	order := &PagingOrder{
		Column:    pagingOrderColumn(pagingOption.CursorColumn),
		Direction: pagingOrderDirection(pagingOption.CursorDirection),
	}

	// where
	where := &PagingWhere{
		Column:      pagingOption.CursorColumn,
		Symbol:      "=",
		Placeholder: defaultWherePlaceholder,
		Data:        pagingOption.CursorValue,
	}

	// offset && where.Symbol
	switch order.Direction {

	case defaultOrderAsc: // asc

		switch {

		case currentPage == 0: // first page
			// offset
			collection.Offset = (gotoPage - 1) * pageSize
			if collection.Offset < 0 {
				collection.Offset = 0
			}
			// order
			order.Direction = defaultOrderAsc
			collection.Order = append(collection.Order, order)

		case jumpNumber < 0: // preceding page
			// where
			where.Symbol = "<="
			collection.Where = append(collection.Where, where)
			// order
			order.Direction = defaultOrderAsc
			collection.Order = append(collection.Order, order)
			// offset
			collection.Offset = (gotoPage - 1) * pageSize
			if collection.Offset < 0 {
				collection.Offset = 0
			}

		default: // next page || page not change(page not change will be jump to next page)
			// where
			where.Symbol = ">"
			collection.Where = append(collection.Where, where)
			// order
			order.Direction = defaultOrderAsc
			collection.Order = append(collection.Order, order)
			// offset
			collection.Offset = (jumpNumber - 1) * pageSize
			if collection.Offset < 0 {
				collection.Offset = 0
			}
		}

	default: // desc

		switch {

		case currentPage == 0: // first page
			// offset
			collection.Offset = (gotoPage - 1) * pageSize
			if collection.Offset < 0 {
				collection.Offset = 0
			}
			// order
			order.Direction = defaultOrderDesc
			collection.Order = append(collection.Order, order)

		case jumpNumber < 0: // preceding page
			// where
			where.Symbol = ">="
			collection.Where = append(collection.Where, where)
			// order
			order.Direction = defaultOrderDesc
			collection.Order = append(collection.Order, order)
			// offset
			collection.Offset = (gotoPage - 1) * pageSize
			if collection.Offset < 0 {
				collection.Offset = 0
			}

		default: // next page || page not change(page not change will be jump to next page)
			// where
			where.Symbol = "<"
			collection.Where = append(collection.Where, where)
			// order
			order.Direction = defaultOrderDesc
			collection.Order = append(collection.Order, order)
			// offset
			collection.Offset = (jumpNumber - 1) * pageSize
			if collection.Offset < 0 {
				collection.Offset = 0
			}
		}
	}
	return collection
}

// page number mode option collection
//
// now have a table(tb_goods) : it has 200 records, and the auto_id is 1,2,3,4,5...200
// set each page show 10 records : PageSize = 10
//
//
//
// # example : order by auto_id asc
//
//		* directly jump to tenth page
// 			SELECT * FROM tb_goods ORDER BY auto_id ASC LIMIT 10 OFFSET 90
//
// 		* tenth page jump to the eleventh page(next page)
// 			SELECT * FROM tb_goods ORDER BY auto_id ASC LIMIT 10 OFFSET 100
//
// 		* tenth page jump to the twelfth page(next page)
// 			SELECT * FROM tb_goods ORDER BY auto_id ASC LIMIT 10 OFFSET 110
//
// 		* tenth page jump to the ninth page(preceding page)
// 			SELECT * FROM tb_goods ORDER BY auto_id ASC LIMIT 10 OFFSET 80
//
// 		* tenth page jump to the eighth page(preceding page)
// 			SELECT * FROM tb_goods ORDER BY auto_id ASC LIMIT 10 OFFSET 70
//
//
//
// # example : order by auto_id desc
//
//		* directly jump to tenth page
// 			SELECT * FROM tb_goods ORDER BY auto_id DESC LIMIT 10 OFFSET 90
//
// 		* tenth page jump to the eleventh page(next page)
// 			SELECT * FROM tb_goods ORDER BY auto_id DESC LIMIT 10 OFFSET 100
//
// 		* tenth page jump to the twelfth page(next page)
// 			SELECT * FROM tb_goods ORDER BY auto_id DESC LIMIT 10 OFFSET 110
//
// 		* tenth page jump to the ninth page(preceding page)
// 			SELECT * FROM tb_goods ORDER BY auto_id DESC LIMIT 10 OFFSET 80
//
// 		* tenth page jump to the eighth page(preceding page)
// 			SELECT * FROM tb_goods ORDER BY auto_id DESC LIMIT 10 OFFSET 70
//
// page number mode option collection
func pageNumberPagingOptionCollection(pagingOption *PagingOption) *PagingOptionCollection {

	limit := pagingOption.PageSize
	offset := pagingOption.PageSize * (pagingOption.GotoPageNumber - 1)
	orderSlice := DefaultPageNumberOrderHandler(pagingOption.OrderBy)

	collection := &PagingOptionCollection{
		Option: pagingOption,
		Limit:  limit,
		Offset: offset,
		Order:  orderSlice,
	}
	return collection
}

// DefaultPageNumberOrderHandler : order
var DefaultPageNumberOrderHandler = func(orderOptions []*PagingOrder) []*PagingOrder {

	var queryOrder []*PagingOrder

	// custom order
	for _, orderBy := range orderOptions {
		queryOrder = append(queryOrder, &PagingOrder{
			Column:    pagingOrderColumn(orderBy.Column),
			Direction: pagingOrderDirection(orderBy.Direction),
		})
	}
	return queryOrder
}

// SetPagingResult : set paging result
func SetPagingResult(optionCollection *PagingOptionCollection, resultCollection *PagingResultCollection) (*PagingResult, error) {

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

	// ListPointer not pointer
	sliceValue := reflect.ValueOf(resultCollection.ListPointer)
	if sliceValue.Kind() != reflect.Ptr {
		err := fmt.Errorf("resultCollection.ListPointer is not a valid slice pointer")
		return pagingResult, err
	}

	// ListPointer not slice
	sliceElem := sliceValue.Elem()
	if sliceElem.Kind() != reflect.Slice {
		err := fmt.Errorf("resultCollection.ListPointer is not a valid slice pointer")
		return pagingResult, err
	}

	// ListPointer slice size
	sliceLen := sliceElem.Len()

	// empty list
	if sliceLen == 0 {
		return pagingResult, nil
	}

	// show from - to
	if sliceLen > 0 {
		pagingResult.ShowFrom = (pagingResult.CurrentPage-1)*pagingOption.PageSize + 1
		pagingResult.ShowTo = pagingResult.ShowFrom + int64(sliceLen) - 1
	}

	// cursor mode .
	// goto preceding page .
	// paging option cursor direction(PagingOption.CursorDirection)
	// and sql order by(OptionCollection.Order.Direction) are the opposite .
	// so reverse the slice(list) .
	// keep data sort same as paging option cursor direction
	if err := DefaultCursorPagingResultHandler(optionCollection, resultCollection); err != nil {
		return pagingResult, err
	}

	// cursor value
	cursorValue, err := DefaultCursorValueHandler(optionCollection, sliceElem.Index(sliceLen - 1).Interface())
	if err != nil {
		return pagingResult, nil
	}
	pagingResult.CursorValue = cursorValue

	return pagingResult, nil
}

// DefaultCursorPagingResultHandler :
// serialization cursor mode paging result
var DefaultCursorPagingResultHandler = func(optionCollection *PagingOptionCollection, resultCollection *PagingResultCollection) error {
	// not reverse
	if !optionCollection.IsReverse {
		return nil
	}

	// reverse slice
	return ReverseSlice(resultCollection.ListPointer)
}

// DefaultCursorValueHandler :
// new cursor value（structPointer : must be a valid struct pointer）
var DefaultCursorValueHandler = func(optionCollection *PagingOptionCollection, structData interface{}) (cursorValue float64, err error) {

	// struct elem
	structElem := reflect.ValueOf(structData)
	if structElem.Kind() == reflect.Ptr {
		structElem = structElem.Elem()
	}

	// not struct
	if structElem.Kind() != reflect.Struct {
		err := fmt.Errorf("structData is not a valid struct")
		return cursorValue, err
	}

	// column name
	columnName := ToCamelString(optionCollection.Option.CursorColumn)
	// column exist in struct
	if !structElem.FieldByName(columnName).IsValid() {
		err = fmt.Errorf("cursor column not exist in struct")
		return cursorValue, err
	}

	// column value
	columnValue := structElem.FieldByName(columnName)

	switch columnValue.Kind() {

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		cursorValue = float64(columnValue.Int())

	case reflect.Float32:
		cursorValue, err = strconv.ParseFloat(fmt.Sprint(columnValue.Interface().(float32)), 64)
		if err != nil {
			err = fmt.Errorf("cursor column value : float32 convert to float64 fail : %v", err)
			return cursorValue, err
		}

	case reflect.Float64:
		cursorValue = columnValue.Float()

	case reflect.String:
		cursorValue, err = strconv.ParseFloat(columnValue.String(), 64)
		if err != nil {
			err = fmt.Errorf("cursor column value : string convert to number fail : %v", err)
			return cursorValue, err
		}

	default:
		err = fmt.Errorf("cursor column value not a number or number string")
		cursorValue = 0
		return cursorValue, err
	}
	return cursorValue, err
}
