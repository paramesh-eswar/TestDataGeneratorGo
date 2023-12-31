package main

import (
	"bufio"
	"fmt"
	"math"
	"math/rand"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"golang.org/x/exp/slices"
)

var output string
var processedRowCount int
var mu sync.Mutex

func testDataGenerator() string {
	var dataWriter bufio.Writer
	startTime := time.Now()
	output += "\n----------------------------------------------------------------\n"
	output += "Test data generation for the row count: " + strconv.Itoa(numOfRows) + "\n"

	// test data generation logic starts
	progressBar.SetValue(0)
	processedRowCount = 1
	var rowBuilder strings.Builder
	for _, jsonAttr := range metaDataJson {
		rowBuilder.WriteString(jsonAttr["name"].(string) + DELIMITER)
	}
	headerRow := rowBuilder.String()
	headerRow = headerRow[:len(headerRow)-1]
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
	attributeLookups := make(map[string][]string)
	for rowCount <= numOfRows {
		senderGroup.Add(1)
		// for i := 0; i < numOfThreads; i++ {
		if rowCount <= numOfRows-1000 {
			endCount = rowCount + 1000
		} else {
			endCount = numOfRows + 1
		}
		go sendeRecord(tdgChannel, &senderGroup, rowCount, endCount, &attributeLookups)
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

func sendeRecord(tdgChannel chan string, wg *sync.WaitGroup, rowCount, endCount int, attrLookups *map[string][]string) {
	// wg.Add(1)
	fmt.Println("sender:", rowCount, ":", endCount-1)
	var rowBuilder strings.Builder
	for i := rowCount; i < endCount; i++ {
		rowBuilder.WriteString("\n")
		for _, jsonAttr := range metaDataJson {
			person := gofakeit.Person()
			switch jsonAttr["datatype"] {
			case "number":
				switch dataGenType[jsonAttr["name"].(string)] {
				case DEFAULT:
					rowBuilder.WriteString(strings.TrimSpace(jsonAttr["default_value"].(string)) + DELIMITER)
				case NATURAL_SEQ:
					rowBuilder.WriteString(fmt.Sprint(i) + DELIMITER)
				case SEQ_IN_RANGE:
					numRange := strings.Split(strings.TrimSpace(jsonAttr["range"].(string)), "~")
					numMin, _ := strconv.Atoi(numRange[0])
					rowBuilder.WriteString(fmt.Sprint(numMin+(i-1)) + DELIMITER)
				case DUP_IN_RANGE:
					numRange := strings.Split(strings.TrimSpace(jsonAttr["range"].(string)), "~")
					numMin, _ := strconv.Atoi(numRange[0])
					numMax, _ := strconv.Atoi(numRange[1])
					rand.NewSource(time.Now().UnixNano())
					rowBuilder.WriteString(fmt.Sprint(numMin+rand.Intn(numMax-numMin)) + DELIMITER)
				case RANDOM:
					rowBuilder.WriteString(fmt.Sprint(gofakeit.Number(1, numOfRows)) + DELIMITER)
					// rand.NewSource(time.Now().UnixNano())
					// rowBuilder.WriteString(fmt.Sprint(1+rand.Intn(numOfRows-1)) + DELIMITER)
				}
			case "text":
				switch dataGenType[jsonAttr["name"].(string)] {
				case DEFAULT:
					rowBuilder.WriteString(strings.TrimSpace(jsonAttr["default_value"].(string)) + DELIMITER)
				case NATURAL_SEQ:
					rowBuilder.WriteString(jsonAttr["name"].(string) + fmt.Sprint(i) + DELIMITER)
				case SEQ_IN_RANGE:
					textRange := jsonAttr["range"].([]interface{})
					rowBuilder.WriteString(textRange[i-1].(string) + DELIMITER)
				case DUP_IN_RANGE:
					textRange := jsonAttr["range"].([]interface{})
					rand.NewSource(time.Now().UnixNano())
					rowBuilder.WriteString(textRange[rand.Intn(len(textRange))].(string) + DELIMITER)
				case RANDOM:
					rowBuilder.WriteString(gofakeit.Word() + DELIMITER)
				}
			case "float":
				switch dataGenType[jsonAttr["name"].(string)] {
				case DEFAULT:
					rowBuilder.WriteString(strings.TrimSpace(jsonAttr["default_value"].(string)) + DELIMITER)
				case NATURAL_SEQ:
					scale, _ := strconv.Atoi(strings.TrimSpace(jsonAttr["scale"].(string)))
					rowBuilder.WriteString(strconv.FormatFloat(float64(i), 'f', scale, 64) + DELIMITER)
				case SEQ_IN_RANGE:
					scale, _ := strconv.Atoi(strings.TrimSpace(jsonAttr["scale"].(string)))
					floatRange := strings.Split(strings.TrimSpace(jsonAttr["range"].(string)), "~")
					floatMin, _ := strconv.ParseFloat(floatRange[0], 64)
					rowBuilder.WriteString(strconv.FormatFloat(floatMin+float64(i-1), 'f', scale, 64) + DELIMITER)
				case DUP_IN_RANGE:
					scale, _ := strconv.Atoi(strings.TrimSpace(jsonAttr["scale"].(string)))
					floatRange := strings.Split(strings.TrimSpace(jsonAttr["range"].(string)), "~")
					floatMin, _ := strconv.ParseFloat(floatRange[0], 64)
					floatMax, _ := strconv.ParseFloat(floatRange[1], 64)
					rand.NewSource(time.Now().UnixNano())
					rowBuilder.WriteString(strconv.FormatFloat(floatMin+rand.NormFloat64()*(floatMax-floatMin), 'f', scale, 64) + DELIMITER)
				case RANDOM:
					scale, _ := strconv.Atoi(strings.TrimSpace(jsonAttr["scale"].(string)))
					rowBuilder.WriteString(strconv.FormatFloat(gofakeit.Float64Range(0, float64(numOfRows)), 'f', scale, 64) + DELIMITER)
				}
			case "date":
				switch dataGenType[jsonAttr["name"].(string)] {
				case DEFAULT:
					rowBuilder.WriteString(strings.TrimSpace(jsonAttr["default_value"].(string)) + DELIMITER)
				case DUP_IN_RANGE:
					dateRange := strings.Split(strings.TrimSpace(jsonAttr["range"].(string)), "~")
					dateMin, _ := time.Parse(jsonAttr["date_format"].(string), dateRange[0])
					dateMax, _ := time.Parse(jsonAttr["date_format"].(string), dateRange[1])
					// duration := dateMax.Sub(dateMin)
					// randDuration := time.Duration(rand.Int63n(int64(duration)))
					// rowBuilder.WriteString(dateMin.Add(randDuration).Format(jsonAttr["date_format"].(string)) + DELIMITER)
					rowBuilder.WriteString(gofakeit.DateRange(dateMin, dateMax).Format(jsonAttr["date_format"].(string)) + DELIMITER)
				case RANDOM:
					rowBuilder.WriteString(gofakeit.Date().Format(jsonAttr["date_format"].(string)) + DELIMITER)
				}
			case "gender":
				switch dataGenType[jsonAttr["name"].(string)] {
				case DEFAULT:
					rowBuilder.WriteString(strings.TrimSpace(jsonAttr["default_value"].(string)) + DELIMITER)
				case RANDOM:
					// rowBuilder.WriteString(person.Gender + DELIMITER)
					genderFormat := strings.TrimSpace(jsonAttr["format"].(string))
					rand.NewSource(time.Now().UnixNano())
					if strings.EqualFold(genderFormat, "long") {
						longGenders := (descriptorJson["gender"].(map[string]interface{}))["range"].([]interface{})
						rowBuilder.WriteString(longGenders[rand.Intn(len(longGenders))].(string) + DELIMITER)
					} else {
						shortGenders := (descriptorJson["gender"].(map[string]interface{}))["short-range"].([]interface{})
						rowBuilder.WriteString(shortGenders[rand.Intn(len(shortGenders))].(string) + DELIMITER)
					}
				}
			case "boolean":
				switch dataGenType[jsonAttr["name"].(string)] {
				case DEFAULT:
					rowBuilder.WriteString(strings.TrimSpace(jsonAttr["default_value"].(string)) + DELIMITER)
				case RANDOM:
					boolFormat := strings.TrimSpace(jsonAttr["format"].(string))
					rand.NewSource(time.Now().UnixNano())
					if strings.EqualFold(boolFormat, "long") {
						longBool := (descriptorJson["boolean"].(map[string]interface{}))["range"].([]interface{})
						rowBuilder.WriteString(longBool[rand.Intn(len(longBool))].(string) + DELIMITER)
					} else {
						shortBool := (descriptorJson["boolean"].(map[string]interface{}))["short-range"].([]interface{})
						rowBuilder.WriteString(shortBool[rand.Intn(len(shortBool))].(string) + DELIMITER)
					}
				}
			case "ssn":
				switch dataGenType[jsonAttr["name"].(string)] {
				case DEFAULT:
					rowBuilder.WriteString(strings.TrimSpace(jsonAttr["default_value"].(string)) + DELIMITER)
				case NATURAL_SEQ:
					// ssnRegEx := `^(?!666|000|9\\d{2})\\d{3}-(?!00)\\d{2}-(?!0{4})\\d{4}$`
					ssnRegEx := (descriptorJson["ssn"].(map[string]interface{}))["format"].(string)
					ssnVal := gofakeit.Regex(ssnRegEx)
					mu.Lock()
					ssnList, found := (*attrLookups)[jsonAttr["name"].(string)]
					if found {
						for slices.Contains(ssnList, ssnVal) {
							ssnVal = gofakeit.Regex(ssnRegEx)
						}
					} else {
						ssnList = make([]string, numOfRows)
					}
					ssnList = append(ssnList, ssnVal)
					mu.Unlock()
					mu.Lock()
					(*attrLookups)[jsonAttr["name"].(string)] = ssnList
					mu.Unlock()
					rowBuilder.WriteString(ssnVal + DELIMITER)
				case RANDOM:
					ssnRegEx := (descriptorJson["ssn"].(map[string]interface{}))["format"].(string)
					rowBuilder.WriteString(gofakeit.Regex(ssnRegEx) + DELIMITER)
				}
			case "email":
				switch dataGenType[jsonAttr["name"].(string)] {
				case DEFAULT:
					rowBuilder.WriteString(strings.TrimSpace(jsonAttr["default_value"].(string)) + DELIMITER)
				case NATURAL_SEQ:
					emailId := person.Contact.Email
					mu.Lock()
					emailList, found := (*attrLookups)[jsonAttr["name"].(string)]
					if found {
						for slices.Contains(emailList, emailId) {
							emailId = gofakeit.Person().Contact.Email
						}
					} else {
						emailList = make([]string, numOfRows)
					}
					emailList = append(emailList, emailId)
					mu.Unlock()
					mu.Lock()
					(*attrLookups)[jsonAttr["name"].(string)] = emailList
					mu.Unlock()
					rowBuilder.WriteString(emailId + DELIMITER)
				case RANDOM:
					rowBuilder.WriteString(person.Contact.Email + DELIMITER)
				}
			case "phonenumber":
				switch dataGenType[jsonAttr["name"].(string)] {
				case DEFAULT:
					rowBuilder.WriteString(strings.TrimSpace(jsonAttr["default_value"].(string)) + DELIMITER)
				case NATURAL_SEQ:
					phoneNumber := person.Contact.Phone
					mu.Lock()
					phoneNumList, found := (*attrLookups)[jsonAttr["name"].(string)]
					if found {
						for slices.Contains(phoneNumList, phoneNumber) {
							phoneNumber = gofakeit.Person().Contact.Phone
						}
					} else {
						phoneNumList = make([]string, numOfRows)
					}
					phoneNumList = append(phoneNumList, phoneNumber)
					mu.Unlock()
					mu.Lock()
					(*attrLookups)[jsonAttr["name"].(string)] = phoneNumList
					mu.Unlock()
					rowBuilder.WriteString(phoneNumber + DELIMITER)
				case RANDOM:
					rowBuilder.WriteString(person.Contact.Phone + DELIMITER)
				}
			case "aadhar":
				switch dataGenType[jsonAttr["name"].(string)] {
				case DEFAULT:
					rowBuilder.WriteString(strings.TrimSpace(jsonAttr["default_value"].(string)) + DELIMITER)
				case NATURAL_SEQ:
					aadharRegEx := (descriptorJson["aadhar"].(map[string]interface{}))["format"].(string)
					aadharNum := gofakeit.Regex(aadharRegEx)
					mu.Lock()
					aadharList, found := (*attrLookups)[jsonAttr["name"].(string)]
					if found {
						for slices.Contains(aadharList, aadharNum) {
							aadharNum = gofakeit.Regex(aadharRegEx)
						}
					} else {
						aadharList = make([]string, numOfRows)
					}
					aadharList = append(aadharList, aadharNum)
					mu.Unlock()
					mu.Lock()
					(*attrLookups)[jsonAttr["name"].(string)] = aadharList
					mu.Unlock()
					rowBuilder.WriteString(aadharNum + DELIMITER)
				case RANDOM:
					aadharRegEx := (descriptorJson["aadhar"].(map[string]interface{}))["format"].(string)
					rowBuilder.WriteString(gofakeit.Regex(aadharRegEx) + DELIMITER)
				}
			case "creditcard":
				switch dataGenType[jsonAttr["name"].(string)] {
				case DEFAULT:
					rowBuilder.WriteString(strings.TrimSpace(jsonAttr["default_value"].(string)) + DELIMITER)
				case NATURAL_SEQ:
					creditCardNumber := getCreditCardNumber(jsonAttr)
					mu.Lock()
					creditCardNumList, found := (*attrLookups)[jsonAttr["name"].(string)]
					if found {
						for slices.Contains(creditCardNumList, creditCardNumber) {
							creditCardNumber = getCreditCardNumber(jsonAttr)
						}
					} else {
						creditCardNumList = make([]string, numOfRows)
					}
					creditCardNumList = append(creditCardNumList, creditCardNumber)
					mu.Unlock()
					mu.Lock()
					(*attrLookups)[jsonAttr["name"].(string)] = creditCardNumList
					mu.Unlock()
					rowBuilder.WriteString(creditCardNumber + DELIMITER)
				case RANDOM:
					rowBuilder.WriteString(getCreditCardNumber(jsonAttr) + DELIMITER)
				}
			case "zipcode":
				switch dataGenType[jsonAttr["name"].(string)] {
				case DEFAULT:
					rowBuilder.WriteString(strings.TrimSpace(jsonAttr["default_value"].(string)) + DELIMITER)
				case RANDOM:
					countryName := person.Address.Country
					fmt.Println(countryName)
					postalMap := getPostalCodeMap(countryName)
					if len(postalMap["Regex"]) > 0 {
						zip, _ := regexp.Compile(postalMap["Regex"])
						fmt.Println(zip.MatchString(gofakeit.Regex(postalMap["Regex"])))
						rowBuilder.WriteString(gofakeit.Regex(postalMap["Regex"]) + DELIMITER)
					} else {
						rowBuilder.WriteString(person.Address.Zip + DELIMITER)
					}
				}
			case "uuid":
				switch dataGenType[jsonAttr["name"].(string)] {
				case DEFAULT:
					rowBuilder.WriteString(strings.TrimSpace(jsonAttr["default_value"].(string)) + DELIMITER)
				case RANDOM:
					rowBuilder.WriteString(gofakeit.UUID() + DELIMITER)
				}
			case "ipaddress":
				switch dataGenType[jsonAttr["name"].(string)] {
				case DEFAULT:
					rowBuilder.WriteString(strings.TrimSpace(jsonAttr["default_value"].(string)) + DELIMITER)
				case RANDOM:
					iptypes := (descriptorJson["ipaddress"].(map[string]interface{}))["iptypes"].([]interface{})
					iptypesList := make([]string, len(iptypes))
					for i, v := range iptypes {
						iptypesList[i] = fmt.Sprintf(v.(string))
					}
					if strings.EqualFold(strings.TrimSpace(jsonAttr["ipaddress_type"].(string)), "ipv4") {
						rowBuilder.WriteString(gofakeit.IPv4Address() + DELIMITER)
					} else if strings.EqualFold(strings.TrimSpace(jsonAttr["ipaddress_type"].(string)), "ipv6") {
						rowBuilder.WriteString(gofakeit.IPv6Address() + DELIMITER)
					} else {
						rand.NewSource(time.Now().UnixNano())
						if rand.Intn(2) == 0 {
							rowBuilder.WriteString(gofakeit.IPv4Address() + DELIMITER)
						} else {
							rowBuilder.WriteString(gofakeit.IPv6Address() + DELIMITER)
						}
					}
				}
			case "timestamp":
				switch dataGenType[jsonAttr["name"].(string)] {
				case DEFAULT:
					rowBuilder.WriteString(strings.TrimSpace(jsonAttr["default_value"].(string)) + DELIMITER)
				case DUP_IN_RANGE:
					timestampRange := strings.Split(strings.TrimSpace(jsonAttr["range"].(string)), "~")
					timestampMin, _ := time.Parse(jsonAttr["timestamp_format"].(string), timestampRange[0])
					timestampMax, _ := time.Parse(jsonAttr["timestamp_format"].(string), timestampRange[1])
					rowBuilder.WriteString(gofakeit.DateRange(timestampMin, timestampMax).Format(jsonAttr["timestamp_format"].(string)) + DELIMITER)
				case RANDOM:
					rowBuilder.WriteString(gofakeit.Date().Format(jsonAttr["timestamp_format"].(string)) + DELIMITER)
				}
			}
		}
		dataRow := rowBuilder.String()
		dataRow = dataRow[:len(dataRow)-1]
		tdgChannel <- dataRow
		rowBuilder.Reset()
	}
	defer wg.Done()
}

func getPostalCodeMap(countryName string) map[string]string {
	for _, postalMap := range postalCodeJson {
		if strings.EqualFold(strings.TrimSpace(postalMap["Country"]), strings.TrimSpace(countryName)) {
			return postalMap
		}
	}
	return nil
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

func getBlankRowPercentage(blankPer int) int {
	// rowCount % (blankRowCount + 1) == 0
	blankRows := float64(numOfRows) * float64(blankPer) / 100
	if blankRows < 1 {
		return 0
	}
	return int(math.Round(float64(numOfRows) / blankRows))
}

func getCreditCardNumber(entry map[string]interface{}) string {
	cctype := strings.TrimSpace(entry["cctype"].(string))
	if strings.EqualFold(cctype, "any") {
		return gofakeit.CreditCardNumber(&gofakeit.CreditCardOptions{Types: []string{}, Gaps: true})
	} else {
		return gofakeit.CreditCardNumber(&gofakeit.CreditCardOptions{Types: []string{cctype}, Gaps: true})
	}
}
