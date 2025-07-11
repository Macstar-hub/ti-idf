package main

import (
	"log"

	"os/exec"
	"sync"
	logger "tf-idf/cmd/logger"
	mysqlconnector "tf-idf/cmd/mysql"
	"tf-idf/cmd/telegram"
	"time"
)

const (
	logFilePath = "../../logs/scheduler/"
	logPrefix   = ".log"
)

/*
- Handle if a job got error and make it continue ...
*/

func main() {

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go cronjob(15, APIPriceUpdateTask, wg, time.Minute)
	// go cronjob(23, HousePriceAnalyze, wg, time.Hour)
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
	logger.Logger(logFilePath, logPrefix, "Gold API Update Completed.", "debug")
}

func HousePriceAnalyze(wg *sync.WaitGroup) {
	log.Println("House Price Analyze Run.")
	output, err := exec.Command("/bin/zsh", "../jobs/run.sh").Output() // Make sure change all directory in run.sh scripts.
	if err != nil {
		log.Println("Cannot make house price with error: ", err)
	}
	log.Printf("output is %s\n", output)
	logger.Logger(logFilePath, logPrefix, string(output), "debug")
}
