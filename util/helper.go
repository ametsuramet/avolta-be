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

func Contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func GetDates(startDate, endDate time.Time) []time.Time {
	dates := []time.Time{}
	currentDate := startDate

	for currentDate.Before(endDate.AddDate(0, 0, 1)) {
		dates = append(dates, currentDate)
		currentDate = currentDate.AddDate(0, 0, 1)
	}
	return dates
}

func FormatDuration(d time.Duration) string {
	// Extract the total hours and minutes from the duration
	hours := int(d.Hours())
	minutes := int(d.Minutes()) % 60

	// Format as "HH:mm"
	return fmt.Sprintf("%02d:%02d", hours, minutes)
}
