package pagination

import (
	"fmt"
	"reflect"
)

// getCursorOptionCollection get cursor mode query option collection
func getCursorOptionCollection(pagingOption *PagingOption, models ...interface{}) (*OptionCollection, error) {
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
var DefaultCursorOptionCollectionHandler = func(pagingOption *PagingOption) *OptionCollection {

	// init cursor query option collection
	pageSize := pagingOption.PageSize

	collection := &OptionCollection{
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
var AnotherCursorOptionCollectionHandler = func(pagingOption *PagingOption) *OptionCollection {

	// init cursor query option collection
	pageSize := pagingOption.PageSize

	collection := &OptionCollection{
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
