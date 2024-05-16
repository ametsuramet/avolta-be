package util

import (
	"avolta/config"
	"fmt"
	"time"
)

func GetURL(path string) string {
	return fmt.Sprintf("%s/%s", config.App.Server.BaseURL, path)
}

func GetCurrentTimestamp() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}
