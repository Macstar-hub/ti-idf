package main

import (
	"crypto/tls"
	b64 "encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"regexp"
	"strconv"
	"strings"
	"syscall"

	mysqlconnector "tf-idf/cmd/mysql"

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
	var totalSquarePrice int
	var persquarPrice int
	// _, _, maskanURL := httpGet("https://divar.ir/s/tehran/buy-residential/ahang?size=65-80", "maskanurls")
	// _, _, maskanURL := httpGet("https://divar.ir/s/tehran/buy-apartment/west-tehran-pars?size=60-70", "maskanurls")
	_, _, maskanURL := httpGet("https://divar.ir/s/tehran/buy-residential/majid-abad?size=65-80", "maskanurls")
	ids, links, allLinks := mysqlconnector.SelectHousePrice()

	if len(ids) < 2 && allLinks == 0 {
		// Condition for make new table with fresh records.
		for id := 0; id < len(maskanURL); id++ {
			fmt.Println("URL: ", maskanURL[id])
			link := b64.StdEncoding.EncodeToString([]byte(maskanURL[id]))
			mysqlconnector.InsertHousePrice(int64(id), link, 0, 0)
		}
		// Condition for countinuse crawling price.
	} else {
		for i := 0; i < len(ids); i++ {
			maskanURL := fmt.Sprintf("%s", links[i])
			_, MaskanPrice, _ := httpGet(maskanURL, "maskan")

			if MaskanPrice == nil {
				fmt.Println("Cannot retrive price form divar site ...: ")
				break
			} else {
				totalSquarePrice = MaskanPrice[1]
				persquarPrice = MaskanPrice[0]
			}

			fmt.Println(persquarPrice, totalSquarePrice)

			idsInt := ids[i]
			mysqlconnector.UpdateHousePrice(idsInt, totalSquarePrice, persquarPrice)
			fmt.Println(MaskanPrice)
		}
	}
	if len(ids) == 0 && allLinks != 0 {
		fmt.Println("All links scraped and Done !")
		signalExit()
	}
}

func signalExit() {
	os.Exit(1) // Exit the Go program
}

func signalHandeler() {
	quite := make(chan os.Signal, 1)
	signal.Notify(quite, syscall.SIGINT, syscall.SIGTERM)
	s := <-quite
	fmt.Println("Termination: ", s)
}

func httpGet(url string, priceType string) (int, []int, []string) {
	var price int
	var maskanPrice []int
	var maskanURLS []string

	netClient := customHttpClient()
	// Add request example:
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Println(err)
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	log.Println("URL: ", url)
	// Add redponse section:
	responseByte, err := netClient.Do(req)

	httpErrorHandeler(err)
	responeBody, err := ioutil.ReadAll(responseByte.Body)
	byteReadErrorHandelete(err)

	responseString := string(responeBody)

	if priceType == "maskanurls" {
		maskanURLS = maskanPriceURL(responseString)
	}
	if priceType == "maskan" {
		maskanPrice = findMaskan(responseString)
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
			priceSlice = append(priceSlice, persianStepAsInt)
		}
	}
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

func maskanPriceURL(html string) []string {

	regexInit := regexp.MustCompile("https://divar.ir/v/.*")
	regexSec := regexp.MustCompile("\"url.*")
	var step2 string
	var step5 []string

	step1 := regexInit.FindAllString(string(html), 20000)

	for i := 0; i < len(step1); i++ {
		step2 = strings.ReplaceAll(step1[i], ",", "\n")
		step3 := regexSec.FindAllString(step2, 1000)

		fmt.Println("All houses founded: ", len(step3))

		for j := 0; j < len(step3); j++ {
			step4 := strings.Replace(step3[j], "\"", "", -1)
			fmt.Println(step4)
			stepBeracket := strings.Replace(step4, "}", "", -1)
			fmt.Println(stepBeracket)
			step5 = append(step5, fmt.Sprintf("%s", strings.Replace(stepBeracket, "url:", "", -1)))
		}
	}

	return step5

}
