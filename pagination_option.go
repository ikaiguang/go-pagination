package pagination

import (
	"fmt"
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
