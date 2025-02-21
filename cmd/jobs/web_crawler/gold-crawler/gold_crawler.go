package main

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	// "time"
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

	fmt.Println(price)

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

	if priceType == "maskanurls" {
		maskanURLS = maskanPriceURL(responseString)
	}
	if priceType == "maskan" {
		maskanPrice = findMaskan(responseString)
	} else {
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
			step5 = append(step5, fmt.Sprintf("%s", strings.Replace(step4, "url:", "", -1)))
		}
	}

	return step5
	// maskanURL := []string{"https://divar.ir/v/فروش-آپارتمان-۷۴-متری-۲-خوابه-در-تهرانپارس-غربی/wZV6v7KC", "https://divar.ir/v/اپارتمان-73-متری-فلکه-سوم-تهرانپارس/wZnSpYPw", "https://divar.ir/v/سالن-پرده-خور-تخلیه-۷۶-متر-۲-خ-فول-بازسازی-عروسکی/wZj22PcD", "https://divar.ir/v/۷۰-متری-۲-خوابه-بهار-تهرانپارس/wZnuHk0f", "https://divar.ir/v/۷۷متر-فول-امکانات-تک-واحدی/wZn226UU", "https://divar.ir/v/فروش-اپارتمان-۶۹-متری-دو-خوابه-در-تهرانپارس-غربی/wZfu3BXE", "https://divar.ir/v/80-متر-2-خ-فول-تراوتن-اخوت-و-طاهری-رو-به-آفتاب/wZnS0WeK", "https://divar.ir/v/۷۷متر-با-اسانسور-وانباری-غرب-فلکه-سوم/wZnu2BkK", "https://divar.ir/v/۷۶-دو-خوابه-فول-امکانات-تاپ-لوکیشن/wZniFE2P", "https://divar.ir/v/۷۷متر-فول-امکانات-بالکن-سالن-پرده-خور-نماتراورتن/wZjCXo5c", "https://divar.ir/v/۶۶-متر-۲-خوابه-بهترین-فرعی-غربیها/wZ1501An", "https://divar.ir/v/68-متر-2خواب-طبقه-اول-مسکن-بزرگ-210/wZnqFPTD", "https://divar.ir/v/اپارتمان-۶۷متری-۱خوابه-تیرانداز/wZdK_AvX", "https://divar.ir/v/پیش-فروش-آپارتمان-78متری-در-قشم-شهرک-مهروماه/wZnujFRf", "https://divar.ir/v/فروش-آپارتمان-۷۵-متری-۲-خوابه-فول-امکانات/wZnOBnSP", "https://divar.ir/v/۷۴-متر-دو-خواب-تهرانپارس-غربی/wZB6kq9G", "https://divar.ir/v/۶۶-متر-تکخواب-تخلیه-فول-امکانات/wZGwkj4F", "https://divar.ir/v/۸۰متری-فول-۴ساله-فلکه-اول/wZnSCUAV", "https://divar.ir/v/آپارتمان-۷۷-متری-۲-خوابه/wZnGCDPj", "https://divar.ir/v/۷۲-متر-۲-خواب-محدوده-بهار/wZnuhi3h", "https://divar.ir/v/۷۳متر-۲-خواب-تــک-واحــدی-فـول-بـازسـازی/wZnqxMYI", "https://divar.ir/v/فروش-آپارتمان-۶۸-متری-۲-خوابه-در-تهرانپارس-غربی/wZh2vhO2", "https://divar.ir/v/فروش-آپارتمان-۷۷-متری-۲-خوابه-در-تهرانپارس-غربی/wZXGoRyn", "https://divar.ir/v/۶۶متر-یک-خوابه-تک-واحدی-فلکه-چهارم-تک-بازدید/wZi6fVom"}

}
