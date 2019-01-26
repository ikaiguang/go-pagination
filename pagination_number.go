package pagination

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
