package utils

import (
	"os"
	"strings"
)

func TryGetEnvValue(keyName string, fallback string) string {
	val, set := os.LookupEnv(keyName)
	if set {
		return val
	}
	return fallback
}

func GetEnvBool(keyName string) bool {
	val, isSet := os.LookupEnv(keyName)
	return isSet && strings.ToUpper(val) == "TRUE"
}
