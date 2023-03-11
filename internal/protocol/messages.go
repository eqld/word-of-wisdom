package protocol

import (
	"fmt"
	"strconv"
	"strings"
)

// FormatChallengeForClient formats message with challenge and difficulty for client.
func FormatChallengeForClient(challenge string, difficulty int) string {
	return fmt.Sprintf("%s:%d", challenge, difficulty)
}

// ParseChallengeForClient parses message with challenge and difficulty for client.
func ParseChallengeForClient(challengeWithDifficulty string) (challenge string, difficulty int, err error) {
	cd := strings.SplitN(strings.TrimRight(challengeWithDifficulty, "\n"), ":", 2)
	if len(cd) != 2 {
		return "", 0, fmt.Errorf("wrong format of server message with challenge")
	}

	challenge = cd[0]

	if difficulty, err = strconv.Atoi(cd[1]); err != nil {
		return "", 0, fmt.Errorf("wrong format of difficulty in server message: %w", err)
	}

	return challenge, difficulty, nil
}
