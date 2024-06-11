package internal

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func RedisConnStringSplitter(fullConnString string) (string, string) {
	url := strings.Split(fullConnString, "@")[1]
	password := strings.Split(fullConnString, "@")[0]
	password = strings.Split(password, "redis://:")[1]
	return url, password
}

func GeneratePostgresConnString(url string, username string, password string, dbName string, port string) string {
	return fmt.Sprintf(
		"host=%v user=%v password=%v dbname=%v port=%v sslmode=disable",
		url,
		username,
		password,
		dbName,
		port,
	)
}

// GetEnvOrDefault returns the value of the environment variable named by the key.
// If the environment variable is empty, the defaultValue is returned.
// If the environment variable is not empty, the value is coerced to the expectedType.
// If the value cannot be coerced to the expectedType, panic.
func GetEnvOrDefault(key string, defaultValue string, expectedType string) interface{} {
	Value := os.Getenv(key)
	if Value == "" {
		Value = defaultValue
	}
	switch expectedType {
	case "bool":
		rVal, ok := coerceBool(Value)
		if !ok {
			panic("invalid bool value for " + key)
		}
		return rVal
	case "int":
		rVal, ok := coerceInt(Value)
		if !ok {
			panic("invalid int value for " + key)
		}
		return rVal
	case "string":
		return Value
	case "[]string":
		return coerceSliceOfStrings(Value)
	default:
		panic("invalid expectedType for " + key)
	}
}

// coerceBool converts a string to a bool.
// If the string is not a valid bool, the second return value is false.
func coerceBool(value string) (bool, bool) {
	if value == "true" {
		return true, true
	}
	if value == "false" {
		return false, true
	}
	return false, false
}

// coerceInt converts a string to an int.
// If the string is not a valid int, the second return value is false.
func coerceInt(value string) (int, bool) {
	convertedInt, err := strconv.Atoi(value)
	if err != nil {
		return 0, false
	}
	return convertedInt, true
}

func coerceSliceOfStrings(value string) []string {
	convertedSlice := []string{}
	if value == "" {
		return convertedSlice
	}
	return strings.Split(value, ",")
}
