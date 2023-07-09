package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

var dataGenType map[string]interface{}

func generateData() string {
	log.Print("Number of rows:", numOfRows)
	log.Print("Metadata file path:", metadataFileName)
	// reading the metadata file
	dat, err := os.ReadFile(metadataFileName)
	var metaDataJson []map[string]interface{}
	if err != nil {
		return err.Error()
	}
	isMetadataSchemaValid := json.Valid(dat)
	var errorMsg string
	if !isMetadataSchemaValid {
		errorMsg = "Invalid metadata json schema\n"
	}
	if err := json.Unmarshal(dat, &metaDataJson); err != nil {
		errorMsg += err.Error()
	}
	if len(errorMsg) > 0 {
		return errorMsg
	}
	// fmt.Println(metaDataJson)
	// fmt.Println(metaDataJson[1]["name"])
	schemaErrors := validateMetadataJsonSchema(metaDataJson)
	if schemaErrors != "" && len(schemaErrors) > 0 {
		return schemaErrors
	}
	fmt.Println(dataGenType)
	return "success"
}

func validateMetadataJsonSchema(metaDataJson []map[string]interface{}) string {
	// we can also use validator package from https://github.com/go-playground/validator
	var errorMessage strings.Builder
	dataGenType = make(map[string]interface{})
	for index, jsonAttr := range metaDataJson {
		fmt.Println("Index:", index)
		fmt.Println("Attr:", jsonAttr)
		switch jsonAttr["datatype"] {
		case "number":
			fmt.Println("Number")
			if (jsonAttr["default_value"] == nil) || (strings.TrimSpace(jsonAttr["default_value"].(string)) == "") {
				if (jsonAttr["duplicates_allowed"] == nil) ||
					(strings.TrimSpace(jsonAttr["duplicates_allowed"].(string)) == "") ||
					(!(strings.EqualFold(jsonAttr["duplicates_allowed"].(string), "no")) && !(strings.EqualFold(jsonAttr["duplicates_allowed"].(string), "yes"))) {
					errorMessage.WriteString(jsonAttr["name"].(string) + ": Invalid value for the property duplicates_allowed\n")
					continue
				}
				if (jsonAttr["range"] != nil) && (strings.TrimSpace(jsonAttr["range"].(string)) != "") {
					numRange := strings.Split(strings.TrimSpace(jsonAttr["range"].(string)), "~")
					if (numRange != nil) && ((len(numRange) != 2) || (numRange[0] == "" || numRange[1] == "")) {
						errorMessage.WriteString(jsonAttr["name"].(string) + ": Invalid value for the property range\n")
						continue
					}
					if (numRange != nil) && strings.EqualFold(jsonAttr["duplicates_allowed"].(string), "no") {
						numMin, errMin := strconv.Atoi(numRange[0])
						numMax, errMax := strconv.Atoi(numRange[1])
						if errMin != nil || errMax != nil {
							errorMessage.WriteString(jsonAttr["name"].(string) + ": range (min and max) should be integer value\n")
							continue
						}
						if numOfRows > (numMax - numMin) {
							errorMessage.WriteString(jsonAttr["name"].(string) + ": range should be greater than or equal to number of rows\n")
							continue
						}
						dataGenType[jsonAttr["name"].(string)] = SEQ_IN_RANGE
					} else {
						dataGenType[jsonAttr["name"].(string)] = DUP_IN_RANGE
					}
				} else {
					if (strings.TrimSpace(jsonAttr["duplicates_allowed"].(string)) != "") && strings.EqualFold(jsonAttr["duplicates_allowed"].(string), "yes") {
						dataGenType[jsonAttr["name"].(string)] = RANDOM
					} else {
						dataGenType[jsonAttr["name"].(string)] = NATURAL_SEQ
					}
				}
			} else {
				dataGenType[jsonAttr["name"].(string)] = DEFAULT
			}
		case "text":
			fmt.Println("Text")
		case "float":
			fmt.Println("Float")
		case "date":
			fmt.Println("Date")
		case "gender":
			fmt.Println("Gender")
		case "boolean":
			fmt.Println("Boolean")
		case "ssn":
			fmt.Println("SSN")
		case "creditcard":
			fmt.Println("CreditCard")
		case "email":
			fmt.Println("Email")
		case "phonenumber":
			fmt.Println("Phone Number")
		case "zipcode":
			fmt.Println("Zip code")
		case "uuid":
			fmt.Println("UUID")
		case "ipaddress":
			fmt.Println("IP Address")
		case "timestamp":
			fmt.Println("Timestamp")
		case "aadhar":
			fmt.Println("Aadhar")
		}
	}
	return errorMessage.String()
}
