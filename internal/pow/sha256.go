package pow

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
)

// GenerateRandomBytes generates a random sequence of bytes of given length for the PoW challenge or solution.
func GenerateRandomBytes(length int) ([]byte, error) {
	if length <= 0 {
		return nil, fmt.Errorf("string length must be positive")
	}

	message := make([]byte, length)
	if _, err := rand.Read(message); err != nil {
		return nil, fmt.Errorf("failed to generate random message: %w", err)
	}

	return message, nil
}

// SolveChallenge attempts to solve given challenge with given difficulty using a solution of given length.
func SolveChallenge(ctx context.Context, challenge []byte, difficulty, length int) (solution []byte, err error) {

	for ctx.Err() == nil && !VerifySolution(challenge, solution, difficulty) {
		if solution, err = GenerateRandomBytes(length); err != nil {
			return nil, err
		}
	}

	return solution, ctx.Err()
}

// VerifySolution takes a PoW challenge, a difficulty level, and a solution as inputs.
// It appends the solution to the challenge, computes the SHA256 hash of the resulting string,
// and checks whether the first `difficulty` bytes of the hash are zero.
// If they are, the function returns `true`, indicating that the solution is valid.
func VerifySolution(challenge, solution []byte, difficulty int) bool {

	hash := sha256.Sum256(append(challenge, solution...))

	for i := 0; i < difficulty; i++ {
		if hash[i] != 0 {
			return false
		}
	}

	return true
}
