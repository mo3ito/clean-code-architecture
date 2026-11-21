package utils

import (
	"os"
	"strconv"
)

func GetEnvInt(key string, fallback int) int {
	v := os.Getenv(key)

	if v == "" {
		return fallback
	}

	i, err := strconv.Atoi(v)
	if err != nil {
		return fallback
	}

	return i
}