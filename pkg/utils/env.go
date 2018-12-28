package utils

import (
	"os"
	"strings"
)

func GetEnvValue(keyName string) string {
	return os.Getenv(keyName)
}

func GetEnvBool(keyName string) bool {
	val, isSet := os.LookupEnv(keyName)
	return isSet && strings.ToUpper(val) == "TRUE"
}
