package protocol

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFormatAndParse(t *testing.T) {
	challenge := "test-challenge"
	difficulty := 4
	solutionLength := 8
	message := "test-challenge:4:8"

	messageActual := FormatChallengeForClient(challenge, difficulty, solutionLength)
	assert.Equal(t, message, messageActual)

	challengeActual, difficultyActual, solutionLengthActual, err := ParseChallengeForClient(messageActual)
	require.NoError(t, err)
	assert.Equal(t, challenge, challengeActual)
	assert.Equal(t, difficulty, difficultyActual)
	assert.Equal(t, solutionLength, solutionLengthActual)
}

func TestParseChallengeForClient_NewLine(t *testing.T) {
	challenge := "test-challenge"
	difficulty := 4
	solutionLength := 8
	message := "test-challenge:4:8\n"

	challengeActual, difficultyActual, solutionLengthActual, err := ParseChallengeForClient(message)
	require.NoError(t, err)
	assert.Equal(t, challenge, challengeActual)
	assert.Equal(t, difficulty, difficultyActual)
	assert.Equal(t, solutionLength, solutionLengthActual)
}

func TestParseChallengeForClient_Error(t *testing.T) {
	f := func(message string) {
		t.Run(message, func(t *testing.T) {
			_, _, _, err := ParseChallengeForClient(message)
			require.Error(t, err)
		})
	}

	f("")
	f(" ")
	f("\n")
	f("test-challenge:4")
	f("test-challenge:4 ")
	f("test-challenge:4:8 ")
	f("test-challenge")
	f("test-challenge4")
	f("4")
	f(":4")
	f("4:8")
	f(":4:8")
	f("test-challenge:4:5:6")
	f("test-challenge:foo:bar")
	f("test-challenge:foo:8")
	f("test-challenge:4:bar")
}
