package main

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/mavihq/persian"
)

type Price struct {
	Dollar       int
	SekkeTamam   int
	SekketGhadim int
	SekkehNim    int
	RobeSekke    int
	Maskan       int
}

func main() {
	var price Price

	usdPrice, _ := httpGet("https://www.tgju.org/profile/price_dollar_rl", "priceGold")
	price.Dollar = usdPrice

	sekkeTamamPrice, _ := httpGet("https://www.tgju.org/profile/sekee", "priceGold")
	price.SekkeTamam = sekkeTamamPrice

	sekkeGhadimPrice, _ := httpGet("https://www.tgju.org/profile/sekeb", "priceGold")
	price.SekketGhadim = sekkeGhadimPrice

	SekkehNimPrice, _ := httpGet("https://www.tgju.org/profile/nim", "priceGold")
	price.SekkehNim = SekkehNimPrice

	SekkehRobePrice, _ := httpGet("https://www.tgju.org/profile/rob", "priceGold")
	price.RobeSekke = SekkehRobePrice

	maskanURL := maskanPriceURL()
	for i := 0; i < len(maskanPriceURL()); i++ {
		_, MaskanPrice := httpGet(maskanURL[i], "maskan")
		fmt.Println(MaskanPrice)
	}

	// _, maskanPrice := httpGet("https//diva.i/v/فروش-آپارتمان-۶۰-متری-۲-خوابه-در-تهرانپارس-غربی/wZUivCI3", "maskan")
	// fmt.Println(maskanPrice)
	fmt.Println(price)
}

func httpGet(url string, priceType string) (int, []int) {
	var price int
	var maskanPrice []int
	netClient := customHttpClient()

	responseByte, err := netClient.Get(url)

	httpErrorHandeler(err)

	responeBody, err := ioutil.ReadAll(responseByte.Body)
	byteReadErrorHandelete(err)

	responseString := string(responeBody)

	if priceType == "maskan" {
		maskanPrice = findMaskan(responseString)
	} else {
		_, price = findSekkeTamam(responseString)
	}

	responseByte.Body.Close()

	return price, maskanPrice
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

func findMaskan(html string) []int {
	var priceSlice []int

	regexInt, _ := regexp.Compile("\"value.*")
	regexSec := regexp.MustCompile("[^0-9]+")
	htmlString := string(html)
	priceStep1Cleanup := strings.Replace(htmlString, ",", "\n", -1)
	priceStep2Cleanup := regexInt.FindAllString(priceStep1Cleanup, 10000)
	persianCharecter := "\u062A\u0648\u0645\u0627\u0646"
	for i := 0; i < len(priceStep2Cleanup); i++ {
		if strings.Contains(priceStep2Cleanup[i], persianCharecter) {
			persianStep1 := persian.ToEnglishDigits(priceStep2Cleanup[i])
			persianStep2 := persian.SwitchToEnglishKey(persianStep1)
			persianStep3 := strings.Replace(persianStep2, "\"value\":\"", "", -1)
			persianStep4 := strings.Replace(persianStep3, "j,lhk\"", "", -1)
			persianStepAsString := regexSec.ReplaceAllString(persianStep4, "")
			persianStepAsInt, _ := strconv.Atoi(persianStepAsString)
			// fmt.Println(persianStepAsInt)
			priceSlice = append(priceSlice, persianStepAsInt)
		}
	}
	// fmt.Println(priceSlice)
	return priceSlice
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

func maskanPriceURL() []string {
	maskanURL := []string{"https://divar.ir/v/اپارتمان-۱۰۰-متری-۲-خوابه/wZlevVAA", "https://divar.ir/v/۹۶مترفلکه-اول-۱۶۲غربی-کلید-نخورده-برند/wZzZ6BDM", "https://divar.ir/v/فروش-و-معاوضه-۸-دستگاه-آپارتمان-۹۰-متری/wZfyUrJe", "https://divar.ir/v/۸۱متر-فول-محله-بهار-تهرانپارس/wZlicEoQ"}
	return maskanURL
}
