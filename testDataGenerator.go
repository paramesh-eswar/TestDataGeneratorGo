package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/brianvoe/gofakeit/v6"
)

var output string
var processedRowCount int

func testDataGenerator() string {
	var dataWriter bufio.Writer
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
	headerRow = headerRow[:len(headerRow)-1] + "\n"
	fmt.Println(headerRow)
	outputFilePath := generateOutputFileName(metadataFileName)
	fileToWrite, err := os.OpenFile(outputFilePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return "Failed to create output file: " + outputFilePath
	}
	defer fileToWrite.Close()
	dataWriter = *bufio.NewWriter(fileToWrite)
	dataWriter.WriteString(headerRow)
	dataWriter.Flush()

	tdgChannel := make(chan string, 1000)
	rowCount := 1
	var endCount int

	var readerGroup sync.WaitGroup
	readerGroup.Add(1)
	go readRecord(tdgChannel, &readerGroup, dataWriter)

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
	// dataWriter.Flush()
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
	var rowBuilder strings.Builder
	for i := rowCount; i < endCount; i++ {
		for _, jsonAttr := range metaDataJson {
			switch jsonAttr["datatype"] {
			case "number":
				switch dataGenType[jsonAttr["name"].(string)] {
				case DEFAULT:
					rowBuilder.WriteString(strings.TrimSpace(jsonAttr["default_value"].(string)) + ",")
				case NATURAL_SEQ:
					rowBuilder.WriteString(fmt.Sprint(i) + ",")
				case SEQ_IN_RANGE:
					numRange := strings.Split(strings.TrimSpace(jsonAttr["range"].(string)), "~")
					numMin, _ := strconv.Atoi(numRange[0])
					rowBuilder.WriteString(fmt.Sprint(numMin+i) + ",")
				case DUP_IN_RANGE:
					numRange := strings.Split(strings.TrimSpace(jsonAttr["range"].(string)), "~")
					numMin, _ := strconv.Atoi(numRange[0])
					numMax, _ := strconv.Atoi(numRange[1])
					rowBuilder.WriteString(fmt.Sprint(numMin+rand.Intn(numMax-numMin)) + ",")
				case RANDOM:
					rowBuilder.WriteString(fmt.Sprint(gofakeit.Number(1, numOfRows)) + ",")
					// rand.NewSource(time.Now().UnixNano())
					// rowBuilder.WriteString(fmt.Sprint(1+rand.Intn(numOfRows-1)) + ",")
				}
			case "text":
				switch dataGenType[jsonAttr["name"].(string)] {
				case DEFAULT:
					rowBuilder.WriteString(strings.TrimSpace(jsonAttr["default_value"].(string)) + ",")
				case NATURAL_SEQ:
				case SEQ_IN_RANGE:
				case DUP_IN_RANGE:
				case RANDOM:
				}
			case "float":
				switch dataGenType[jsonAttr["name"].(string)] {
				case DEFAULT:
					rowBuilder.WriteString(strings.TrimSpace(jsonAttr["default_value"].(string)) + ",")
				case NATURAL_SEQ:
				case SEQ_IN_RANGE:
				case DUP_IN_RANGE:
				case RANDOM:
				}
			case "date":
				switch dataGenType[jsonAttr["name"].(string)] {
				case DEFAULT:
					rowBuilder.WriteString(strings.TrimSpace(jsonAttr["default_value"].(string)) + ",")
				case NATURAL_SEQ:
				case SEQ_IN_RANGE:
				case DUP_IN_RANGE:
				case RANDOM:
				}
			case "gender":
				switch dataGenType[jsonAttr["name"].(string)] {
				case DEFAULT:
					rowBuilder.WriteString(strings.TrimSpace(jsonAttr["default_value"].(string)) + ",")
				case NATURAL_SEQ:
				case SEQ_IN_RANGE:
				case DUP_IN_RANGE:
				case RANDOM:
				}
			case "boolean":
				switch dataGenType[jsonAttr["name"].(string)] {
				case DEFAULT:
					rowBuilder.WriteString(strings.TrimSpace(jsonAttr["default_value"].(string)) + ",")
				case NATURAL_SEQ:
				case SEQ_IN_RANGE:
				case DUP_IN_RANGE:
				case RANDOM:
				}
			case "ssn":
				switch dataGenType[jsonAttr["name"].(string)] {
				case DEFAULT:
					rowBuilder.WriteString(strings.TrimSpace(jsonAttr["default_value"].(string)) + ",")
				case NATURAL_SEQ:
				case SEQ_IN_RANGE:
				case DUP_IN_RANGE:
				case RANDOM:
				}
			case "email":
				switch dataGenType[jsonAttr["name"].(string)] {
				case DEFAULT:
					rowBuilder.WriteString(strings.TrimSpace(jsonAttr["default_value"].(string)) + ",")
				case NATURAL_SEQ:
				case SEQ_IN_RANGE:
				case DUP_IN_RANGE:
				case RANDOM:
				}
			case "phonenumber":
				switch dataGenType[jsonAttr["name"].(string)] {
				case DEFAULT:
					rowBuilder.WriteString(strings.TrimSpace(jsonAttr["default_value"].(string)) + ",")
				case NATURAL_SEQ:
				case SEQ_IN_RANGE:
				case DUP_IN_RANGE:
				case RANDOM:
				}
			case "aadhar":
				switch dataGenType[jsonAttr["name"].(string)] {
				case DEFAULT:
					rowBuilder.WriteString(strings.TrimSpace(jsonAttr["default_value"].(string)) + ",")
				case NATURAL_SEQ:
				case SEQ_IN_RANGE:
				case DUP_IN_RANGE:
				case RANDOM:
				}
			case "creditcard":
				switch dataGenType[jsonAttr["name"].(string)] {
				case DEFAULT:
					rowBuilder.WriteString(strings.TrimSpace(jsonAttr["default_value"].(string)) + ",")
				case NATURAL_SEQ:
				case SEQ_IN_RANGE:
				case DUP_IN_RANGE:
				case RANDOM:
				}
			case "zipcode":
				switch dataGenType[jsonAttr["name"].(string)] {
				case DEFAULT:
					rowBuilder.WriteString(strings.TrimSpace(jsonAttr["default_value"].(string)) + ",")
				case NATURAL_SEQ:
				case SEQ_IN_RANGE:
				case DUP_IN_RANGE:
				case RANDOM:
				}
			case "uuid":
				switch dataGenType[jsonAttr["name"].(string)] {
				case DEFAULT:
					rowBuilder.WriteString(strings.TrimSpace(jsonAttr["default_value"].(string)) + ",")
				case NATURAL_SEQ:
				case SEQ_IN_RANGE:
				case DUP_IN_RANGE:
				case RANDOM:
				}
			case "ipaddress":
				switch dataGenType[jsonAttr["name"].(string)] {
				case DEFAULT:
					rowBuilder.WriteString(strings.TrimSpace(jsonAttr["default_value"].(string)) + ",")
				case NATURAL_SEQ:
				case SEQ_IN_RANGE:
				case DUP_IN_RANGE:
				case RANDOM:
				}
			case "timestamp":
				switch dataGenType[jsonAttr["name"].(string)] {
				case DEFAULT:
					rowBuilder.WriteString(strings.TrimSpace(jsonAttr["default_value"].(string)) + ",")
				case NATURAL_SEQ:
				case SEQ_IN_RANGE:
				case DUP_IN_RANGE:
				case RANDOM:
				}
			}
		}
		dataRow := rowBuilder.String()
		dataRow = dataRow[:len(dataRow)-1] + "\n"
		tdgChannel <- dataRow
		rowBuilder.Reset()
	}
	defer wg.Done()
}

func readRecord(tdgChannel chan string, rg *sync.WaitGroup, dataWriter bufio.Writer) {
	for message := range tdgChannel {
		fmt.Println("reader", ":", message)
		dataWriter.WriteString(message)
		dataWriter.Flush()
		// refresh the progress bar to show the current progress
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
	outputFileName := inputFileNameWithoutExt + ".csv"
	finalOutputFileName := filepath.Join(inputFileDir, outputFileName)
	fmt.Println("Output file path:", finalOutputFileName)
	return finalOutputFileName
}
