package validator

import (
    "regexp"
    "slices"
    "strings"
    "unicode/utf8"
)

var EmailRX = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

type Validator struct {
    NonFieldErrors []string
    FieldErrors    map[string]string
}

func (v *Validator) Valid() bool {
    return len(v.FieldErrors) == 0 && len(v.NonFieldErrors) == 0
}

func (v *Validator) AddNonFieldError(message string) {
    v.NonFieldErrors = append(v.NonFieldErrors, message)
}

func (v *Validator) AddFieldError(key, message string) {
    // must initialize map first if it has not already been initialized
    if v.FieldErrors == nil {
        v.FieldErrors = make(map[string]string)
    }

    if _, exists := v.FieldErrors[key]; !exists {
        v.FieldErrors[key] = message
    }
}

// CheckField() adds an error message to the FieldErrors map
func (v *Validator) CheckField(ok bool, key, message string) {
    if !ok {
        v.AddFieldError(key, message)
    }
}

// NotBlank() returns true if a value is not an empty string
func NotBlank(value string) bool {
    return strings.TrimSpace(value) != ""
}

// MaxChars() retunrs true if a value contains no more than n characters
func MaxChars(value string, n int) bool {
    return utf8.RuneCountInString(value) <= n
}

func MinChars(value string, n int) bool {
    return utf8.RuneCountInString(value) >= n
}

// PermittedValue() returns true if a value is in a list of specific
// permitted values
func PermittedValue[T comparable](value T, permittedValues ...T) bool {
    return slices.Contains(permittedValues, value)
}

func Matches(value string, rx *regexp.Regexp) bool {
    return rx.MatchString(value)
}
