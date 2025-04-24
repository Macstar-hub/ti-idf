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
	go cronjob(15, APIPriceUpdateTask, wg)
	wg.Wait()
}

func cronjob(schedule int, functionName func(), wg *sync.WaitGroup) {
	for {
		functionName()
		time.Sleep(time.Duration(schedule) * time.Minute)
		// wg.Done()
	}
}

func APIPriceUpdateTask() {
	log.Println("Update cache and database is running .")
	telegram.GetCoinPrice()
	mysqlconnector.UpdatePrice()

}

// func DBPriceUpdateTask() {
// 	log.Println("Hello from cronjob2.")
// 	mysqlconnector.UpdatePrice()

// }
