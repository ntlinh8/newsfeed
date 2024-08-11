package validator

import (
	"fmt"
	"regexp"
)

func ValidateNotEmpty(s string, fieldName string) error {
	if len(s) == 0 {
		return fmt.Errorf("%s must be not empty", fieldName)
	}
	return nil
}

func ValidateLength(s string, min, max int, fieldName string) error {
	if len(s) < min {
		return fmt.Errorf("%s length must be gte %d", fieldName, min)
	} else if len(s) > max {
		return fmt.Errorf("%s length must be lte %d", fieldName, max)
	}
	return nil
}

func ValidateAsciiCharacters(s string, fieldName string) error {
	match, err := regexp.MatchString("^[a-zA-Z0-9]$", s)
	if err != nil {
		return fmt.Errorf("failed to match regex: %s", err)
	}
	if !match {
		return fmt.Errorf("%s must only has [a-zA-Z0-9]")
	}
	return nil
}
