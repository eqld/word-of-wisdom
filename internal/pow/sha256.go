package pow

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

// GenerateRandomString generates a random string of given length for the PoW challenge or solution.
// The message is returned as a string of hexadecimal digits.
func GenerateRandomString(length int) (string, error) {
	if length <= 0 {
		return "", fmt.Errorf("string length must be positive")
	}

	l := length
	if l%2 == 1 {
		l++
	}

	message := make([]byte, l/2)
	if _, err := rand.Read(message); err != nil {
		return "", fmt.Errorf("failed to generate random message: %w", err)
	}

	return hex.EncodeToString(message)[:length], nil
}

// SolveChallenge attempts to solve given challenge with given difficulty using a solution of given length.
func SolveChallenge(ctx context.Context, challenge string, difficulty, length int) (solution string, err error) {

	for ctx.Err() == nil && !VerifySolution(challenge, difficulty, solution) {
		if solution, err = GenerateRandomString(length); err != nil {
			return "", err
		}
	}

	return solution, ctx.Err()
}

// VerifySolution takes a PoW challenge, a difficulty level, and a solution as inputs.
// It appends the solution to the challenge, computes the SHA256 hash of the resulting string,
// and checks whether the first `difficulty` bytes of the hash are zero.
// If they are, the function returns `true`, indicating that the solution is valid.
func VerifySolution(challenge string, difficulty int, solution string) bool {

	challengeWithSolution := challenge + solution
	hash := sha256.Sum256([]byte(challengeWithSolution))

	for i := 0; i < difficulty; i++ {
		if hash[i] != 0 {
			return false
		}
	}

	return true
}
