package env

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

const testEnvName = "WOW_TEST"

func TestMustReadIntEnv(t *testing.T) {
	os.Setenv(testEnvName, "123")
	defer os.Setenv(testEnvName, "")

	exitFunc = func(int) { t.Fail() }

	v := MustReadIntEnv(testEnvName, 456, 5)
	assert.Equal(t, 123, v)
}

func TestMustReadIntEnv_DefaultValue(t *testing.T) {
	os.Setenv(testEnvName, "")

	exitFunc = func(int) { t.Fail() }

	v := MustReadIntEnv(testEnvName, 456, 5)
	assert.Equal(t, 456, v)
}

func TestMustReadIntEnv_WrongFormat(t *testing.T) {
	os.Setenv(testEnvName, "foo")
	defer os.Setenv(testEnvName, "")

	exitCode := 5

	exitFuncCalled := false
	exitFunc = func(code int) {
		assert.Equal(t, exitCode, code)
		exitFuncCalled = true
	}

	MustReadIntEnv(testEnvName, 456, exitCode)
	assert.True(t, exitFuncCalled)
}
