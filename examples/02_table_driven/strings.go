package tabledriven

import (
	"strings"
	"unicode/utf8"
)

func ToUpper(s string) string {
	return strings.ToUpper(s)
}

func Reverse(s string) string {
	r := []rune(s)
	for i := 0; i < len(r)/2; i++ {
		r[i], r[len(r)-1-i] = r[len(r)-1-i], r[i]
	}
	return string(r)
}

func IsPalindrome(s string) bool {
	s = strings.ToLower(s)
	s = strings.ReplaceAll(s, " ", "")

	runes := []rune(s)
	for i := 0; i < len(runes)/2; i++ {
		if runes[i] != runes[len(runes)-1-i] {
			return false
		}
	}
	return true
}

func CountVowels(s string) int {
	count := 0
	vowels := "aeiouAEIOU"
	for _, r := range s {
		if strings.ContainsRune(vowels, r) {
			count++
		}
	}
	return count
}

func IsValidEmail(email string) bool {
	if !strings.Contains(email, "@") {
		return false
	}
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return false
	}
	if len(parts[0]) == 0 || len(parts[1]) == 0 {
		return false
	}
	if !strings.Contains(parts[1], ".") {
		return false
	}
	return true
}

func TruncateString(s string, maxLen int) string {
	if maxLen <= 0 {
		return ""
	}

	if utf8.RuneCountInString(s) <= maxLen {
		return s
	}

	runes := []rune(s)
	if maxLen <= 3 {
		return string(runes[:maxLen])
	}

	return string(runes[:maxLen-3]) + "..."
}
