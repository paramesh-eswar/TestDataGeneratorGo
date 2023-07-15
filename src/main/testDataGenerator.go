package main

import (
	"fmt"
	"strconv"
	"sync"
	"time"
)

var output string

func testDataGenerator() string {
	startTime := time.Now()
	output += "\n----------------------------------------------------------------\n"
	output += "Test data generation is in progress ...\n"

	// test data generation logic starts
	tdgChannel := make(chan string, 1000)
	numOfReaders := 10
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
		if rowCount <= numOfRows-10 {
			endCount = rowCount + 10
		} else {
			endCount = numOfRows + 1
		}
		go sendeRecord(tdgChannel, &senderGroup, rowCount, endCount, 1)
		rowCount = endCount
		if rowCount >= numOfRows {
			break
		}
		// }
	}

	go func() {
		senderGroup.Wait()
		close(tdgChannel)
	}()

	readerGroup.Wait()
	// test data generation logic ends

	endTime := time.Now()
	output += "Time taken to generate test data: " + strconv.FormatInt(((endTime.UnixMilli()-startTime.UnixMilli())/1000), 10) + " sec\n"
	output += "Data generated successfully!"
	return output
}

func sendeRecord(tdgChannel chan string, wg *sync.WaitGroup, rowCount, endCount int, threadNum int) {
	// wg.Add(1)
	fmt.Println("sender-", threadNum, ":", rowCount, ":", endCount)
	for i := rowCount; i < endCount; i++ {
		// time.Sleep(1 * time.Second)
		tdgChannel <- strconv.Itoa(i)
	}
	defer wg.Done()
}

func readRecord(tdgChannel chan string, rg *sync.WaitGroup, readerNum int) {
	fmt.Println("reader", readerNum)
	for message := range tdgChannel {
		// time.Sleep(2 * time.Second)
		fmt.Println("reader", readerNum, ":", message)
	}
	defer rg.Done()
}
