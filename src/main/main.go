package main

import (
	"fmt"
	"strconv"
	"sync"
)

func init() {
	fmt.Println("=======================")
	fmt.Println("Test Data Generator")
	fmt.Println("=======================")
}

func reader(tdgChannel chan string, rg *sync.WaitGroup, readerNum int) {
	fmt.Println("reader", readerNum)
	for message := range tdgChannel {
		// time.Sleep(2 * time.Second)
		fmt.Println("reader", readerNum, ":", message)
	}
	rg.Done()
}

func main() {
	fmt.Println("Test Data Generator tool started ...")
	tdgChannel := make(chan string, 10)

	var rg sync.WaitGroup
	rg.Add(2)

	for i := 1; i <= 2; i++ {
		go reader(tdgChannel, &rg, i)
	}

	var sg sync.WaitGroup

	numOfRows := 100
	rowCount := 1
	numOfThreads := 3
	var endCount int

	for rowCount <= numOfRows {
		for i := 0; i < numOfThreads; i++ {
			if rowCount <= numOfRows-10 {
				endCount = rowCount + 10
			} else {
				endCount = numOfRows + 1
			}
			go sender(tdgChannel, &sg, rowCount, endCount, i)
			rowCount = endCount
			if rowCount >= numOfRows {
				break
			}
		}
	}

	go func() {
		sg.Wait()
		close(tdgChannel)
	}()
	rg.Wait()
}

func sender(tdgChannel chan string, wg *sync.WaitGroup, rowCount, encCount int, threadNum int) {
	wg.Add(1)
	fmt.Println("sender-", threadNum, ":", rowCount, ":", encCount)
	for i := rowCount; i < encCount; i++ {
		// time.Sleep(1 * time.Second)
		tdgChannel <- strconv.Itoa(i)
	}
	defer wg.Done()
}
