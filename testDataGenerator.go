package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
)

var output string
var dataWriter bufio.Writer
var processedRowCount int

func testDataGenerator() string {
	startTime := time.Now()
	output += "\n----------------------------------------------------------------\n"
	output += "Test data generation is in progress ...\n"

	// test data generation logic starts
	progressBar.SetValue(0)
	processedRowCount = 1
	var rowBuilder strings.Builder
	for _, jsonAttr := range metaDataJson {
		rowBuilder.WriteString(jsonAttr["name"].(string) + ",")
	}
	headerRow := rowBuilder.String()
	headerRow = headerRow[:len(headerRow)-1]
	fmt.Println(headerRow)
	outputFilePath := generateOutputFileName(metadataFileName)
	fileToWrite, err := os.OpenFile(outputFilePath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return "Failed to create output file: " + outputFilePath
	}
	defer fileToWrite.Close()
	dataWriter = *bufio.NewWriter(fileToWrite)
	dataWriter.WriteString(headerRow)

	tdgChannel := make(chan string, 1000)
	const numOfReaders = 10
	// numOfThreads := 10
	rowCount := 1
	var endCount int

	var readerGroup sync.WaitGroup
	readerGroup.Add(numOfReaders)
	for i := 1; i <= numOfReaders; i++ {
		go readRecord(tdgChannel, &readerGroup, i)
	}

	var senderGroup sync.WaitGroup
	for rowCount <= numOfRows {
		senderGroup.Add(1)
		// for i := 0; i < numOfThreads; i++ {
		if rowCount <= numOfRows-500 {
			endCount = rowCount + 500
		} else {
			endCount = numOfRows + 1
		}
		go sendeRecord(tdgChannel, &senderGroup, rowCount, endCount)
		rowCount = endCount
		if rowCount > numOfRows {
			break
		}
		// }
	}

	go func() {
		senderGroup.Wait()
		close(tdgChannel)
	}()

	readerGroup.Wait()
	dataWriter.Flush()
	progressBar.SetValue(1)
	// test data generation logic ends

	endTime := time.Now()
	output += "Time taken to generate test data: " + strconv.FormatInt(((endTime.UnixMilli()-startTime.UnixMilli())/1000), 10) + " sec\n"
	output += "Output file location: " + outputFilePath + "\n"
	output += "Data generated successfully!"
	return output
}

func sendeRecord(tdgChannel chan string, wg *sync.WaitGroup, rowCount, endCount int) {
	// wg.Add(1)
	fmt.Println("sender:", rowCount, ":", endCount-1)
	for i := rowCount; i < endCount; i++ {
		tdgChannel <- strconv.Itoa(i)
	}
	defer wg.Done()
}

func readRecord(tdgChannel chan string, rg *sync.WaitGroup, readerNum int) {
	for message := range tdgChannel {
		fmt.Println("reader", processedRowCount, ":", message)
		PercentageCompleted = (float64(processedRowCount) / float64(numOfRows))
		// externalFloat.Reload()
		progressBar.SetValue(PercentageCompleted)
		processedRowCount++
	}
	defer rg.Done()
}

func generateOutputFileName(inputFilePath string) string {
	inputFileDir := filepath.Dir(inputFilePath)
	inputFileName := filepath.Base(inputFilePath)
	inputFileExt := filepath.Ext(inputFileName)
	inputFileNameWithoutExt := strings.TrimSuffix(inputFileName, inputFileExt) + "_output"
	outputFileName := inputFileNameWithoutExt + inputFileExt
	finalOutputFileName := filepath.Join(inputFileDir, outputFileName)
	fmt.Println("Output file path:", finalOutputFileName)
	return finalOutputFileName
}
