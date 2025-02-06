package telegram

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	mysqlconnector "tf-idf/cmd/mysql"
	webcrawler "tf-idf/cmd/web_crawler" // tjgu crawler source.
)

type AllPriceNew struct {
	GoldPrice []struct {
		Gold
	} `json:"gold"`
}

type Gold struct {
	// Date       string
	// Time       string
	Name  string
	Price int
	// ChangeRate float64
	// Unit       string
}

func GetURL() string {
	return fmt.Sprintf("https://brsapi.ir/FreeTsetmcBourseApi/Api_Free_Gold_Currency_v2.json")
}

func loadFile(filePath string) string {
	date, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Println("Cannot read file with error: ", err)
	}

	return string(date)
}

func GetCoinPrice() (string, []int) {
	// For crawl with tjgu source plz uncomment below line.
	webcrawler.GetPrice() // tjgu source
	getAllMessage, err := http.Get(GetURL())
	if err != nil {
		fmt.Println("Cannot get update with: ", err)
	}

	// Make convert resp.http to string:
	body, err := io.ReadAll(getAllMessage.Body)
	if err != nil {
		fmt.Println("Cannot make resp.http to string with: ", err)
	}

	var gold AllPriceNew

	// Make convert string to byte to make json unmarshall:
	getAllMessageByte := []byte(body)
	errJsonUmarshall := json.Unmarshal(getAllMessageByte, &gold)
	if errJsonUmarshall != nil {
		fmt.Println("Cannot make unmarshall json with: ", errJsonUmarshall)
	}

	// Make slice from all gold price.
	priceSlice := []int{}
	fmt.Println(gold.GoldPrice[0].Price)
	priceSlice = append(priceSlice, gold.GoldPrice[6].Price)
	priceSlice = append(priceSlice, gold.GoldPrice[0].Price)
	priceSlice = append(priceSlice, gold.GoldPrice[1].Price)
	priceSlice = append(priceSlice, gold.GoldPrice[2].Price)

	mysqlconnector.UpdatePrice(priceSlice)
	return string(body), priceSlice
}
