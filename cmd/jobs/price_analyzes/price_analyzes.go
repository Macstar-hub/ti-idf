package main

import (
	"database/sql"
	"fmt"
	"math"
	"strconv"

	b64 "encoding/base64"

	ztable "github.com/gregscott94/z-table-golang"
	"github.com/montanaflynn/stats"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"

	mysqlconnector "tf-idf/cmd/mysql"

	// "strconv"
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

type TableInfo struct {
	PerSquar string `json:"per_squar"`
}

var DBConnection = MakeConnectionToDB()

func main() {
	zScoreColumn()
	dbCleaner()

	ZScore(TableName)
	lowerBound, upperBound := IQR(TableName)
	fmt.Println("Lower Bound: ", lowerBound, "Upper Bound: ", upperBound)
	fmt.Println("Fine tune average pice: ", averagePriceZscore(TableName))
	repeatedHousePrice("house_price", "house_price_majidieh_1741164255")
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

func makeAveragePrice(tableName string) (int, []int) {

	id := 0
	sumPrice := 0
	var averagePrice float64
	averagePrice = 0.0

	var priceList []int

	perSquar := TableInfo{
		PerSquar: `json:"per_squar"`,
	}
	var tableCreationTime = fmt.Sprintf("SELECT per_squar from %v where not per_squar = 0", tableName)
	timeTable, err := DBConnection.Query(tableCreationTime)

	if err != nil {
		fmt.Println("Cannot find creation table with error: ", err)
	}

	for timeTable.Next() {
		err := timeTable.Scan(&perSquar.PerSquar)
		if err != nil {
			fmt.Println("Cannot find table creation time with: ", err)
		}
		id++
		priceInt, _ := strconv.Atoi(perSquar.PerSquar)
		priceList = append(priceList, priceInt)
	}

	for i := 0; i < len(priceList); i++ {
		sumPrice = priceList[i] + sumPrice
	}
	averagePrice = float64(sumPrice) / float64(id)

	return int(averagePrice), priceList
}

func ZScore(tableName string) {
	var a int
	var b int
	sigmaPrice := 0
	var deviationPopulation float64
	zTable := ztable.NewZTable(nil)

	_, allPrice := makeAveragePrice(tableName)

	for i := 0; i < len(allPrice); i++ {
		sigmaPrice = allPrice[i] + sigmaPrice
	}
	meanPoplutaion := sigmaPrice / len(allPrice)

	for i := 0; i < len(allPrice); i++ {
		sigmaPrice = (allPrice[i] - meanPoplutaion)
		a = (sigmaPrice * sigmaPrice)
		b = a + b
	}
	sigmaPrice = int(b)
	c := sigmaPrice / len(allPrice)
	deviationPopulation = math.Sqrt(float64(c))

	for i := 0; i < len(allPrice); i++ {
		z := (float64(allPrice[i]) - float64(meanPoplutaion)) / deviationPopulation
		mysqlconnector.UpdateHousePriceZscore(allPrice[i], (zTable.FindPercentage(z) * 100), TableName)
	}
}

func IQR(tableName string) (int, int) {
	var allPriceFloat []float64

	_, allPrice := makeAveragePrice(tableName)
	for i := 0; i < len(allPrice); i++ {
		allPriceFloat = append(allPriceFloat, float64(allPrice[i]))
	}

	Q1, _ := stats.Percentile(allPriceFloat, 25)
	Q3, _ := stats.Percentile(allPriceFloat, 75)

	IQR := Q3 - Q1

	lowerBound := Q1 - 1.5*IQR
	UpperBound := Q3 + 1.5*IQR

	return int(lowerBound), int(UpperBound)
}

func averagePriceZscore(tableName string) int {

	id := 0
	sumPrice := 0
	var averagePrice float64
	averagePrice = 0.0

	var priceList []int

	perSquar := TableInfo{
		PerSquar: `json:"per_squar"`,
	}
	var tableCreationTime = fmt.Sprintf("SELECT per_squar from %v where (not per_squar = 0  and z_score  BETWEEN 35 AND 75)", tableName)
	timeTable, err := DBConnection.Query(tableCreationTime)

	if err != nil {
		fmt.Println("Cannot find creation table with error: ", err)
	}

	for timeTable.Next() {
		err := timeTable.Scan(&perSquar.PerSquar)
		if err != nil {
			fmt.Println("Cannot find table creation time with: ", err)
		}
		id++
		priceInt, _ := strconv.Atoi(perSquar.PerSquar)
		priceList = append(priceList, priceInt)
	}

	for i := 0; i < len(priceList); i++ {
		sumPrice = priceList[i] + sumPrice
	}
	averagePrice = float64(sumPrice) / float64(id)

	return int(averagePrice)
}

func histPlot() {
	plot := plot.New()

	testSlice := []float64{1, 2, 3, 7, 5, 6, 6, 3, 9, 10}

	plot.Title.Text = "Price Histogram"

	hist, err := plotter.NewHist(plotter.Values(testSlice), 20)
	if err != nil {
		fmt.Println("Cannot create histogram with error: ", err)
	}

	plot.Add(hist)
	plot.Save(4*vg.Inch, 4*vg.Inch, "hist.png")

	/*
		For more plot info:
		https://golangdocs.com/plotting-in-golang-histogram-barplot-boxplot
	*/
}

func dbCleaner() {
	var dbCleaner = fmt.Sprintf("delete from %v where per_squar = 0 ;", TableName)
	_, err := DBConnection.Query(dbCleaner)

	if err != nil {
		fmt.Println("Cannot cleanup DB's: ", err)
	}
}

func zScoreColumn() {
	var addZscoreColumn = fmt.Sprintf("alter table %v add z_score float", TableName)
	_, err := DBConnection.Query(addZscoreColumn)

	if err != nil {
		fmt.Println("Cannot make zscore column with error: ", err)
	}
}

func repeatedHousePrice(todayTable string, targetTable string) {
	linkList := mysqlconnector.RepeatedHousePrice(todayTable, targetTable)

	for i := 0; i < len(linkList); i++ {
		link, _ := b64.StdEncoding.DecodeString(linkList[i])
		fmt.Println(string(link))
		fmt.Println()
	}

}
