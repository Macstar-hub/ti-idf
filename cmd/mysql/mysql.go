package mysqlconnector

// package main

import (
	"database/sql"
	b64 "encoding/base64"
	"fmt"
	"strconv"

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
type HouseInfo struct {
	ID        string `json:"id"`
	Links     string `json:"links"`
	LinkList  string `json: "id"`
	TableName string
}
type Asset struct {
	GoldPrice     int `json:"goldprice"`
	NewCoinPrice  int `json:"newcoinprice"`
	OldCoinPrice  int `json:"oldcoinprice"`
	SemiCoinPrice int `json:"semicoinprice"`
}

type ShowLinksStruct struct {
	Link   string `json:"link"`
	Name   string `json:"name"`
	Label  string `json:"label"`
	Label1 string `json:"label1"`
	Label2 string `json:"label2"`
}

type ShowLinksStructList struct {
	Link   []string
	Name   []string
	Label  []string
	Label1 []string
	Label2 []string
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

func UpdatePrice(goldPrice []int) {
	tableInfo := Tabelinfo{
		TableName: "price",
	}
	db := MakeConnectionToDB()
	defer db.Close()
	var updateQuery = fmt.Sprintf("update %v set goldprice = %v, newcoinprice = %v, oldcoinprice = %v, semicoinprice = %v ", tableInfo.TableName, goldPrice[0], goldPrice[1], goldPrice[2], goldPrice[3])
	update, err := db.Query(updateQuery)
	if err != nil {
		panic(err.Error())

	}
	defer update.Close()

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

func InsertHousePrice(id int64, links string, per_squar int64, total_squar int64) {
	tableName := "house_price"
	db := MakeConnectionToDB()
	defer db.Close()
	var insertQuery = fmt.Sprintf("insert into %v(id, links, per_squar, total_squar) values ('%v', '%v', '%v', '%v')", tableName, id, links, per_squar, total_squar)
	insert, err := db.Query(insertQuery)
	if err != nil {
		panic(err.Error())
	}
	defer insert.Close()
}

func SelectHousePrice() ([]int, []string, int) {

	var linkList []int
	var id []int
	var links []string
	var i = 0

	tableInfo := HouseInfo{
		ID:       `json:id`,
		Links:    `json:"links"`,
		LinkList: `json:"id"`,
	}

	db := MakeConnectionToDB()
	selectQuery, err := db.Query("select id, links from house_price where per_squar = 0")
	if err != nil {
		panic(err.Error())
	}

	defer db.Close()

	for selectQuery.Next() {
		err = selectQuery.Scan(&tableInfo.ID, &tableInfo.Links)
		if err != nil {
			panic(err.Error())
		}

		intID, _ := strconv.Atoi(tableInfo.ID)
		id = append(id, intID)
		decodeLink, _ := b64.StdEncoding.DecodeString(tableInfo.Links)
		links = append(links, string(decodeLink))
		i++
		if i >= 2 {
			break
		}
		defer db.Close()
	}

	linkListSelectQuey, err := db.Query("select id from house_price")
	if err != nil {
		panic(err.Error())
	}

	for linkListSelectQuey.Next() {
		err = linkListSelectQuey.Scan(&tableInfo.LinkList)
		if err != nil {
			panic(err.Error())
		}

		allLinkCount, _ := strconv.Atoi(tableInfo.LinkList)
		linkList = append(id, allLinkCount)
		defer db.Close()
	}

	return id, links, len(linkList)
}

func RepeatedHousePrice(todayTable string, targetTable string) []string {
	db := MakeConnectionToDB()
	var linkList []string

	innerJoin := fmt.Sprintf("SELECT  %v.links  FROM %v INNER JOIN %v  ON %v.links= %v.links   where (not house_price.per_squar = 0  and house_price.z_score  BETWEEN 35 AND 75);", todayTable, todayTable, targetTable, targetTable, todayTable)
	innerJoinQuery, err := db.Query(innerJoin)
	if err != nil {
		fmt.Println("Cannot make inner join with error: ", err)
	}

	for innerJoinQuery.Next() {

		var links string
		// var persquar_source int
		// var persqaur_target int

		innerJoinQuery.Scan(&links)
		innerJoinQuery.Scan(fmt.Sprintf("%v.per_squar", todayTable))
		innerJoinQuery.Scan(fmt.Sprintf("%v.per_squar", targetTable))

		linkList = append(linkList, links)
	}

	defer db.Close()
	return linkList
}

func TelegramHousePriceSend(todayTable string) []string {
	db := MakeConnectionToDB()
	var linkList []string

	innerJoin := fmt.Sprintf("SELECT  links  FROM %v  where (not per_squar = 0  and z_score  BETWEEN 35 AND 75);", todayTable)
	innerJoinQuery, err := db.Query(innerJoin)
	if err != nil {
		fmt.Println("Cannot make inner join with error: ", err)
	}

	for innerJoinQuery.Next() {

		var links string
		// var persquar_source int
		// var persqaur_target int

		innerJoinQuery.Scan(&links)
		innerJoinQuery.Scan(fmt.Sprintf("%v.per_squar", todayTable))
		linkList = append(linkList, links)
	}

	defer db.Close()
	return linkList
}

func UpdateHousePrice(id int, per_squar int, total_squar int) {
	tableName := "house_price"
	db := MakeConnectionToDB()
	defer db.Close()
	var updateQuery = fmt.Sprintf("update %v set per_squar = %v, total_squar = %v where id = %v", tableName, per_squar, total_squar, id)
	update, err := db.Query(updateQuery)
	if err != nil {
		panic(err.Error())
	}
	defer update.Close()
}

func UpdateHousePriceZscore(per_squar int, z_score float64, TableName string) {
	tableName := TableName
	db := MakeConnectionToDB()
	defer db.Close()
	var updateQuery = fmt.Sprintf("update %v set z_score = %v where per_squar = %v", tableName, z_score, per_squar)
	update, err := db.Query(updateQuery)
	if err != nil {
		panic(err.Error())
	}
	defer update.Close()
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
	}
	defer db.Close()
	return asstStruct.GoldPrice, asstStruct.NewCoinPrice, asstStruct.OldCoinPrice, asstStruct.SemiCoinPrice
}

// Render all links in table.
func ShowLinks() ShowLinksStructList {

	var Link []string
	var Name []string
	var Lable []string
	var Label1 []string
	var Lable2 []string
	// var usersInfosTable Tabelinfo
	var showLinksStruct ShowLinksStruct

	db := MakeConnectionToDB()
	selectQuery, err := db.Query("select * from links")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	for selectQuery.Next() {
		err = selectQuery.Scan(&showLinksStruct.Link, &showLinksStruct.Name, &showLinksStruct.Label, &showLinksStruct.Label1, &showLinksStruct.Label2)

		Link = append(Link, showLinksStruct.Link)
		Name = append(Name, showLinksStruct.Name)
		Lable = append(Lable, showLinksStruct.Label)
		Label1 = append(Label1, showLinksStruct.Label1)
		Lable2 = append(Lable2, showLinksStruct.Label2)

		if err != nil {
			panic(err.Error())
		}
	}
	fmt.Println("++++++++++++++++", showLinksStruct.Link, showLinksStruct.Name, showLinksStruct.Label, showLinksStruct.Label1, showLinksStruct.Label2)

	// -----------------------------> Just searchQuery debug.
	// SearchRecord("test") // Just for debug.

	showLinksStructList := ShowLinksStructList{
		Link:   Link,
		Name:   Name,
		Label:  Lable,
		Label1: Label1,
		Label2: Lable2,
	}

	defer db.Close()
	return showLinksStructList
}

// Make search function:
func SearchRecord(searchWord string) (ShowLinksStructList, int) {

	var showLinksStruct ShowLinksStruct
	var Link []string
	var Name []string
	var Label []string
	var Label1 []string
	var Label2 []string

	/* Note: we can use switch, case in this funciton in mode of search filed
	by labels, name OR even links regex.
	*/

	query := fmt.Sprintf("select * from links where name = '%s'", searchWord)

	db := MakeConnectionToDB()
	searchQuery, err := db.Query(query)
	if err != nil {
		panic(err.Error())
	}
	for searchQuery.Next() {
		err = searchQuery.Scan(&showLinksStruct.Link, &showLinksStruct.Name, &showLinksStruct.Label, &showLinksStruct.Label1, &showLinksStruct.Label2)
		Link = append(Link, showLinksStruct.Link)
		Name = append(Name, showLinksStruct.Name)
		Label = append(Label, showLinksStruct.Label)
		Label1 = append(Label1, showLinksStruct.Label1)
		Label2 = append(Label2, showLinksStruct.Label2)
		if err != nil {
			panic(err.Error())
		}
	}

	showLinksStructList := ShowLinksStructList{
		Link:   Link,
		Name:   Name,
		Label:  Label,
		Label1: Label1,
		Label2: Label2,
	}
	element := len(showLinksStructList.Name)
	fmt.Println("Debug from searhQuery function: ", showLinksStructList)
	defer db.Close()
	return showLinksStructList, element
}

// Make edit function:
// func EditRecord() {
// 	query := fmt.Sprintf("select * from links where name = '%s'", searchWord)
// }
