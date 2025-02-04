package validator

import (
	"regexp"
	"slices"
	"strings"
	"unicode/utf8"
)

var EmailRX = regexp.MustCompile(`^[a-zA-Z0-9.!#$%&'*+/=?^_` + "`" + `{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$`)

type Validator struct {
	FieldErrors map[string]string
}

func (val *Validator) Valid() bool {
	return len(val.FieldErrors) == 0
}

func (val *Validator) AddfieldError(key, message string) {
	if val.FieldErrors == nil {
		val.FieldErrors = make(map[string]string)
	}
	if _, exists := val.FieldErrors[key]; !exists {
		val.FieldErrors[key] = message
	}
}

func (val *Validator) CheckField(ok bool, key, message string) {
	if !ok {
		val.AddfieldError(key, message)
	}
}

func NotBlank(value string) bool {
	return strings.TrimSpace(value) != ""
}

func MaxChars(value string, n int) bool {
	return utf8.RuneCountInString(value) <= n
}

func PermittedValue[T comparable](value T, permittedValues ...T) bool {
	return slices.Contains(permittedValues, value)
}

func MinChars(value string, n int) bool {
	return utf8.RuneCountInString(value) >= n
}

func Matches(value string, rx *regexp.Regexp) bool {
	return rx.MatchString(value)
}
