package mysqlconnector

// package main

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

var TableName = "word"

type Tabelinfo struct {
	Word      string  `json:"word"`
	Count     int     `json:"count"`
	TF        float64 `json:"tf"`
	IDF       float64 `json:"idf"`
	TFDIDF    float64 `json:"tfidf"`
	TableName string
}
type Asset struct {
	GoldPrice     int `json:"goldprice"`
	NewCoinPrice  int `json:"newcoinprice"`
	OldCoinPrice  int `json:"oldcoinprice"`
	SemiCoinPrice int `json:"semicoinprice"`
}

type SqlConfig struct {
	Password     string
	UserName     string
	MysqlIP      string
	MysqlPort    int
	DatabaseName string
	TableName    string
}

type UsersInfoTable struct {
	Firstnames   []string
	LastName     []string
	Email        []string
	TicketNumber []int
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

func Insert(word string, count int, tf float64, idf float64) {

	tableInfo := Tabelinfo{
		Word:      `json:"word"`,
		TableName: TableName,
	}
	db := MakeConnectionToDB()
	// defer db.Close()

	var insertQuery = fmt.Sprintf("insert into %v(word, count, tf, idf) values ('%v', '%v', '%v', %v)", tableInfo.TableName, word, count, tf, idf)
	insert, err := db.Query(insertQuery)

	if err != nil {
		panic(err.Error())

	}
	// fmt.Println(insert.Columns())
	defer insert.Close()

}

func Update(word string, idf float64) {

	tableInfo := Tabelinfo{
		Word:      `json:"word"`,
		TableName: TableName,
	}
	db := MakeConnectionToDB()
	defer db.Close()
	var uodateQuery = fmt.Sprintf("update %v set idf = %v where word = '%v'", tableInfo.TableName, idf, word)
	update, err := db.Query(uodateQuery)
	if err != nil {
		panic(err.Error())

	}
	defer update.Close()

	var updateTFIDF = fmt.Sprintf("update %v set tfidf = %v", tableInfo.TableName, 0)
	updatetfidf, err := db.Query(updateTFIDF)
	if err != nil {
		panic(err.Error())

	}
	defer updatetfidf.Close()

}

// func SelectQury() UsersInfoTable {
func SelectQury() {
	var word []string
	var tf []float64
	var idf []float64
	tableInfo := Tabelinfo{
		Word:      `json:"word"`,
		TableName: TableName,
	}
	// var usersInfosTable Tabelinfo
	var usersTable Tabelinfo

	db := MakeConnectionToDB()
	selectQuery, err := db.Query("select * from word") // For example: db.Query("select * from users")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	for selectQuery.Next() {
		err = selectQuery.Scan(&usersTable.Word, &usersTable.Count, &usersTable.TF, &usersTable.IDF, &usersTable.TFDIDF)
		if err != nil {
			panic(err.Error())
		}
		// fmt.Println("Value from database: ", usersTable.FirstName, usersTable.Lastname, usersTable.Email, usersTable.TicketNumber)
		word = append(word, usersTable.Word)
		tf = append(tf, usersTable.TF)
		idf = append(idf, usersTable.IDF)

		// Add TF-IDF each word:
		var updateQuery = fmt.Sprintf("update %v set tfidf = %v where word = '%v'", tableInfo.TableName, usersTable.TF*usersTable.IDF, usersTable.Word)
		update, err := db.Query(updateQuery)
		update.Close()
		if err != nil {
			panic(err.Error())

		}
		defer update.Close()

		fmt.Println("Word: ", usersTable.Word, "And TF * IDF is: ", usersTable.TF*usersTable.IDF)
		defer db.Close()
	}
}

// Fuction for add links and labels name.
func InsertLabels(link string, name string, lable1 string, lable2 string, lable3 string) {
	tableName := "links"
	db := MakeConnectionToDB()
	defer db.Close()
	var insertQuery = fmt.Sprintf("insert into %v(link, name, lable, lable1, lable2) values ('%v', '%v', '%v', '%v', '%v')", tableName, link, name, lable1, lable2, lable3)
	insert, err := db.Query(insertQuery)
	if err != nil {
		panic(err.Error())
	}
	defer insert.Close()
}

// Function for show all labels.
func ShowLabels() {
	TableName := "links"
	db := MakeConnectionToDB()
	defer db.Close()
	var selectAllLabels = fmt.Sprintf("select * from %v", TableName)
	selectAllValues, err := db.Query(selectAllLabels)
	if err != nil {
		panic(err.Error())
	}
	defer selectAllValues.Close()
}

// Function For Show Pice Gold.
func SelectPrice() (*sql.Rows, *sql.Rows, *sql.Rows, *sql.Rows) {
	TableName := "price"
	db := MakeConnectionToDB()
	defer db.Close()
	var GoldPrice = fmt.Sprintf("select goldprice from %v", TableName)
	goldPrice, err := db.Query(GoldPrice)
	if err != nil {
		panic(err.Error())
	}
	defer goldPrice.Close()

	var NewCoinPrice = fmt.Sprintf("select newcoinprice from %v", TableName)
	newCoinPrice, err := db.Query(NewCoinPrice)
	if err != nil {
		panic(err.Error())
	}
	defer newCoinPrice.Close()

	var OldCoinPrice = fmt.Sprintf("select newcoinprice from %v", TableName)
	oldCoinPrice, err := db.Query(OldCoinPrice)
	if err != nil {
		panic(err.Error())
	}
	defer oldCoinPrice.Close()

	var SemiCoinPrice = fmt.Sprintf("select newcoinprice from %v", TableName)
	semiCoinPrice, err := db.Query(SemiCoinPrice)
	if err != nil {
		panic(err.Error())
	}
	defer semiCoinPrice.Close()
	return goldPrice, newCoinPrice, oldCoinPrice, semiCoinPrice
}

func SelectPriceGold() (int, int, int, int) {

	// var usersInfosTable Tabelinfo
	var asstStruct Asset

	db := MakeConnectionToDB()
	selectQuery, err := db.Query("select * from price") // For example: db.Query("select * from users")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	for selectQuery.Next() {
		err = selectQuery.Scan(&asstStruct.GoldPrice, &asstStruct.NewCoinPrice, &asstStruct.OldCoinPrice, &asstStruct.SemiCoinPrice)
		if err != nil {
			panic(err.Error())
		}
		fmt.Println(asstStruct.GoldPrice, asstStruct.NewCoinPrice, asstStruct.OldCoinPrice, asstStruct.SemiCoinPrice)
		defer db.Close()
	}
	return asstStruct.GoldPrice, asstStruct.NewCoinPrice, asstStruct.OldCoinPrice, asstStruct.SemiCoinPrice
}
