package httppost

import (
	"fmt"
	"net/http"

	strconv "strconv"
	mysqlconnector "tf-idf/cmd/mysql"
	"tf-idf/cmd/telegram"

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

	assetGeram, _ := strconv.Atoi(body.PostForm("assetGeram"))
	newCoin, _ := strconv.Atoi(body.PostForm("newCoin"))
	oldCoin, _ := strconv.Atoi(body.PostForm("oldCoin"))
	semiCoin, _ := strconv.Atoi(body.PostForm("semiCoin"))

	// Get price from mysql
	telegram.GetCoinPrice()
	goldPrice, newCoinPrice, oldCoinPrice, semiCoinPrice := mysqlconnector.SelectPriceGold()

	totalAsset := (assetGeram * goldPrice) + (newCoin * newCoinPrice) + (oldCoin * oldCoinPrice) + (semiCoin * semiCoinPrice)

	// Render all Gold asset
	body.HTML(http.StatusOK, "assetCalc.html", gin.H{"totalAsset": totalAsset})
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
