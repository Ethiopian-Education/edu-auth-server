package utils

import (
	"fmt"
	"time"
)


func FormatDatesToSimilarZone(date_1 string, date_2 string) (string, string, error) {
	t1, err := time.Parse(time.RFC3339Nano, date_1)
	if err != nil {
		fmt.Println("Error parsing dateTime1:", err)
		return "" , "", err
	}
	formatted_date_1 := t1.Format("2006-01-02 15:04:05 MST")

	t2, err := time.Parse(time.RFC3339Nano, date_2)
	if err != nil {
		fmt.Println("Error parsing dateTime1:", err)
		return "", "", err
	}
	formatted_date_2 := t2.Format("2006-01-02 15:04:05 MST")

	return formatted_date_1, formatted_date_2, nil
}