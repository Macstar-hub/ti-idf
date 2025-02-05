package main

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

type Price struct {
	Dollar       int
	SekkeTamam   int
	SekketGhadim int
	SekkehNim    int
	RobeSekke    int
}

func main() {
	var price Price

	usdPrice := httpGet("https://www.tgju.org/profile/price_dollar_rl")
	price.Dollar = usdPrice

	sekkeTamamPrice := httpGet("https://www.tgju.org/profile/sekee")
	price.SekkeTamam = sekkeTamamPrice

	sekkeGhadimPrice := httpGet("https://www.tgju.org/profile/sekeb")
	price.SekketGhadim = sekkeGhadimPrice

	SekkehNimPrice := httpGet("https://www.tgju.org/profile/nim")
	price.SekkehNim = SekkehNimPrice

	SekkehRobePrice := httpGet("https://www.tgju.org/profile/rob")
	price.RobeSekke = SekkehRobePrice

	fmt.Println(price)
}

func httpGet(url string) int {
	netClient := customHttpClient()

	responseByte, err := netClient.Get(url)

	httpErrorHandeler(err)

	responeBody, err := ioutil.ReadAll(responseByte.Body)
	byteReadErrorHandelete(err)

	responseString := string(responeBody)

	_, price := findSekkeTamam(responseString)
	// fmt.Println(price)
	responseByte.Body.Close()

	return price
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

	return price
}
