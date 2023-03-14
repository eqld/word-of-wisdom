package pow

import (
	"context"
	"encoding/hex"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerateRandomBytes(t *testing.T) {
	for i := range make([]struct{}, 256) {
		length := i + 1

		v, err := GenerateRandomBytes(length)
		require.NoError(t, err)
		assert.Len(t, v, length)
	}
}

func TestGenerateRandomBytes_Error(t *testing.T) {
	_, err := GenerateRandomBytes(0)
	require.Error(t, err)

	_, err = GenerateRandomBytes(-1)
	require.Error(t, err)
}

func TestSolveAndVerify(t *testing.T) {
	ctx := context.Background()

	challenge := []byte("test-challenge")
	difficulty := 2
	solutionLength := 8

	solution, err := SolveChallenge(ctx, challenge, difficulty, solutionLength)
	require.NoError(t, err)

	correct := VerifySolution(challenge, solution, difficulty)
	assert.True(t, correct)

	correct = VerifySolution(challenge, []byte("wrong-solution"), difficulty)
	assert.False(t, correct)

	wrongDifficulty := 999
	correct = VerifySolution(challenge, solution, wrongDifficulty)
	assert.False(t, correct)
}

func TestVerifySolution_Empty(t *testing.T) {
	challengeWithEmptySolution, err := hex.DecodeString("f3a8b76fd9afdebe9871d2962893c8bde78266a055b3960846ac78d3304a")
	require.NoError(t, err)

	valid := VerifySolution(challengeWithEmptySolution, []byte{}, 2)
	assert.False(t, valid)

	valid = VerifySolution(challengeWithEmptySolution, nil, 2)
	assert.False(t, valid)
}

func TestSolveChallenge_Timeout(t *testing.T) {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	challenge := []byte("test-challenge")
	difficulty := 999
	solutionLength := 1

	_, err := SolveChallenge(ctx, challenge, difficulty, solutionLength)
	require.ErrorIs(t, err, context.DeadlineExceeded)
}
