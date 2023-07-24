package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"golang.org/x/exp/slices"
)

var dataGenType map[string]string
var metaDataJson []map[string]interface{}
var descriptorJson map[string]interface{}

func validateMetadata() string {
	log.Print("Number of rows:", numOfRows)
	log.Print("Metadata file path:", metadataFileName)
	// reading the metadata file
	metadataFileReader, err := os.ReadFile(metadataFileName)
	if err != nil {
		return err.Error()
	}
	isMetadataSchemaValid := json.Valid(metadataFileReader)
	var errorMsg string
	if !isMetadataSchemaValid {
		errorMsg = "Invalid metadata json schema\n"
	}
	if err := json.Unmarshal(metadataFileReader, &metaDataJson); err != nil {
		errorMsg += err.Error()
	}
	if len(errorMsg) > 0 {
		return errorMsg
	}
	// fmt.Println(metaDataJson)
	// fmt.Println(metaDataJson[1]["name"])

	descriptorFilePath, err := filepath.Abs("./resources/descriptor.json")
	if err != nil {
		return err.Error()
	}
	desciptorJsonReader, err := os.ReadFile(descriptorFilePath)
	if err != nil {
		return err.Error()
	}
	isDesciptorJsonValid := json.Valid(desciptorJsonReader)
	if !isDesciptorJsonValid {
		errorMsg = "Invalid descriptor json file\n"
	}
	if err := json.Unmarshal(desciptorJsonReader, &descriptorJson); err != nil {
		errorMsg += err.Error()
	}
	if len(errorMsg) > 0 {
		return errorMsg
	}

	schemaErrors := validateMetadataJsonSchema(metaDataJson, descriptorJson)
	if schemaErrors != "" && len(schemaErrors) > 0 {
		return schemaErrors
	}
	log.Print(dataGenType)
	return "success"
}

func validateMetadataJsonSchema(metaDataJson []map[string]interface{}, descriptorJson map[string]interface{}) string {
	// we can also use validator package from https://github.com/go-playground/validator
	var errorMessage strings.Builder
	dataGenType = make(map[string]string)
	for _, jsonAttr := range metaDataJson {
		switch jsonAttr["datatype"] {
		case "number":
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
				if _, err := strconv.Atoi(strings.TrimSpace(jsonAttr["default_value"].(string))); err != nil {
					errorMessage.WriteString(jsonAttr["name"].(string) + ": default value should be of type number\n")
					continue
				}
				dataGenType[jsonAttr["name"].(string)] = DEFAULT
			}
		case "text":
			if (jsonAttr["default_value"] == nil) || (strings.TrimSpace(jsonAttr["default_value"].(string)) == "") {
				if (jsonAttr["duplicates_allowed"] == nil) ||
					(strings.TrimSpace(jsonAttr["duplicates_allowed"].(string)) == "") ||
					(!(strings.EqualFold(jsonAttr["duplicates_allowed"].(string), "no")) && !(strings.EqualFold(jsonAttr["duplicates_allowed"].(string), "yes"))) {
					errorMessage.WriteString(jsonAttr["name"].(string) + ": Invalid value for the property duplicates_allowed\n")
					continue
				}
				if jsonAttr["range"] != nil {
					textRange := jsonAttr["range"].([]interface{})
					if (textRange != nil) && strings.EqualFold(jsonAttr["duplicates_allowed"].(string), "no") {
						if len(textRange) == 0 {
							errorMessage.WriteString(jsonAttr["name"].(string) + ": Invalid value for the property range\n")
							continue
						}
						if numOfRows > len(textRange) {
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
		case "float":
			if (jsonAttr["default_value"] == nil) || (strings.TrimSpace(jsonAttr["default_value"].(string)) == "") {
				if (jsonAttr["duplicates_allowed"] == nil) ||
					(strings.TrimSpace(jsonAttr["duplicates_allowed"].(string)) == "") ||
					(!(strings.EqualFold(jsonAttr["duplicates_allowed"].(string), "no")) && !(strings.EqualFold(jsonAttr["duplicates_allowed"].(string), "yes"))) {
					errorMessage.WriteString(jsonAttr["name"].(string) + ": Invalid value for the property duplicates_allowed\n")
					continue
				}
				if strings.TrimSpace(jsonAttr["scale"].(string)) != "" {
					if _, err := strconv.Atoi(strings.TrimSpace(jsonAttr["scale"].(string))); err != nil {
						errorMessage.WriteString(jsonAttr["name"].(string) + ": Invalid value for the property scale\n")
						continue
					}
				}
				if (jsonAttr["range"] != nil) && (strings.TrimSpace(jsonAttr["range"].(string)) != "") {
					floatRange := strings.Split(strings.TrimSpace(jsonAttr["range"].(string)), "~")
					if (floatRange != nil) && ((len(floatRange) != 2) || (floatRange[0] == "" || floatRange[1] == "")) {
						errorMessage.WriteString(jsonAttr["name"].(string) + ": Invalid value for the property range\n")
						continue
					}
					if (floatRange != nil) && strings.EqualFold(jsonAttr["duplicates_allowed"].(string), "no") {
						floatMin, errMin := strconv.ParseFloat(floatRange[0], 64)
						floatMax, errMax := strconv.ParseFloat(floatRange[1], 64)
						if errMin != nil || errMax != nil {
							errorMessage.WriteString(jsonAttr["name"].(string) + ": range (min and max) should be floating (decimal) value\n")
							continue
						}
						if numOfRows > (int(floatMax) - int(floatMin)) {
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
				if _, err := strconv.ParseFloat(strings.TrimSpace(jsonAttr["default_value"].(string)), 64); err != nil {
					errorMessage.WriteString(jsonAttr["name"].(string) + ": default value should be of type float\n")
					continue
				}
				dataGenType[jsonAttr["name"].(string)] = DEFAULT
			}
		case "date":
			if (jsonAttr["default_value"] == nil) || (strings.TrimSpace(jsonAttr["default_value"].(string)) == "") {
				if (jsonAttr["range"] != nil) && (strings.TrimSpace(jsonAttr["range"].(string)) != "") {
					dateRange := strings.Split(strings.TrimSpace(jsonAttr["range"].(string)), "~")
					if dateRange != nil {
						if (len(dateRange) != 2) || (dateRange[0] == "" || dateRange[1] == "") {
							errorMessage.WriteString(jsonAttr["name"].(string) + ": Invalid value for the property range\n")
							continue
						}
						if (jsonAttr["date_format"] != nil) && (strings.TrimSpace(jsonAttr["date_format"].(string)) != "") {
							// layout param in parse function should always a reference date e.g., 01/02/2006 03:04:05PM
							dateMin, errMin := time.Parse(jsonAttr["date_format"].(string), dateRange[0])
							dateMax, errMax := time.Parse(jsonAttr["date_format"].(string), dateRange[1])
							if errMin != nil {
								errorMessage.WriteString(jsonAttr["name"].(string) + ": invalid date_format or date value\n")
								errorMessage.WriteString(errMin.Error())
								continue
							}
							if errMax != nil {
								errorMessage.WriteString(jsonAttr["name"].(string) + ": invalid date_format or date value\n")
								errorMessage.WriteString(errMax.Error())
								continue
							}
							if dateMax.Compare(dateMin) <= 0 {
								errorMessage.WriteString(jsonAttr["name"].(string) + ": min date should always less than are equal to max date\n")
								continue
							}
						}
					}
					dataGenType[jsonAttr["name"].(string)] = DUP_IN_RANGE
				} else {
					dataGenType[jsonAttr["name"].(string)] = RANDOM
				}
			} else {
				if _, err := time.Parse(jsonAttr["date_format"].(string), strings.TrimSpace(jsonAttr["default_value"].(string))); err != nil {
					errorMessage.WriteString(jsonAttr["name"].(string) + ": default value should be in the choosen date format " + jsonAttr["date_format"].(string) + "\n")
					continue
				}
				dataGenType[jsonAttr["name"].(string)] = DEFAULT
			}
		case "gender":
			if (jsonAttr["default_value"] == nil) || (strings.TrimSpace(jsonAttr["default_value"].(string)) == "") {
				if (jsonAttr["format"] != nil) && (strings.TrimSpace(jsonAttr["format"].(string)) != "") {
					genderFormat := strings.TrimSpace(jsonAttr["format"].(string))
					if (genderFormat != "") && (!strings.EqualFold(genderFormat, "long") && !strings.EqualFold(genderFormat, "short")) {
						errorMessage.WriteString(jsonAttr["name"].(string) + ": Invalid value for the property format\n")
						continue
					}
				}
				dataGenType[jsonAttr["name"].(string)] = RANDOM
			} else {
				genderDefaultVal := strings.TrimSpace(jsonAttr["default_value"].(string))
				if (jsonAttr["format"] != nil) && (strings.TrimSpace(jsonAttr["format"].(string)) != "") {
					genderFormat := strings.TrimSpace(jsonAttr["format"].(string))
					if genderFormat != "" && strings.EqualFold(genderFormat, "long") &&
						!strings.EqualFold(genderDefaultVal, "male") && !strings.EqualFold(genderDefaultVal, "female") && !strings.EqualFold(genderDefaultVal, "others") {
						errorMessage.WriteString(jsonAttr["name"].(string) + ": default value should be in the list ['male', 'female', 'others']\n")
						continue
					}
					if genderFormat != "" && strings.EqualFold(genderFormat, "short") &&
						!strings.EqualFold(genderDefaultVal, "m") && !strings.EqualFold(genderDefaultVal, "f") && !strings.EqualFold(genderDefaultVal, "o") {
						errorMessage.WriteString(jsonAttr["name"].(string) + ": default value should be in the list ['m', 'f', 'o']\n")
						continue
					}
				} else {
					if !strings.EqualFold(genderDefaultVal, "male") && !strings.EqualFold(genderDefaultVal, "female") &&
						!strings.EqualFold(genderDefaultVal, "others") && !strings.EqualFold(genderDefaultVal, "m") &&
						!strings.EqualFold(genderDefaultVal, "f") && !strings.EqualFold(genderDefaultVal, "o") {
						errorMessage.WriteString(jsonAttr["name"].(string) + ": default value should be in the list ['male', 'female', 'others', 'm', 'f', 'o']\n")
						continue
					}
				}
				dataGenType[jsonAttr["name"].(string)] = DEFAULT
			}
		case "boolean":
			if (jsonAttr["default_value"] == nil) || (strings.TrimSpace(jsonAttr["default_value"].(string)) == "") {
				if (jsonAttr["format"] != nil) && (strings.TrimSpace(jsonAttr["format"].(string)) != "") {
					boolFormat := strings.TrimSpace(jsonAttr["format"].(string))
					if (boolFormat != "") && (!strings.EqualFold(boolFormat, "long") && !strings.EqualFold(boolFormat, "short")) {
						errorMessage.WriteString(jsonAttr["name"].(string) + ": Invalid value for the property format\n")
						continue
					}
				}
				dataGenType[jsonAttr["name"].(string)] = RANDOM
			} else {
				boolDefaultVal := strings.TrimSpace(jsonAttr["default_value"].(string))
				if (jsonAttr["format"] != nil) && (strings.TrimSpace(jsonAttr["format"].(string)) != "") {
					boolFormat := strings.TrimSpace(jsonAttr["format"].(string))
					if boolFormat != "" && strings.EqualFold(boolFormat, "long") &&
						!strings.EqualFold(boolDefaultVal, "true") && !strings.EqualFold(boolDefaultVal, "false") {
						errorMessage.WriteString(jsonAttr["name"].(string) + ": default value should be in the list ['true', 'false']\n")
						continue
					}
					if boolFormat != "" && strings.EqualFold(boolFormat, "short") &&
						!strings.EqualFold(boolDefaultVal, "t") && !strings.EqualFold(boolDefaultVal, "f") {
						errorMessage.WriteString(jsonAttr["name"].(string) + ": default value should be in the list ['t', 'f']\n")
						continue
					}
				} else {
					if !strings.EqualFold(boolDefaultVal, "true") && !strings.EqualFold(boolDefaultVal, "false") &&
						!strings.EqualFold(boolDefaultVal, "t") && !strings.EqualFold(boolDefaultVal, "f") {
						errorMessage.WriteString(jsonAttr["name"].(string) + ": default value should be in the list ['true', 'false', 't', 'f']\n")
						continue
					}
				}
				dataGenType[jsonAttr["name"].(string)] = DEFAULT
			}
		case "ssn", "email", "phonenumber", "aadhar":
			if (jsonAttr["default_value"] == nil) || (strings.TrimSpace(jsonAttr["default_value"].(string)) == "") {
				if (jsonAttr["duplicates_allowed"] == nil) ||
					(strings.TrimSpace(jsonAttr["duplicates_allowed"].(string)) == "") ||
					(!(strings.EqualFold(jsonAttr["duplicates_allowed"].(string), "no")) && !(strings.EqualFold(jsonAttr["duplicates_allowed"].(string), "yes"))) {
					errorMessage.WriteString(jsonAttr["name"].(string) + ": Invalid value for the property duplicates_allowed\n")
					continue
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
		case "creditcard":
			if (jsonAttr["default_value"] == nil) || (strings.TrimSpace(jsonAttr["default_value"].(string)) == "") {
				if (jsonAttr["duplicates_allowed"] == nil) ||
					(strings.TrimSpace(jsonAttr["duplicates_allowed"].(string)) == "") ||
					(!(strings.EqualFold(jsonAttr["duplicates_allowed"].(string), "no")) && !(strings.EqualFold(jsonAttr["duplicates_allowed"].(string), "yes"))) {
					errorMessage.WriteString(jsonAttr["name"].(string) + ": Invalid value for the property duplicates_allowed\n")
					continue
				} else {
					if (strings.TrimSpace(jsonAttr["duplicates_allowed"].(string)) != "") && strings.EqualFold(jsonAttr["duplicates_allowed"].(string), "yes") {
						dataGenType[jsonAttr["name"].(string)] = RANDOM
					} else {
						dataGenType[jsonAttr["name"].(string)] = NATURAL_SEQ
					}
				}
				if (jsonAttr["cctype"] == nil) || (strings.TrimSpace(jsonAttr["cctype"].(string)) == "") {
					errorMessage.WriteString(jsonAttr["name"].(string) + ": Invalid value for the property cctype\n")
					continue
				} else {
					cctypes := (descriptorJson["creditcard"].(map[string]interface{}))["cctypes"].([]interface{})
					cctypesList := make([]string, len(cctypes))
					for i, v := range cctypes {
						cctypesList[i] = fmt.Sprintf(v.(string))
					}
					if !strings.EqualFold(jsonAttr["cctype"].(string), "any") && !slices.Contains(cctypesList, jsonAttr["cctype"].(string)) {
						errorMessage.WriteString(jsonAttr["name"].(string) + ": cctype value should be in the list " + strings.Join(cctypesList, ", ") + "\n")
						continue
					}
				}
			} else {
				dataGenType[jsonAttr["name"].(string)] = DEFAULT
			}
		case "zipcode", "uuid":
			if (jsonAttr["default_value"] == nil) || (strings.TrimSpace(jsonAttr["default_value"].(string)) == "") {
				dataGenType[jsonAttr["name"].(string)] = RANDOM
			} else {
				dataGenType[jsonAttr["name"].(string)] = DEFAULT
			}
		case "ipaddress":
			if (jsonAttr["default_value"] == nil) || (strings.TrimSpace(jsonAttr["default_value"].(string)) == "") {
				if (jsonAttr["ipaddress_type"] == nil) || (strings.TrimSpace(jsonAttr["ipaddress_type"].(string)) == "") {
					errorMessage.WriteString(jsonAttr["name"].(string) + ": Invalid value for the property ipaddress_type\n")
					continue
				} else {
					iptypes := (descriptorJson["ipaddress"].(map[string]interface{}))["iptypes"].([]interface{})
					iptypesList := make([]string, len(iptypes))
					for i, v := range iptypes {
						iptypesList[i] = fmt.Sprintf(v.(string))
					}
					if !strings.EqualFold(jsonAttr["ipaddress_type"].(string), "any") && !slices.Contains(iptypesList, jsonAttr["ipaddress_type"].(string)) {
						errorMessage.WriteString(jsonAttr["name"].(string) + ": ipaddress_type value should be in the list " + strings.Join(iptypesList, ", ") + "\n")
						continue
					}
				}
				dataGenType[jsonAttr["name"].(string)] = RANDOM
			} else {
				dataGenType[jsonAttr["name"].(string)] = DEFAULT
			}
		case "timestamp":
			if (jsonAttr["default_value"] == nil) || (strings.TrimSpace(jsonAttr["default_value"].(string)) == "") {
				if (jsonAttr["range"] != nil) && (strings.TrimSpace(jsonAttr["range"].(string)) != "") {
					timestampRange := strings.Split(strings.TrimSpace(jsonAttr["range"].(string)), "~")
					if timestampRange != nil {
						if (len(timestampRange) != 2) || (timestampRange[0] == "" || timestampRange[1] == "") {
							errorMessage.WriteString(jsonAttr["name"].(string) + ": Invalid value for the property range\n")
							continue
						}
						if (jsonAttr["timestamp_format"] != nil) && (strings.TrimSpace(jsonAttr["timestamp_format"].(string)) != "") {
							// layout param in parse function should always a reference date e.g., 01/02/2006 03:04:05PM
							timestampMin, errMin := time.Parse(jsonAttr["timestamp_format"].(string), timestampRange[0])
							timestampMax, errMax := time.Parse(jsonAttr["timestamp_format"].(string), timestampRange[1])
							if errMin != nil {
								errorMessage.WriteString(jsonAttr["name"].(string) + ": invalid timestamp_format or timestamp value\n")
								errorMessage.WriteString(errMin.Error())
								continue
							}
							if errMax != nil {
								errorMessage.WriteString(jsonAttr["name"].(string) + ": invalid timestamp_format or timestamp value\n")
								errorMessage.WriteString(errMax.Error())
								continue
							}
							if timestampMax.Compare(timestampMin) <= 0 {
								errorMessage.WriteString(jsonAttr["name"].(string) + ": min timestamp should always less than are equal to max timestamp\n")
								continue
							}
						}
					}
					dataGenType[jsonAttr["name"].(string)] = DUP_IN_RANGE
				} else {
					dataGenType[jsonAttr["name"].(string)] = RANDOM
				}
			} else {
				if _, err := time.Parse(jsonAttr["timestamp_format"].(string), strings.TrimSpace(jsonAttr["default_value"].(string))); err != nil {
					errorMessage.WriteString(jsonAttr["name"].(string) + ": default value should be in the choosen timestamp format " + jsonAttr["timestamp_format"].(string) + "\n")
					continue
				}
				dataGenType[jsonAttr["name"].(string)] = DEFAULT
			}
		}
	}
	return errorMessage.String()
}
