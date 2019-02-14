# paginate

```bash

protoc -I. -I$GOPATH/src --go_out=. ./*.proto

```

## base mode : jump to page number

```

// page number mode
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
// page number mode

```

## special model : cursor mode

```

// cursor mode
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
// cursor mode

```

### another cursor mode

```

// cursor mode
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
// cursor mode

```
