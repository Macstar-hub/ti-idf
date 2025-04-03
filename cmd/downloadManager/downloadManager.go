package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	startTime := time.Now()

	totalSize := chunkManager("http://185.81.97.192")
	numWorker := makeWorkerPool()

	fmt.Println("Total Size is: ", totalSize, "And Also All Worker Number: ", numWorker)
	downloadURL("http://185.81.97.192", 0, 100, 1)

	fmt.Println("Download latency: ", time.Since(startTime))
}

func chunkManager(url string) int {
	response, err := http.Head(url)
	if err != nil {
		log.Println("Cannot make head webpage with: ", err)
	}
	contentLeanth := int(response.ContentLength)

	return int(contentLeanth)
}

func makeWorkerPool() int {
	worker := 10

	return worker
}

func downloadURL(url string, startRange int, endRange int, index int) {
	startTime := time.Now()

	client := &http.Client{}
	dataRange := fmt.Sprintf("bytes=%d-%d", startRange, endRange)

	request, errorReq := http.NewRequest("GET", url, nil)
	if errorReq != nil {
		log.Println("Cannot make request with error: ", errorReq)
	}

	request.Header.Add("Range", dataRange)

	response, errorResponse := client.Do(request)

	if errorResponse != nil {
		log.Println("Cannot make response: ", errorResponse)
	}

	body, _ := ioutil.ReadAll(response.Body)

	fileName, errCreateFile := os.Create(fmt.Sprintf("/tmp/%d.txt", index))
	if errCreateFile != nil {
		log.Println("Cannot make file with: ", errCreateFile)
	}
	// Make write all contenct to temp file:
	errOSWrite := os.WriteFile(fileName.Name(), body, 0644)
	if errOSWrite != nil {
		log.Println("Cannot write to file with error: ", errOSWrite)
	}

	fmt.Println("Download Latency: ", time.Since(startTime))
}

/*
	All tjgu IPs:
	92.246.136.221
	213.176.127.46
	92.246.136.221
	5.75.204.214
*/
