package service

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func GetLocationName(latitude, longitude float64) (string, error) {
	url := fmt.Sprintf("https://nominatim.openstreetmap.org/reverse?lat=%f&lon=%f&format=json", latitude, longitude)
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return "", err
	}

	if displayName, ok := result["display_name"].(string); ok {
		return displayName, nil
	}

	return "", fmt.Errorf("location name not found")
}
