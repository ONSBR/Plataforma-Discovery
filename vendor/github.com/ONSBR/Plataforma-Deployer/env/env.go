package env

import (
	"os"
)

//Get key from Environment Variables or return default value
func Get(key, _default string) string {
	value := os.Getenv(key)
	if value == "" {
		return _default
	}
	return value
}
