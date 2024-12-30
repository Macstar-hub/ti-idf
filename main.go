package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	fmt.Println("Hello from TF-IDF project.")
	DocReader("./test/sampleFile.txt")
}

func DocReader(fileName string) string {
	file, err := os.ReadFile(fileName)
	if err != nil {
		fmt.Printf("Cannot read file name: %v with error: %v \n", fileName, err)
	}

	fmt.Println("File size: ", len(file))
	WordCount(string(file))
	return string(file)
}

func WordCount(fileText string) {
	words := strings.Fields(fileText)
	fmt.Println("All words count: ", len(words))
	for _, word := range words {
		count := 0
		var ratio float64
		for _, targetWord := range words {
			if word == targetWord {
				count = count + 1
			}
		}
		ratio = float64(count) / float64(len(words))
		fmt.Printf("Word '%v': repeat with %v, and with ratio %v \n", word, count, ratio*100)
	}
}
