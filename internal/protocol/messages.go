package protocol

import (
	"fmt"
	"strconv"
	"strings"
)

// FormatChallengeForClient formats message with challenge and difficulty for client.
func FormatChallengeForClient(challenge string, difficulty, solutionLength int) string {
	return fmt.Sprintf("%s:%d:%d", challenge, difficulty, solutionLength)
}

// ParseChallengeForClient parses message with challenge and difficulty for client.
func ParseChallengeForClient(message string) (challenge string, difficulty, solutionLength int, err error) {
	messageParsed := strings.SplitN(strings.TrimRight(message, "\n"), ":", 3)
	if len(messageParsed) != 3 {
		return "", 0, 0, fmt.Errorf("wrong format of server message with challenge")
	}

	challenge = messageParsed[0]
	if challenge == "" {
		return "", 0, 0, fmt.Errorf("empty challenge in server message: %w", err)
	}

	if difficulty, err = strconv.Atoi(messageParsed[1]); err != nil {
		return "", 0, 0, fmt.Errorf("wrong format of difficulty in server message: %w", err)
	}

	if solutionLength, err = strconv.Atoi(messageParsed[2]); err != nil {
		return "", 0, 0, fmt.Errorf("wrong format of solution length in server message: %w", err)
	}

	return challenge, difficulty, solutionLength, nil
}
