package main

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	// "sync"
	"time"
	// "time"
	// "github.com/mavihq/persian"
)

type Price struct {
	Dollar       int
	SekkeTamam   int
	SekketGhadim int
	SekkehNim    int
	RobeSekke    int
	Gold18       int
	GoldDast2    int
}

func main() {

	startTime := time.Now()

	// responseChannel := make(chan int, 1024)
	// wg := *&sync.WaitGroup{}

	var price Price

	usdPrice, _, _ := httpGet("https://www.tgju.org/profile/price_dollar_rl", "priceGold")
	price.Dollar = usdPrice

	sekkeTamamPrice, _, _ := httpGet("https://www.tgju.org/profile/sekee", "priceGold")
	price.SekkeTamam = sekkeTamamPrice

	sekkeGhadimPrice, _, _ := httpGet("https://www.tgju.org/profile/sekeb", "priceGold")
	price.SekketGhadim = sekkeGhadimPrice

	SekkehNimPrice, _, _ := httpGet("https://www.tgju.org/profile/nim", "priceGold")
	price.SekkehNim = SekkehNimPrice

	SekkehRobePrice, _, _ := httpGet("https://www.tgju.org/profile/rob", "priceGold")
	price.RobeSekke = SekkehRobePrice

	Gold18, _, _ := httpGet("https://www.tgju.org/profile/geram18", "priceGold")
	price.Gold18 = Gold18

	GoldDast2, _, _ := httpGet("https://www.tgju.org/profile/gold_mini_size", "priceGold")
	price.GoldDast2 = GoldDast2

	fmt.Println(price)

	log.Println("Total latency: ", time.Since(startTime))

}

func httpGet(url string, priceType string) (int, []int, []string) {
	var price int
	var maskanPrice []int
	var maskanURLS []string

	netClient := customHttpClient()

	responseByte, err := netClient.Get(url)

	httpErrorHandeler(err)

	responeBody, err := ioutil.ReadAll(responseByte.Body)
	byteReadErrorHandelete(err)

	responseString := string(responeBody)

	if priceType == "priceGold" {
		_, price = findSekkeTamam(responseString)
	}

	responseByte.Body.Close()

	return price, maskanPrice, maskanURLS
}

func httpErrorHandeler(err error) error {
	if err != nil {
		fmt.Println("Cannot http call with error: ", err)
	}
	return err
}

func byteReadErrorHandelete(err error) error {
	if err != nil {
		fmt.Println("Cannot read as byte: ", err)
	}
	return err
}

func customHttpClient() http.Client {
	config := &tls.Config{
		InsecureSkipVerify: true,
	}

	transport := &http.Transport{
		TLSClientConfig: config,
	}

	netClient := &http.Client{
		Transport: transport,
	}
	return *netClient
}

func findSekkeTamam(html string) (string, int) {
	regex, _ := regexp.Compile("info.last_trade.PDrCotVal.*")
	price := regex.FindString(html)
	priceInt := priceCleaner(price)
	return price, priceInt
}

func priceCleaner(priceString string) int {
	regexInt, _ := regexp.Compile("[0-9].*")

	// Make Clean "info.last_trade.PDrCotVal">195,000,000</span>"
	priceByte := regexInt.FindString(priceString)
	someString := string(priceByte)

	// Make Clean "</span>" from "195,000,000</span>""
	someString2 := strings.Replace(someString, "</span>", "", -1)

	// Make Clean all "," in "195,000,000"
	priceInString := strings.Replace(someString2, ",", "", -1)

	// Make int format.
	price, _ := strconv.Atoi(priceInString)
	fmt.Println(price)
	return price
}
