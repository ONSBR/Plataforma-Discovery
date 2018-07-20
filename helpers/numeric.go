package helpers

import (
	"strconv"
)

//IsInt checks if a string s is a integer number
func IsInt(s string) bool {
	_, err := strconv.Atoi(s)
	return err == nil
}

func IsFloat(s string) bool {
	_, err := strconv.ParseFloat(s, 64)
	return err == nil
}

func IsNumeric(s string) bool {
	return IsFloat(s) || IsInt(s)
}
