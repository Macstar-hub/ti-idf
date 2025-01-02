package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
)

const dockPath string = "./test/sampleFileTrim.txt"
const stopWordFilePath string = "./configs/stopWords.txt"

func main() {
	fmt.Println("Hello from TF-IDF project.")
	DocReader(dockPath)
	DocCleanUp(stopWordFilePath)
}

func DocReader(fileName string) string {
	file, err := os.ReadFile(fileName)
	if err != nil {
		fmt.Printf("Cannot read file name: %v with error: %v \n", fileName, err)
	}

	fmt.Println("File size: ", len(file))
	return string(file)
}

func DocCleanUp(stopWordFilePath string) string {
	var cleanContent string
	stopWordFile2, _ := os.Open(stopWordFilePath)
	uncleanContent := DocReader(dockPath)
	uncleanContentLowerCase := strings.ToLower(uncleanContent)

	Lines := bufio.NewScanner(stopWordFile2)
	for Lines.Scan() {
		fmt.Println(Lines.Text())
		stopWordRegex := regexp.MustCompile(Lines.Text())
		fmt.Println(stopWordRegex)
		cleanContent = stopWordRegex.ReplaceAllString(uncleanContentLowerCase, " ")
		uncleanContentLowerCase = cleanContent
	}
	if err := Lines.Err(); err != nil {
		log.Fatal(err)
	}

	cleanData := regexp.MustCompile(`\.\s*`).ReplaceAllString(uncleanContentLowerCase, "\n")
	fmt.Println(uncleanContentLowerCase)
	WordCount(string(uncleanContentLowerCase))
	fmt.Println(cleanData)
	return cleanData
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
