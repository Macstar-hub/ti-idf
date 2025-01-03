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
	TFcount(string(uncleanContentLowerCase))
	IDFcount(cleanData)
	fmt.Println("all sentences: \n", cleanData)
	return cleanData
}

func TFcount(fileText string) {
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

func IDFcount(cleanContent string) {
	var documentNumber int
	documentLines := strings.Split(string(cleanContent), "\n")

	// Find total document numbers:
	for index, _ := range documentLines {
		documentNumber = index
	}

	// Iterate over document:
	for _, documentLine := range documentLines {
		documentLineNew := strings.Fields(documentLine)
		// fmt.Printf("All words in document %v are %v \n", index, documentLineNew)
		var wordSlices []string

		for _, word := range documentLineNew {
			fmt.Println("Debug each word in document line: ", word)
			wordSlices = append(wordSlices, word)
			fmt.Println("Debug word slices: ", wordSlices)
		}

		for index := range len(wordSlices) {
			fmt.Println("$$$$$$", wordSlices[index])
		}
	}

	fmt.Println("Debug from IDF count: ", documentNumber)
}
