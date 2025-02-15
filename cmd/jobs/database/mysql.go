package main

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type SqlConfig struct {
	Password     string
	UserName     string
	MysqlIP      string
	MysqlPort    int
	DatabaseName string
	TableName    string
}

const (
	DBName    = "words"
	TableName = "house_price"
)

type CreateTime struct {
	CREATE_TIME string `json:"CREATE_TIME"`
}

var DBConnection = MakeConnectionToDB()

func main() {
	err := CopyTable(TableName)
	if err == nil {
		TableCleanUP(TableName)
	}
}

func MakeConnectionToDB() *sql.DB {
	SqlConfig := SqlConfig{
		Password:     "test@test",
		UserName:     "root",
		MysqlIP:      "127.0.0.1",
		MysqlPort:    3306,
		DatabaseName: "words",
		TableName:    "word",
	}
	connectioString := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v", SqlConfig.UserName, SqlConfig.Password, SqlConfig.MysqlIP, SqlConfig.MysqlPort, SqlConfig.DatabaseName)
	db, err := sql.Open("mysql", connectioString)

	if err != nil {
		panic(err.Error())
	}
	// defer db.Close()
	return db // Make correct return db.

}

func CheckTableTime() {

	createTime := CreateTime{
		CREATE_TIME: `json:"CREATE_TIME"`,
	}
	var tableCreationTime = fmt.Sprintf("SELECT create_time FROM INFORMATION_SCHEMA.TABLES WHERE table_schema = '%v' and table_name = '%v'", DBName, TableName)
	timeTable, err := DBConnection.Query(tableCreationTime)

	if err != nil {
		fmt.Println("Cannot find creation table with error: ", err)
	}

	for timeTable.Next() {
		err := timeTable.Scan(&createTime.CREATE_TIME)
		if err != nil {
			fmt.Println("Cannot find table creation time with: ", err)
		}
	}
	fmt.Println(createTime.CREATE_TIME)
}

func CopyTable(tableName string) error {
	currentTimeUnix := time.Now().Unix()

	var copyTableQuery = fmt.Sprintf("create table house_price_majidieh_%v as select * from house_price", currentTimeUnix)
	_, err := DBConnection.Query(copyTableQuery)

	if err != nil {
		fmt.Println("Cannot make copy with error: ", err)
	}
	return err
}

func TableCleanUP(tableName string) {
	var copyTableQuery = fmt.Sprintf("delete from house_price")
	_, err := DBConnection.Query(copyTableQuery)

	if err != nil {
		fmt.Println("Cannot DROP table with error: ", err)
	}

}
