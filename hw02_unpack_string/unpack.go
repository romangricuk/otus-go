package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(s string) (string, error) {
	sb := strings.Builder{}

	var multiplier rune
	var letter rune

	for _, char := range s {
		if !isValidChar(char) {
			return "", ErrInvalidString
		}

		if isLetter(char) {
			if letter > 0 {
				sb.WriteString(string(letter))
			}
			letter = char
		} else if isDigit(char) {
			if multiplier > 0 || letter == 0 {
				return "", ErrInvalidString
			}
			multiplier = char
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
	return char >= '0' && char <= '9'
}

func isLetter(char rune) bool {
	return char >= 'a' && char <= 'z'
}

func isValidChar(char rune) bool {
	return isDigit(char) || isLetter(char)
}
