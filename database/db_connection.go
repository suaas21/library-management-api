package database

import (
	"fmt"
	"xorm.io/core"

	"github.com/go-xorm/xorm"
	_ "github.com/lib/pq"
)

var (
	eng    *xorm.Engine
	tables []interface{}
)

type Config struct {
	DBPort     string
	DBPassword string
	DBName     string
	BDUser     string
}

var Cfg Config


func InitializeDB() {
	fmt.Println("Connecting database........")
	var err error
	tables = append(tables, new(User), new(Book), new(BookLoanHistory), new(BookLoanRequest))
	eng, err = GetPostgresClient()
	if err != nil {
		fmt.Println("Unable to connect ORM engine, reason: ", err)
	}

	if err := eng.Ping(); err != nil {
		fmt.Println("Unable to ping store reason: ", err.Error())
	}

	eng.SetTableMapper(core.SameMapper{})
	eng.SetColumnMapper(core.SnakeMapper{})
	// sync data table to the database
	if err := eng.Sync2(tables...); err != nil {
		fmt.Println("Unable to sync struct to store table: reason: ", err.Error())
	}
}
func GetPostgresClient() (*xorm.Engine, error) {
	// generate database connection string
	// then connect to the database
	cnnstr := fmt.Sprintf("user=%s password=%s host=127.0.0.1 port=%v dbname=%s sslmode=disable", Cfg.BDUser, Cfg.DBPassword, Cfg.DBPort, Cfg.DBName)
	return xorm.NewEngine("postgres", cnnstr)
}