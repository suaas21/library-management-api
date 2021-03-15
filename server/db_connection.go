package server

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

func (svr Server) InitializeDB() {
	fmt.Println("Connecting database........")
	var err error
	tables = append(tables, new(model.User), new(model.Book), new(model.BookLoanHistory))
	eng, err = svr.GetPostgresClient()
	if err != nil {
		fmt.Println("Unable to connect ORM engine, reason: ", err)
	}

	if err := eng.Ping(); err != nil {
		fmt.Println("Unable to ping store reason: ", err.Error())
	}

	eng.SetTableMapper(core.SameMapper{})
	eng.SetColumnMapper(core.SnakeMapper{})
	if err := eng.Sync2(tables...); err != nil {
		fmt.Println("Unable to sync struct to store table: reason: ", err.Error())
	}
}
func (svr Server) GetPostgresClient() (*xorm.Engine, error) {
	cnnstr := fmt.Sprintf("user=%s password=%s host=127.0.0.1 port=%v dbname=%s sslmode=disable", svr.DBUser, svr.DBPassword, svr.DBPort, svr.DBName)
	return xorm.NewEngine("postgres", cnnstr)
}