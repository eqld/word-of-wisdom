package protocol

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFormatAndParse(t *testing.T) {
	challenge := "test-challenge"
	difficulty := 4
	message := "test-challenge:4"

	messageActual := FormatChallengeForClient(challenge, difficulty)
	assert.Equal(t, message, messageActual)

	challengeActual, difficultyActual, err := ParseChallengeForClient(messageActual)
	require.NoError(t, err)
	assert.Equal(t, challenge, challengeActual)
	assert.Equal(t, difficulty, difficultyActual)
}

func TestParseChallengeForClient_NewLine(t *testing.T) {
	challenge := "test-challenge"
	difficulty := 4
	message := "test-challenge:4\n"

	challengeActual, difficultyActual, err := ParseChallengeForClient(message)
	require.NoError(t, err)
	assert.Equal(t, challenge, challengeActual)
	assert.Equal(t, difficulty, difficultyActual)
}

func TestParseChallengeForClient_Error(t *testing.T) {
	f := func(message string) {
		t.Run(message, func(t *testing.T) {
			_, _, err := ParseChallengeForClient(message)
			require.Error(t, err)
		})
	}

	f("")
	f(" ")
	f("\n")
	f("test-challenge:4 ")
	f("test-challenge")
	f("test-challenge4")
	f("4")
	f("test-challenge:4:5")
	f("test-challenge:foo")
}
