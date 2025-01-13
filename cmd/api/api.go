package httppost

import (
	"fmt"
	"net/http"

	strconv "strconv"
	mysqlconnector "tf-idf/cmd/mysql"

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

	body.Redirect(http.StatusFound, "/")
}

func CalcAsset(body *gin.Context) {

	assetGeram, _ := strconv.Atoi(body.PostForm("assetGeram"))
	newCoin, _ := strconv.Atoi(body.PostForm("newCoin"))
	oldCoin, _ := strconv.Atoi(body.PostForm("oldCoin"))
	semiCoin, _ := strconv.Atoi(body.PostForm("semiCoin"))

	// Get price from mysql
	// goldPrice, newCoinPrice, oldCoinPrice, semiCoinPrice := mysqlconnector.SelectPrice()
	// GoldPrice, _ := strconv.Atoi(goldPrice)
	// NewCoinPrice, _ := strconv.Atoi(newCoinPrice)
	// OldCoinPrice, _ := strconv.Atoi(oldCoinPrice)
	// SemiCoinPrice, _ := strconv.Atoi(semiCoinPrice)

	GoldPrice := 5322100
	NewCoinPrice := 56950000
	OldCoinPrice := 54650000
	SemiCoinPrice := 31500000

	fmt.Println(assetGeram, newCoin, oldCoin, semiCoin)
	totalAsset := (assetGeram * GoldPrice) + (newCoin * NewCoinPrice) + (oldCoin * OldCoinPrice) + (semiCoin * SemiCoinPrice)
	body.HTML(http.StatusOK, "assetCalc.html", gin.H{"totalAsset": totalAsset})

	fmt.Println("Total asset is: ", totalAsset)
}
