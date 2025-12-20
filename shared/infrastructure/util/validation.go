package util

import (
	"strconv"
	"time"
)

func IsNumberOnly(s string) bool {
	numberOnly := true

	_, err := strconv.Atoi(s)
	if err != nil {
		numberOnly = false
	}

	return numberOnly
}

func IsContainsDateTime(s string) bool {
	_, err := time.Parse("2006-01-02 15:04:05", s)
	return err == nil
}
