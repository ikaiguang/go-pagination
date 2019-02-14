package pagination

import (
	"fmt"
	"reflect"
	"strconv"
)

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
	//if optionCollection.Option.PagingMode != PagingModeCursor {
	//	return 0, nil
	//}

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
