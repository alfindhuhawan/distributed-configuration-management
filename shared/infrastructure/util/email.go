package util

import (
	"regexp"
)

func ValidateFormatEmail(email string) bool {
	// REGEX PATTERN FOR EMAIL VALIDATION
	regexPattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`

	// COMPILE REGEX
	regex := regexp.MustCompile(regexPattern)

	// DO MATCH PATTERN USING REGEX EMAIL PATERN
	result := regex.MatchString(email)

	return result
}
