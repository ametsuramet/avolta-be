package util

import (
	"avolta/config"
	"encoding/json"
	"fmt"
	"log"
	"time"
)

func GetURL(path string) string {
	return fmt.Sprintf("%s/%s", config.App.Server.BaseURL, path)
}

func GetCurrentTimestamp() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

func IntToLetters(number int32) (letters string) {
	number--
	if firstLetter := number / 26; firstLetter > 0 {
		letters += IntToLetters(firstLetter)
		letters += string('A' + number%26)
	} else {
		letters += string('A' + number)
	}

	return
}

func LogJson(data interface{}) {
	jsonString, _ := json.MarshalIndent(data, "", "  ")
	fmt.Println(string(jsonString))
}

func SaveLogJson(data interface{}) {
	jsonString, _ := json.MarshalIndent(data, "", "  ")
	log.Println(string(jsonString))
}
