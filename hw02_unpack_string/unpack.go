package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(s string) (string, error) {
	sb := strings.Builder{}

	var multiplier rune
	var letter rune

	for _, char := range s {
		if isDigit(char) {
			if multiplier > 0 || letter == 0 {
				return "", ErrInvalidString
			}
			multiplier = char
		} else {
			if letter > 0 {
				sb.WriteString(string(letter))
			}
			letter = char
		}

		if multiplier > 0 && letter > 0 {
			multiplierInt, _ := strconv.Atoi(string(multiplier))
			sb.WriteString(strings.Repeat(string(letter), multiplierInt))

			letter = 0
			multiplier = 0
		}
	}

	if letter > 0 {
		sb.WriteString(string(letter))
	}

	return sb.String(), nil
}

func isDigit(char rune) bool {
	return unicode.IsDigit(char)
}
