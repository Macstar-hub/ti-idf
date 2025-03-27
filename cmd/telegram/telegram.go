package telegram

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"sync"
	mysqlconnector "tf-idf/cmd/mysql"
	"time"
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

func GetCoinPrice() *Price {

	startTime := time.Now()

	priceSlice := []int{}

	UrlList := []string{
		"https://www.tgju.org/profile/price_dollar_rl",
		"https://www.tgju.org/profile/sekee",
		"https://www.tgju.org/profile/sekeb",
		"https://www.tgju.org/profile/nim",
		"https://www.tgju.org/profile/rob",
		"https://www.tgju.org/profile/geram18",
		"https://www.tgju.org/profile/gold_mini_size",
	}

	Symbol := []string{
		"Dollar",
		"SekkeTamam",
		"SekketGhadim",
		"SekkehNim",
		"RobeSekke",
		"Gold18",
		"GoldDast2",
	}

	p := new(Price)
	responseChannel := make(chan Price, 1024)
	wg := &sync.WaitGroup{}

	for i, url := range UrlList {
		wg.Add(1)
		go p.getPrice(url, "priceGold", responseChannel, wg, Symbol[i])
	}

	wg.Wait()
	close(responseChannel)

	finalPrice := new(Price)

	for responseChann := range responseChannel {
		finalPrice.Dollar = responseChann.Dollar + finalPrice.Dollar
		finalPrice.SekkeTamam = responseChann.SekkeTamam + finalPrice.SekkeTamam
		finalPrice.SekketGhadim = responseChann.SekketGhadim + finalPrice.SekketGhadim
		finalPrice.SekkehNim = responseChann.SekkehNim + finalPrice.SekkehNim
		finalPrice.RobeSekke = responseChann.RobeSekke + finalPrice.RobeSekke
		finalPrice.Gold18 = responseChann.Gold18 + finalPrice.Gold18
		finalPrice.GoldDast2 = responseChann.GoldDast2 + finalPrice.GoldDast2
	}

	log.Println("From channel is: ", finalPrice)
	log.Println("Total latency: ", time.Since(startTime))

	priceSlice = append(priceSlice, finalPrice.Gold18)
	priceSlice = append(priceSlice, finalPrice.SekkeTamam)
	priceSlice = append(priceSlice, finalPrice.SekketGhadim)
	priceSlice = append(priceSlice, finalPrice.SekkehNim)

	mysqlconnector.UpdatePrice(priceSlice)

	return finalPrice

}

func (p Price) getPrice(url string, priceType string, responceChannel chan Price, wg *sync.WaitGroup, getPriceType string) {
	defer wg.Done()

	var price int
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

	if getPriceType == "Dollar" {
		p.Dollar = price
	}
	if getPriceType == "SekkeTamam" {
		p.SekkeTamam = price
	}
	if getPriceType == "SekketGhadim" {
		p.SekketGhadim = price
	}
	if getPriceType == "SekkehNim" {
		p.SekkehNim = price
	}
	if getPriceType == "RobeSekke" {
		p.RobeSekke = price
	}
	if getPriceType == "Gold18" {
		p.Gold18 = price
	}
	if getPriceType == "GoldDast2" {
		p.GoldDast2 = price
	}

	responceChannel <- p

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

	// Enable just for debug:
	// fmt.Println(price)

	return price
}
