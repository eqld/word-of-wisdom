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
	f("wrong-format:4:8")
}

func TestSolutionEncodeDecode(t *testing.T) {
	solution := []byte("test-solution")
	message := "746573742d736f6c7574696f6e"

	messageActual := SolutionEncode(solution)
	assert.Equal(t, message, messageActual)

	solutionActual, err := SolutionDecode(messageActual)
	require.NoError(t, err)
	assert.Equal(t, solution, solutionActual)
}

func TestSolutionDecode_Error(t *testing.T) {
	_, err := SolutionDecode("")
	require.Error(t, err)

	_, err = SolutionDecode("wrong-message")
	require.Error(t, err)
}

func TestQuoteEncodeDecode(t *testing.T) {
	quote := "Foo is Bar.\n        -- anonymous"
	message := "Rm9vIGlzIEJhci4KICAgICAgICAtLSBhbm9ueW1vdXM="

	messageActual := QuoteEncode(quote)
	assert.Equal(t, message, messageActual)

	quoteActual, err := QuoteDecode(messageActual)
	require.NoError(t, err)
	assert.Equal(t, quote, quoteActual)
}

func TestQuoteDecode_NewLine(t *testing.T) {
	quote, err := QuoteDecode("Rm9vIGlzIEJhci4KICAgICAgICAtLSBhbm9ueW1vdXM=\n")
	require.NoError(t, err)
	assert.Equal(t, "Foo is Bar.\n        -- anonymous", quote)
}

func TestQuoteDecode_Error(t *testing.T) {
	_, err := QuoteDecode("")
	require.Error(t, err)

	_, err = QuoteDecode("\n")
	require.Error(t, err)

	_, err = QuoteDecode("wrong-message")
	require.Error(t, err)
}
