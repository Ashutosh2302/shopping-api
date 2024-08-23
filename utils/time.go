package utils

import (
	"fmt"
	"time"
)

func GetEpochTime(timestamp string) (int64, error) {
	t, err := time.Parse(time.RFC3339, timestamp)
	if err != nil {
		fmt.Println("Error parsing ISO 8601 timestamp:", err)
		return 0, err
	}
	return t.Unix(), nil
}
