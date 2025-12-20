package util

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
)

func StringContains(slices []string, comparizon string) bool {
	for _, a := range slices {
		if a == comparizon {
			return true
		}
	}

	return false
}

func ConvertByteToString(inBytes []byte) string {
	myString := string(inBytes[:])

	return myString
}

func ChangeEnterToSpace(s string) string {
	re := regexp.MustCompile(`\r?\n`)
	returnString := re.ReplaceAllString(s, " ")

	return returnString
}

func DebuggingStruct(input any) (string, error) {
	dataInByte, err := json.Marshal(input)
	if err != nil {
		return "", err
	}

	myString := string(dataInByte[:])

	return myString, nil
}

func DebuggingStructWithoutError(input any) string {
	var myString string
	dataInByte, err := json.Marshal(input)
	if err != nil {
		fmt.Println(err)
	} else {
		myString = string(dataInByte[:])
	}

	return myString
}

func ConvertEstimatedTime(estimatedTime int64) string {
	stringEstimatedTime := strconv.Itoa(int(estimatedTime))
	var result string
	if len(stringEstimatedTime) > 3 {
		result = fmt.Sprintf("%s%s%s", stringEstimatedTime[:2], ":", stringEstimatedTime[2:])
	} else {
		result = fmt.Sprintf("0%s%s%s", stringEstimatedTime[:1], ":", stringEstimatedTime[1:])
	}

	return result
}

func DebuggingResponseFromExternal(input []byte) (string, error) {
	var anyData interface{}
	err := json.Unmarshal(input, &anyData)
	if err != nil {
		fmt.Println(err.Error())
	}

	myString, err := DebuggingStruct(anyData)
	if err != nil {
		fmt.Println(err.Error())
	}

	return myString, nil
}
