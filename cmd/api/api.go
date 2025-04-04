package httppost

import (
	"fmt"
	"reflect"
	"sync"

	// "log"
	"net/http"
	strconv "strconv"
	mysqlconnector "tf-idf/cmd/mysql"
	redisclient "tf-idf/cmd/redisClient"

	// "tf-idf/cmd/telegram"
	"time"

	"github.com/gin-gonic/gin"
)

// var Symbol = []string{
// 	"Gold18",
// 	"SekkeTamam",
// 	"SekketGhadim",
// 	"SekkehNim",
// 	"GoldDast2",
// 	"RobeSekke",
// 	"Dollar",
// }

type Symbol struct {
	Gold18       int
	SekkeTamam   int
	SekketGhadim int
	SekkehNim    int
	GoldDast2    int
	RobeSekke    int
	Dollar       int
}

var goldPrice = 0
var newCoinPrice = 0
var oldCoinPrice = 0
var semiCoinPrice = 0
var stockGoldPrice = 0
var quarteCoinPrice = 0
var usdDollar = 0

func PostLabels(body *gin.Context) {

	// Input valuse sections:
	links := body.PostForm("link")
	name := body.PostForm("name")
	label1 := body.PostForm("label1")
	label2 := body.PostForm("label2")
	label3 := body.PostForm("label3")

	// Insert data to mysql
	mysqlconnector.InsertLabels(links, name, label1, label2, label3)

	body.Redirect(http.StatusFound, "/api/v1/linkslist")
}

func CalcAsset(body *gin.Context) {

	startTime := time.Now

	assetGeram, _ := strconv.Atoi(body.PostForm("assetGeram"))
	newCoin, _ := strconv.Atoi(body.PostForm("newCoin"))
	oldCoin, _ := strconv.Atoi(body.PostForm("oldCoin"))
	semiCoin, _ := strconv.Atoi(body.PostForm("semiCoin"))

	newSymbol := new(Symbol)

	redisPriceChannel := make(chan Symbol, 1024)
	wg := &sync.WaitGroup{}

	key := reflect.ValueOf(*newSymbol)
	for i := 0; i < key.NumField(); i++ {
		wg.Add(1)
		go newSymbol.getPriceFromRedis(wg, redisPriceChannel, string(key.Type().Field(i).Name), int(key.Field(i).Int()))
		fmt.Println(string(key.Type().Field(i).Name), int(key.Field(i).Int()))
	}
	wg.Wait()
	close(redisPriceChannel)

	for price := range redisPriceChannel {
		goldPrice = price.Gold18
		newCoinPrice = price.SekkeTamam
		oldCoinPrice = price.SekketGhadim
		semiCoinPrice = price.SekkehNim
		stockGoldPrice = price.SekketGhadim
		quarteCoinPrice = price.RobeSekke
		usdDollar = price.Dollar
	}

	totalAsset := (assetGeram * goldPrice) + (newCoin * newCoinPrice) + (oldCoin * oldCoinPrice) + (semiCoin * semiCoinPrice)

	// Render all Gold asset
	body.HTML(http.StatusOK, "assetCalc.html", gin.H{
		"totalAsset":      totalAsset,
		"goldPrice":       goldPrice,
		"newCoin":         newCoinPrice,
		"oldCoinPrice":    oldCoinPrice,
		"semiCoinPrice":   semiCoinPrice,
		"quarteCoinPrice": quarteCoinPrice,
		"stockGoldPrice":  stockGoldPrice,
		"usdDollar":       usdDollar,
	})

	fmt.Println("Total price calculation latency: ", time.Since(startTime()))
}

func (s *Symbol) getPriceFromRedis(wg *sync.WaitGroup, redisPriceChannel chan Symbol, symbol string, structValue int) {
	startTime := time.Now
	defer wg.Done()
	structValue = redisclient.RedisGetOPS(symbol)

	// Make set value to struct:
	reflect.ValueOf(s).Elem().FieldByName(symbol).SetInt(int64(structValue))
	redisPriceChannel <- *s

	fmt.Println("Latency in make redis get ops to update price: ", time.Since(startTime()))
}

// Render all links in table.
func ShowLinks(body *gin.Context) {

	// Get price from mysql
	showLinksStruct := mysqlconnector.ShowLinks()
	allRecords := len(showLinksStruct.Link)
	var links []gin.H
	for i := 0; i < allRecords; i++ {

		links = append(links, gin.H{
			"Links":  showLinksStruct.Link[i],
			"Name":   showLinksStruct.Name[i],
			"Label":  showLinksStruct.Label[i],
			"Label1": showLinksStruct.Label1[i],
			"Label2": showLinksStruct.Label2[i],
		})
		fmt.Println("===============", showLinksStruct)
	}
	body.HTML(http.StatusOK, "allLinks.html", gin.H{
		"Links": links,
	})
}

// Search function:
func Search(body *gin.Context) {

	// Search based on label and name strings.
	searchWord := body.PostForm("search")
	showLinksStruct, allRecords := mysqlconnector.SearchRecord(searchWord)
	fmt.Println("______________________________", allRecords)
	var links []gin.H
	for i := 0; i < allRecords; i++ {

		links = append(links, gin.H{
			"Links":  showLinksStruct.Link[i],
			"Name":   showLinksStruct.Name[i],
			"Label":  showLinksStruct.Label[i],
			"Label1": showLinksStruct.Label1[i],
			"Label2": showLinksStruct.Label2[i],
		})
		fmt.Println("+++++++++++++++", showLinksStruct)
	}
	body.HTML(http.StatusOK, "allLinks.html", gin.H{
		"Links": links,
	})
}
