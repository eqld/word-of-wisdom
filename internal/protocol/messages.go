package protocol

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
)

// ChallengeEncode encodes message with challenge, difficulty and solution length.
func ChallengeEncode(challenge []byte, difficulty, solutionLength int) (message string) {
	challengeStr := hex.EncodeToString(challenge)
	return fmt.Sprintf("%s:%d:%d", challengeStr, difficulty, solutionLength)
}

// ChallengeDecode decodes message with challenge, difficulty and solution length.
func ChallengeDecode(message string) (challenge []byte, difficulty, solutionLength int, err error) {
	messageParsed := strings.SplitN(strings.TrimRight(message, "\n"), ":", 3)
	if len(messageParsed) != 3 {
		return nil, 0, 0, fmt.Errorf("wrong format of server message with challenge")
	}

	challengeStr := messageParsed[0]
	if challengeStr == "" {
		return nil, 0, 0, fmt.Errorf("empty challenge in server message")
	}

	if challenge, err = hex.DecodeString(challengeStr); err != nil {
		return nil, 0, 0, fmt.Errorf("wrong format of challenge in server message: %w", err)
	}

	if difficulty, err = strconv.Atoi(messageParsed[1]); err != nil {
		return nil, 0, 0, fmt.Errorf("wrong format of difficulty in server message: %w", err)
	}

	if solutionLength, err = strconv.Atoi(messageParsed[2]); err != nil {
		return nil, 0, 0, fmt.Errorf("wrong format of solution length in server message: %w", err)
	}

	return challenge, difficulty, solutionLength, nil
}

// SolutionEncode encodes message with solution.
func SolutionEncode(solution []byte) (message string) {
	return hex.EncodeToString(solution)
}

// SolutionDecode decodes message with solution.
func SolutionDecode(message string) (solution []byte, err error) {
	if message == "" {
		return nil, fmt.Errorf("empty solution in client message")
	}
	return hex.DecodeString(message)
}

// QuoteEncode encodes a quote to base64 message to preserve arbitrary newlines.
func QuoteEncode(quote string) (message string) {
	return base64.StdEncoding.EncodeToString([]byte(quote))
}

// QuoteDecode decodes a quote from base64 message
func QuoteDecode(message string) (quote string, err error) {
	message = strings.TrimRight(message, "\n")

	if message == "" {
		return "", fmt.Errorf("empty quote in server message")
	}

	quoteBytes, err := base64.StdEncoding.DecodeString(message)
	if err != nil {
		return "", fmt.Errorf("wrong quote message format: %w", err)
	}

	return string(quoteBytes), nil
}
