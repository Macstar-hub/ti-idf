package httppost

import (
	"fmt"
	// "log"
	"net/http"
	strconv "strconv"
	mysqlconnector "tf-idf/cmd/mysql"
	redisclient "tf-idf/cmd/redisClient"
	"tf-idf/cmd/telegram"

	// "tf-idf/cmd/telegram"
	"time"

	"github.com/gin-gonic/gin"
)

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

	// Get price from mysql
	telegram.GetCoinPrice()
	mysqlconnector.UpdatePrice()

	// Just for debug:
	fmt.Println("Just Before Received channel:")
	goldPrice := redisclient.RedisGetOPS("Gold18")          //receivePrice.Gold18
	newCoinPrice := redisclient.RedisGetOPS("SekkeTamam")   //receivePrice.SekkeTamam
	oldCoinPrice := redisclient.RedisGetOPS("SekketGhadim") //receivePrice.SekketGhadim
	semiCoinPrice := redisclient.RedisGetOPS("SekkehNim")   //receivePrice.SekkehNim
	stockGoldPrice := redisclient.RedisGetOPS("GoldDast2")  //receivePrice.GoldDast2
	quarteCoinPrice := redisclient.RedisGetOPS("RobeSekke") //receivePrice.RobeSekke
	usdDollar := redisclient.RedisGetOPS("Dollar")          //receivePrice.Dollar
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

	// Make a channel to receive all data
	// select {
	// case receivePrice := <-telegram.FrontPriceChannel:
	// 	fmt.Println("From channel: ", receivePrice)
	// 	goldPrice := redisclient.RedisGetOPS("Gold18")          //receivePrice.Gold18
	// 	newCoinPrice := redisclient.RedisGetOPS("SekkeTamam")   //receivePrice.SekkeTamam
	// 	oldCoinPrice := redisclient.RedisGetOPS("SekketGhadim") //receivePrice.SekketGhadim
	// 	semiCoinPrice := redisclient.RedisGetOPS("SekkehNim")   //receivePrice.SekkehNim
	// 	stockGoldPrice := redisclient.RedisGetOPS("GoldDast2")  //receivePrice.GoldDast2
	// 	quarteCoinPrice := redisclient.RedisGetOPS("RobeSekke") //receivePrice.RobeSekke
	// 	usdDollar := redisclient.RedisGetOPS("Dollar")          //receivePrice.Dollar
	// 	totalAsset := (assetGeram * goldPrice) + (newCoin * newCoinPrice) + (oldCoin * oldCoinPrice) + (semiCoin * semiCoinPrice)

	// 	// Render all Gold asset
	// 	body.HTML(http.StatusOK, "assetCalc.html", gin.H{
	// 		"totalAsset":      totalAsset,
	// 		"goldPrice":       goldPrice,
	// 		"newCoin":         newCoinPrice,
	// 		"oldCoinPrice":    oldCoinPrice,
	// 		"semiCoinPrice":   semiCoinPrice,
	// 		"quarteCoinPrice": quarteCoinPrice,
	// 		"stockGoldPrice":  stockGoldPrice,
	// 		"usdDollar":       usdDollar,
	// 	})

	// case <-time.After(1 * time.Millisecond):
	// 	log.Println("Timeout meet.")
	// 	close(telegram.FrontPriceChannel)
	// }

	fmt.Println("Total Price Latency Is: ", time.Since(startTime()))
}

// Render all links in table.
func ShowLinks(body *gin.Context) {

	// // Get price from mysql
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
