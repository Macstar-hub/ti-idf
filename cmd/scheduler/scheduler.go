package main

import (
	"log"
	"sync"
	mysqlconnector "tf-idf/cmd/mysql"
	"tf-idf/cmd/telegram"
	"time"
)

func main() {

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go cronjob(1, APIPriceUpdateTask, wg)
	wg.Wait()
}

func cronjob(schedule int, functionName func(), wg *sync.WaitGroup) {
	for {
		go functionName()
		time.Sleep(time.Duration(schedule) * time.Minute)
		// wg.Done()
	}
}

func APIPriceUpdateTask() {
	log.Println("Update api calls from tjgu and cache.")
	telegram.GetCoinPrice()
	mysqlconnector.UpdatePrice()

}

func DBPriceUpdateTask() {
	log.Println("Update DB's records.")
	mysqlconnector.UpdatePrice()
}
