package main

import (
	"fmt"
	"log"
	"os/exec"
	"sync"
	mysqlconnector "tf-idf/cmd/mysql"
	"tf-idf/cmd/telegram"
	"time"
)

func main() {

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go cronjob(15, APIPriceUpdateTask, wg, time.Minute)
	go cronjob(23, HousePriceAnalyze, wg, time.Hour)
	wg.Wait()
}

func cronjob(schedule int, functionName func(wg *sync.WaitGroup), wg *sync.WaitGroup, unit time.Duration) {
	for {
		go functionName(wg)
		time.Sleep(time.Duration(schedule) * unit)
	}
}

func APIPriceUpdateTask(wg *sync.WaitGroup) {
	log.Println("Update api calls from tjgu and cache.")
	telegram.GetCoinPrice()
	mysqlconnector.UpdatePrice()

}

func HousePriceAnalyze(wg *sync.WaitGroup) {
	log.Println("House Price Analyze Run.")
	output, err := exec.Command("/bin/zsh", "../jobs/run.sh").Output() // Make sure change all directory in run.sh scripts.
	if err != nil {
		log.Println("Cannot make house price with error: ", err)
	}
	fmt.Printf("output is %s\n", output)
}
