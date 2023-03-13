package protocol

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestChallengeEncodeDecode(t *testing.T) {
	challenge := []byte("test-challenge")
	difficulty := 4
	solutionLength := 8
	message := "746573742d6368616c6c656e6765:4:8"

	messageActual := ChallengeEncode(challenge, difficulty, solutionLength)
	assert.Equal(t, message, messageActual)

	challengeActual, difficultyActual, solutionLengthActual, err := ChallengeDecode(messageActual)
	require.NoError(t, err)
	assert.Equal(t, challenge, challengeActual)
	assert.Equal(t, difficulty, difficultyActual)
	assert.Equal(t, solutionLength, solutionLengthActual)
}

func TestChallengeDecode_NewLine(t *testing.T) {
	challenge := []byte("test-challenge")
	difficulty := 4
	solutionLength := 8
	message := "746573742d6368616c6c656e6765:4:8\n"

	challengeActual, difficultyActual, solutionLengthActual, err := ChallengeDecode(message)
	require.NoError(t, err)
	assert.Equal(t, challenge, challengeActual)
	assert.Equal(t, difficulty, difficultyActual)
	assert.Equal(t, solutionLength, solutionLengthActual)
}

func TestChallengeDecode_Error(t *testing.T) {
	f := func(message string) {
		t.Run(message, func(t *testing.T) {
			_, _, _, err := ChallengeDecode(message)
			require.Error(t, err)
		})
	}

	f("")
	f(" ")
	f("\n")
	f("746573742d6368616c6c656e6765:4")
	f("746573742d6368616c6c656e6765:4 ")
	f("746573742d6368616c6c656e6765:4:8 ")
	f("746573742d6368616c6c656e6765")
	f("746573742d6368616c6c656e67654")
	f("4")
	f(":4")
	f("4:8")
	f(":4:8")
	f("746573742d6368616c6c656e6765:4:5:6")
	f("746573742d6368616c6c656e6765:foo:bar")
	f("746573742d6368616c6c656e6765:foo:8")
	f("746573742d6368616c6c656e6765:4:bar")
}
