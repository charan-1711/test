package utils

import "strings"

func IsLength(s string, n int) bool {
	return len(strings.TrimSpace(s)) >= n
}

func Contains(validValues []string, val string) bool {
	for _, el := range validValues {
		if el == val {
			return true
		}
	}

	return false
}
