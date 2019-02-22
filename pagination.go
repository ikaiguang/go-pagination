package pagination

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// pagination mode
const (
	PagingModeNumber int64 = 1 // pagination mode : number
	PagingModeCursor int64 = 2 // pagination mode : cursor
)

const (
	defaultCurrentPageNumber int64   = 0      // current page number : which page (default : 1)
	defaultGotoPageNumber    int64   = 1      // goto page number : which page (default : 1)
	defaultPageSize          int64   = 15     // show records number (default : 15)
	defaultCursorValue       float64 = 0      // cursor value (default : 0)
	defaultWherePlaceholder          = "?"    // where param placeholder
	defaultOrderColumn               = "id"   // default order column
	defaultOrderDirection            = "desc" // default order direction
	defaultOrderAsc                  = "asc"  // order direction : asc
	defaultOrderDesc                 = "desc" // order direction : desc
)

// DefaultPagingOption : default paging option
func DefaultPagingOption() *PagingOption {
	return &PagingOption{
		PagingMode:        PagingModeNumber,
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
	pagingOption.PagingMode = getPagingMode(pagingOption.PagingMode)

	// which page
	pagingOption.CurrentPageNumber = getCurrentPageNumber(pagingOption.CurrentPageNumber)

	// which page
	pagingOption.GotoPageNumber = getGotoPageNumber(pagingOption.GotoPageNumber)

	// page size
	pagingOption.PageSize = getPageSize(pagingOption.PageSize)

	// cursor column
	pagingOption.CursorColumn = getOrderColumn(pagingOption.CursorColumn)

	// cursor direction : asc or desc
	pagingOption.CursorDirection = getOrderDirection(pagingOption.CursorDirection)

	// order by
	//if pagingOption.OrderBy == nil {
	//	pagingOption.OrderBy = []*PagingOrder{}
	//}
}

// getPagingMode paging mode
func getPagingMode(pagingMode int64) int64 {

	switch pagingMode {

	case PagingModeCursor: // paging mode cursor
		return PagingModeCursor

	default:
		return PagingModeNumber
	}
}

// getCurrentPageNumber paging current page number
func getCurrentPageNumber(currentPageNumber int64) int64 {

	if currentPageNumber < 0 {
		return 0
	}
	return currentPageNumber
}

// getGotoPageNumber paging goto page number
func getGotoPageNumber(gotoPageNumber int64) int64 {

	if gotoPageNumber < 1 {
		return 1
	}
	return gotoPageNumber
}

// getPageSize paging size
func getPageSize(pageSize int64) int64 {

	if pageSize < 1 {
		return defaultPageSize
	}
	return pageSize
}

// getOrderColumn paging order column
func getOrderColumn(column string) string {
	column = strings.TrimSpace(column)

	if len(column) == 0 {
		return defaultOrderColumn
	}
	return column
}

// getOrderDirection paging order direction
func getOrderDirection(direction string) string {

	direction = strings.TrimSpace(direction)

	// asc || desc
	if direction != defaultOrderAsc {
		direction = defaultOrderDesc
	}
	return direction
}

//////////////////////////////////////////////////////////////////////////////////////////

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

// GetOptionCollection : get paging option collection
func GetOptionCollection(pagingOption *PagingOption, models ...interface{}) (*PagingOptionCollection, error) {

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

//////////////////////////////////////////////////////////////////////////////////////////

// getNumberOptionCollection page number mode option collection
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
// getNumberOptionCollection page number mode option collection
func getNumberOptionCollection(pagingOption *PagingOption) *PagingOptionCollection {

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
			Column:    getOrderColumn(orderBy.Column),
			Direction: getOrderDirection(orderBy.Direction),
		})
	}
	return queryOrder
}

//////////////////////////////////////////////////////////////////////////////////////////

// getCursorOptionCollection get cursor mode query option collection
func getCursorOptionCollection(pagingOption *PagingOption, models ...interface{}) (*PagingOptionCollection, error) {
	// check cursor column
	if err := DefaultCursorColumnCheckHandler(pagingOption, models...); err != nil {
		return nil, err
	}
	return DefaultCursorOptionCollectionHandler(pagingOption), nil
}

// DefaultCursorColumnCheckHandler : check cursor column
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

// DefaultCursorColumnHandler : model struct has field(cursor column)
// return error if s not a struct
var DefaultCursorColumnHandler = func(pagingOption *PagingOption, model interface{}) (bool, error) {

	fieldName := StringToCamel(pagingOption.CursorColumn)

	// reflect.Value
	modelValue := reflect.ValueOf(model)

	// is pointer
	if modelValue.Kind() == reflect.Ptr {
		modelValue = reflect.ValueOf(modelValue.Elem().Interface())
	}

	if modelValue.Kind() != reflect.Struct {
		return false, fmt.Errorf("model isnot struct")
	}
	return modelValue.FieldByName(fieldName).IsValid(), nil
}

// DefaultCursorOptionCollectionHandler :
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
var DefaultCursorOptionCollectionHandler = func(pagingOption *PagingOption) *PagingOptionCollection {

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
		Column:    getOrderColumn(pagingOption.CursorColumn),
		Direction: getOrderDirection(pagingOption.CursorDirection),
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

// AnotherCursorOptionCollectionHandler : cursor mode option collection
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
// AnotherCursorOptionCollectionHandler cursor mode option collection
var AnotherCursorOptionCollectionHandler = func(pagingOption *PagingOption) *PagingOptionCollection {

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
		Column:    getOrderColumn(pagingOption.CursorColumn),
		Direction: getOrderDirection(pagingOption.CursorDirection),
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

//////////////////////////////////////////////////////////////////////////////////////////

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

// PagingResultInfo  calc ResultSlice
type PagingResultInfo struct {
	SliceLen    int64
	CursorValue float64
}

// DefaultCalcResultSliceHandler calc ResultSlice
var DefaultCalcResultSliceHandler = func(optionCollection *PagingOptionCollection, resultCollection *PagingResultCollection) (*PagingResultInfo, error) {
	var res = new(PagingResultInfo)

	// ResultSlice interface
	var sInterface = resultCollection.ResultSlice
	var sReflectValue reflect.Value

	// ResultSlice Value
	sReflectValue = reflect.ValueOf(sInterface)

	// is pointer slice
	if sReflectValue.Kind() == reflect.Ptr {
		sInterface = sReflectValue.Elem().Interface()
		sReflectValue = reflect.ValueOf(sInterface)
	}

	// not slice
	if sReflectValue.Kind() != reflect.Slice {
		return nil, fmt.Errorf("ResultSlice not a slice")
	}

	// ResultSlice slice size
	sLen := sReflectValue.Len()
	res.SliceLen = int64(sLen)

	// empty list
	if sLen == 0 {
		return res, nil
	}

	// cursor mode .
	// goto preceding page .
	// paging option cursor direction(PagingOption.CursorDirection)
	// and sql order by(PagingOptionCollection.Order.Direction) is opposite .
	// so reverse the ResultSlice.
	// keep data sort same as paging option cursor direction
	// not reverse
	if optionCollection.IsReverse {
		swap := reflect.Swapper(sInterface)

		for i, j := 0, sReflectValue.Len()-1; i < j; i, j = i+1, j-1 {
			swap(i, j)
		}
	}

	// CursorValue
	cursorValue, err := DefaultCursorValueHandler(optionCollection, sReflectValue.Index(sLen - 1).Interface())
	if err != nil {
		return nil, err
	}
	res.CursorValue = cursorValue

	return res, nil
}

// DefaultCursorValueHandler : calc PagingResult.CursorValue
var DefaultCursorValueHandler = func(optionCollection *PagingOptionCollection, modelStruct interface{}) (float64, error) {
	// not cursor mode
	if optionCollection.Option.PagingMode != PagingModeCursor {
		return 0, nil
	}

	mReflectValue := reflect.ValueOf(modelStruct)

	// is pointer struct
	if mReflectValue.Kind() == reflect.Ptr {
		mReflectValue = mReflectValue.Elem()
	}

	// not struct
	if mReflectValue.Kind() != reflect.Struct {
		return 0, fmt.Errorf("ResultSlice value isnot struct")
	}

	// column name
	columnName := StringToCamel(optionCollection.Option.CursorColumn)

	// column exist in struct
	if !mReflectValue.FieldByName(columnName).IsValid() {
		return 0, fmt.Errorf("CursorColumn not exist in ResultSlice struct")
	}

	// column value
	columnValue := mReflectValue.FieldByName(columnName)

	switch columnValue.Kind() {

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return float64(columnValue.Int()), nil

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return float64(columnValue.Uint()), nil

	case reflect.Float32:
		cursorValue, err := strconv.ParseFloat(fmt.Sprint(columnValue.Interface().(float32)), 32)
		if err != nil {
			err = fmt.Errorf("CursorColumn float32 convert to float64 fail : %v", err)
		}
		return cursorValue, err

	case reflect.Float64:
		return columnValue.Float(), nil

	case reflect.String:
		cursorValue, err := strconv.ParseFloat(columnValue.String(), 64)
		if err != nil {
			err = fmt.Errorf("CursorColumn string convert to float64 fail : %v", err)
		}
		return cursorValue, err

	default:
		return 0, fmt.Errorf("CursorColumn value isnot numeric")
	}
}
