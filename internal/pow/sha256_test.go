package pow

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerateRandomString(t *testing.T) {
	for i := range make([]struct{}, 256) {
		length := i + 1

		str, err := GenerateRandomString(length)
		require.NoError(t, err)
		assert.Len(t, str, length)
	}
}

func TestGenerateRandomString_Error(t *testing.T) {
	_, err := GenerateRandomString(0)
	require.Error(t, err)

	_, err = GenerateRandomString(-1)
	require.Error(t, err)
}

func TestSolveAndVerify(t *testing.T) {
	ctx := context.Background()

	challenge := "test-challenge"
	difficulty := 2
	solutionLength := 8

	solution, err := SolveChallenge(ctx, challenge, difficulty, solutionLength)
	require.NoError(t, err)

	correct := VerifySolution(challenge, difficulty, solution)
	assert.True(t, correct)

	correct = VerifySolution(challenge, difficulty, "wrong-solution")
	assert.False(t, correct)

	wrongDifficulty := 3
	correct = VerifySolution(challenge, wrongDifficulty, solution)
	assert.False(t, correct)
}

func TestSolveChallenge_Timeout(t *testing.T) {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	challenge := "test-challenge"
	difficulty := 999
	solutionLength := 1

	_, err := SolveChallenge(ctx, challenge, difficulty, solutionLength)
	require.ErrorIs(t, err, context.DeadlineExceeded)
}

func TestVerifySolution(t *testing.T) {
	for _, c := range []struct {
		challenge  string
		difficulty int
		solution   string
	}{{
		challenge:  "25fd47ba9fc2a7da",
		difficulty: 2,
		solution:   "d554",
	}, {
		challenge:  "ea35fe5dd4f7f494",
		difficulty: 2,
		solution:   "71e3",
	}, {
		challenge:  "07ab6d42479d95e2",
		difficulty: 2,
		solution:   "cfc0",
	}, {
		challenge:  "f55aa53dbf683a28",
		difficulty: 2,
		solution:   "76b1949a",
	}, {
		challenge:  "c911614e25ee696f",
		difficulty: 2,
		solution:   "17afde65",
	}} {
		t.Run(c.challenge, func(t *testing.T) {
			correct := VerifySolution(c.challenge, c.difficulty, c.solution)
			assert.True(t, correct)
		})
	}
}
