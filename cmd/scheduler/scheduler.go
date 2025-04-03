package main

import (
	"log"
	mysqlconnector "tf-idf/cmd/mysql"
	"tf-idf/cmd/telegram"
)

func main() {
	// Get price from mysql
	log.Println("Update cache and database is running .")
	telegram.GetCoinPrice()
	mysqlconnector.UpdatePrice()
}
