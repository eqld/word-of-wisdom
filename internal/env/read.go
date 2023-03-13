package env

import (
	"log"
	"os"
	"strconv"
)

var exitFunc = os.Exit

// MustReadIntEnv reads integer value from given environment variable
// or terminates the program with special exit code if the value is not a valid integer.
func MustReadIntEnv(envName string, defaultValue int, exitCode int) int {
	str := os.Getenv(envName)
	if str == "" {
		return defaultValue
	}

	v, err := strconv.Atoi(str)
	if err != nil {
		log.Printf("wrong value of '%v' environment variable: %v\n", envName, err)
		exitFunc(exitCode)
	}

	return v
}
