package utils

import (
	"os"
	"strconv"
)

func GetEnv(key string, defaultValue string) string {
	val, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}
	return val
}

func GetEnvAsString(key string, defaultValue string) string {
	return GetEnv(key, defaultValue)
}

func GetEnvAsBool(key string, defaultValue bool) bool {
	valStr := GetEnv(key, "")
	if valStr == "" {
		return defaultValue
	}
	val, err := strconv.ParseBool(valStr)
	if err != nil {
		return defaultValue
	}
	return val
}

func GetEnvAsInt(key string, defaultValue int) int {
	valStr := GetEnv(key, "")
	if valStr == "" {
		return defaultValue
	}
	val, err := strconv.Atoi(valStr)
	if err != nil {
		return defaultValue
	}
	return val
}

func GetEnvAsBytes(key string, defaultValue []byte) []byte {
	valStr := GetEnv(key, "")
	if valStr == "" {
		return defaultValue
	}
	return []byte(valStr)
}
