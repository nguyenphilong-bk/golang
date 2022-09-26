package utils

import "regexp"

var IsStringAlphabetic = regexp.MustCompile(`^[a-zA-Z0-9]*$`).MatchString