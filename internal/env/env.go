package env

import (
	"fmt"
	"os"
)

func Get(k string) (string, error) {
	if v, ok := os.LookupEnv(k); ok {
		return v, nil
	}
	return "", fmt.Errorf("environment variable (%s) not found", k)
}

func GetOrDefault(k string, dv string) string {
	v, ok := os.LookupEnv(k)
	if ok {
		return v
	}
	return dv
}
