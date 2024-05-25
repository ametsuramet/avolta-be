package util

import (
	"avolta/config"
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
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

func SavedString(str *string) string {
	if str != nil {
		return *str
	}
	return ""
}
func FileURL(str string) string {
	if str != "" {
		return fmt.Sprintf("%s/%s", config.App.Server.BaseURL, str)
	}
	return ""
}
func SavedFloat(str *float64) float64 {
	if str != nil {
		return *str
	}
	return 0
}
func IntegerToRoman(number int) string {
	maxRomanNumber := 3999
	if number > maxRomanNumber {
		return strconv.Itoa(number)
	}

	conversions := []struct {
		value int
		digit string
	}{
		{1000, "M"},
		{900, "CM"},
		{500, "D"},
		{400, "CD"},
		{100, "C"},
		{90, "XC"},
		{50, "L"},
		{40, "XL"},
		{10, "X"},
		{9, "IX"},
		{5, "V"},
		{4, "IV"},
		{1, "I"},
	}

	var roman strings.Builder
	for _, conversion := range conversions {
		for number >= conversion.value {
			roman.WriteString(conversion.digit)
			number -= conversion.value
		}
	}

	return roman.String()
}

func ExtractPercentage(s string) float64 {
	re := regexp.MustCompile(`\d+(\.\d+)?`)

	match := re.FindString(s)
	if match == "" {
		return 0
	}

	number, err := strconv.ParseFloat(match, 64)
	if err != nil {
		return 0
	}

	return number
}

func ParseThousandSeparatedNumber(s string) float64 {
	s = strings.Replace(s, ",", "", -1)

	number, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0
	}

	return number
}
