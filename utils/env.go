package utils

import "os"

func GetEnvDefault(key, defVal string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defVal
	}

	return value
}
