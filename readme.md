# go-pagination

The fantastic pagination library for Golang, aims to be developer friendly.

[![Go Report Card](https://goreportcard.com/badge/github.com/ikaiguang/go-pagination)](https://goreportcard.com/report/github.com/ikaiguang/go-pagination)
[![GoDoc](https://godoc.org/github.com/ikaiguang/go-pagination?status.svg)](https://godoc.org/github.com/ikaiguang/go-pagination)

## Install

install && get source code

`go get -v github.com/ikaiguang/go-pagination`

> go_1.11+ : `GO111MODULE=off go get -v github.com/ikaiguang/go-pagination`

## Test

before run test : you must rewrite database connection and generate test data

```

go get -v github.com/ikaiguang/go-pagination-example
# GO111MODULE=off go get -v github.com/ikaiguang/go-pagination-example

go test -v $GOPATH/src/github.com/ikaiguang/go-pagination-example

```

## Overview

see [pagination.md](pagination.md)

## Getting Started

```shell

# get source code

cd $GOPATH/src/github.com/ikaiguang/go-pagination

git checkout test

```

1. rewrite database connection and generate test data
    - see [go-pagination-example/test_data.go](https://github.com/ikaiguang/go-pagination-example/blob/master/test_data.go)
2. declare a table model
    - see [go-pagination-example/model.go](https://github.com/ikaiguang/go-pagination-example/blob/master/model.go)
3. declare a controller
    - see [go-pagination-example/controller.go](https://github.com/ikaiguang/go-pagination-example/blob/master/controller.go)
4. run test
    - see [go-pagination-example/example_test.go](https://github.com/ikaiguang/go-pagination-example/blob/master/example_test.go)

### create a table and insert test data

> see [go-pagination-example/test_data.go](https://github.com/ikaiguang/go-pagination-example/blob/master/test_data.go)

```sql

# create table

CREATE TABLE pagination_users (
  id   INT          NOT NULL AUTO_INCREMENT
  COMMENT 'user table id',
  name VARCHAR(255) NOT NULL DEFAULT ''
  COMMENT 'user name',
  age  TINYINT(3)   NOT NULL DEFAULT '0'
  COMMENT 'user age',
  PRIMARY KEY (id)
)
  ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COMMENT = 'user table';

# insert test data

INSERT INTO pagination_users (id, name, age)
VALUES
  (1, 'user_name_1', 1),
  (2, 'user_name_2', 2),
  (3, 'user_name_3', 3),
  (4, 'user_name_4', 4),
  (5, 'user_name_5', 5),
  (6, 'user_name_6', 6),
  (7, 'user_name_7', 7),
  (8, 'user_name_8', 8),
  (9, 'user_name_9', 9);

```

### declare a table model

> see [go-pagination-example/model.go](https://github.com/ikaiguang/go-pagination-example/blob/master/model.go)

```go

package example

import (
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/jinzhu/gorm"
	"fmt"
	page "github.com/ikaiguang/go-pagination"
)

// new db connection
func newDbConnection() (*gorm.DB, error) {

	dbDiver := "mysql"
	dbDsn := "root:Mysql.123456@tcp(127.0.0.1:3306)/test?"
	dbDsn += "charset=utf8mb4&parseTime=True&loc=Local"

	db, err := gorm.Open(dbDiver, dbDsn)

	if err != nil {
		err = fmt.Errorf("gorm.Open dsn error : %v", err)
		return nil, err
	}
	//defer db.Close()

	db.LogMode(true)

	return db, nil
}

// where condition
type WhereCondition struct {
	page.PagingWhere
}

// where conditions
func WhereConditions(dbConn *gorm.DB, whereConditions []*WhereCondition) *gorm.DB {

	// where
	for _, where := range whereConditions {
		// db.Where("id = ?", id)
		whereStr := fmt.Sprintf("%s %s %s", where.Column, where.Symbol, where.Placeholder)
		dbConn = dbConn.Where(whereStr, where.Data)
	}
	return dbConn
}

// pagination
func Pagination(dbConn *gorm.DB, pagingOptionCollection *page.OptionCollection) *gorm.DB {
	// limit offset
	dbConn = dbConn.Limit(pagingOptionCollection.Limit).Offset(pagingOptionCollection.Offset)

	// where
	for _, where := range pagingOptionCollection.Where {
		// db.Where("id = ?", id)
		whereStr := fmt.Sprintf("%s %s %s", where.Column, where.Symbol, where.Placeholder)
		dbConn = dbConn.Where(whereStr, where.Data)
	}

	// order
	for _, order := range pagingOptionCollection.Order {
		dbConn = dbConn.Order(fmt.Sprintf("%s %s", order.Column, order.Direction))
	}
	return dbConn
}

// user model
type UserModel struct {
	Id   int64  `gorm:"PRIMARY_KEY;COLUMN:id"` // id
	Name string `gorm:"COLUMN:name"`           // name
	Age  int64  `gorm:"COLUMN:age"`            // age
}

// user table name
func (m *UserModel) TableName() string {
	return "pagination_users"
}

// new model
func (m *UserModel) NewModel() *UserModel {
	return new(UserModel)
}

// user model list
func (m *UserModel) List(whereConditions []*WhereCondition, pagingOptionCollection *page.OptionCollection) (*[]UserModel, int64, error) {
	var count int64
	var list []UserModel

	// db conn
	db, err := newDbConnection()
	if err != nil {
		err = fmt.Errorf("UserModel.List newDbConnection error : %v", err)
		return &list, count, err
	}
	defer db.Close()

	// query where
	userDb := WhereConditions(db, whereConditions)
	defer userDb.Close()

	if err := userDb.Table(m.TableName()).Count(&count).Error; err != nil {
		err = fmt.Errorf("UserModel.List Count error : %v", err)
		return &list, count, err
	} else if count == 0 {
		return &list, count, err // empty
	}

	// pagination
	if err := Pagination(userDb, pagingOptionCollection).Find(&list).Error; err != nil {
		err = fmt.Errorf("UserModel.List Find error : %v", err)
		return &list, count, err
	}
	return &list, count, err
}

```

### declare a controller

> see [go-pagination-example/controller.go](https://github.com/ikaiguang/go-pagination-example/blob/master/controller.go)

```go

package example

import (
	"fmt"
	page "github.com/ikaiguang/go-pagination"
)

type UserController struct {
	Model *UserModel
}

// get user list
func (c *UserController) List(option *page.PagingOption) (*[]UserModel, *page.PagingResult, error) {

	// get paging query option collection
	pagingOptionCollection, err := page.GetOptionCollection(option, c.Model.NewModel())
	if err != nil {
		err := fmt.Errorf("UserController.List GetOptionCollection error : %v", err)
		return nil, nil, err
	}

	// list
	list, count, err := c.Model.List([]*WhereCondition{}, pagingOptionCollection)
	if err != nil {
		err := fmt.Errorf("UserController.List model.List error : %v", err)
		return nil, nil, err
	}

	// set paging result
	pagingResultCollection := &page.PagingResultCollection{
		TotalRecords: count,
		ListPointer:  list,
	}
	pagingResult, err := page.SetPagingResult(pagingOptionCollection, pagingResultCollection)
	if err != nil {
		err := fmt.Errorf("UserController.List SetPagingResult error : %v", err)
		return nil, nil, err
	}
	return list, pagingResult, nil
}

```

## run example

> before run test : you must rewrite database connection and generate test data

`go test -v $GOPATH/src/github.com/ikaiguang/go-pagination-example/`

> see [go-pagination-example/example_test.go](https://github.com/ikaiguang/go-pagination-example/blob/master/example_test.go)

```go

// paging mode : page number
func testPageNumberMode(t *testing.T) {

	var controller UserController           // controller
	var list *[]UserModel                   // list
	var pagingResult *page.PagingResult     // page result
	var err error                           // error
	var pageNumberOption *page.PagingOption // option

	// ===== paging mode : page number ===== //
	// ===== goto 3rd page ===== //
	// page number option
	pageNumberOption = page.DefaultPagingOption()
	pageNumberOption.PageSize = 2                                        // page size : 2
	orderBy := &page.PagingOrder{Column: "age", Direction: "desc"}       // order by age desc
	pageNumberOption.OrderBy = append(pageNumberOption.OrderBy, orderBy) // order by age desc
	pageNumberOption.GotoPageNumber = 3                                  // goto 3rd page

	list, pagingResult, err = controller.List(pageNumberOption)
	if err != nil {
		t.Errorf("testing : controller.List error : %v", err)
		return
	} else {
		format := "\n page number : order age desc && goto 3rd page \n"
		format += "\n paging result : %v \n"
		format += "\n list : %v \n"

		t.Logf(format, pagingResult, list)
	}
}

```

## License

Released under the [MIT License](License)