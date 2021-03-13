package db

import (
	"fmt"
	"github.com/suaas21/library-management-api/model"
	"xorm.io/core"

	"github.com/go-xorm/xorm"
	_ "github.com/lib/pq"
)

var (
	eng    *xorm.Engine
	tables []interface{}
)

func InitializeDB(dbUser, dbPassword, dbName, dbPort string) {
	fmt.Println("Connecting database........")
	var err error
	tables = append(tables, new(model.UserDB), new(model.Book), new(model.BookHistory))
	eng, err = GetPostgresClient(dbUser, dbPassword, dbName, dbPort)
	if err != nil {
		fmt.Println("Unable to connect ORM engine, reason: ", err)
	}

	if err := eng.Ping(); err != nil {
		fmt.Println("Unable to ping db reason: ", err.Error())
	}

	eng.SetTableMapper(core.SameMapper{})
	eng.SetColumnMapper(core.SnakeMapper{})
	if err := eng.Sync2(tables...); err != nil {
		fmt.Println("Unable to sync struct to db table: reason: ", err.Error())
	}
}
func GetPostgresClient(dbUser, dbPassword, dbName, dbPort string) (*xorm.Engine, error) {
	cnnstr := fmt.Sprintf("user=%s password=%s host=127.0.0.1 port=%v dbname=%s sslmode=disable", dbUser, dbPassword, dbPort, dbName)
	return xorm.NewEngine("postgres", cnnstr)
}
