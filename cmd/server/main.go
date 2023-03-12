package main

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/eqld/word-of-wisdom/internal/env"
	"github.com/eqld/word-of-wisdom/internal/pow"
	"github.com/eqld/word-of-wisdom/internal/protocol"
)

const (
	exitCodeWrongParam = iota + 1
	exitCodeFailedToListen
)

const (
	defaultDifficulty         = 2
	defaultChallengeLength    = 16
	defaultSolutionLength     = 8
	defaultConnTimeoutSeconds = 3
)

func main() {

	difficulty := env.MustReadIntEnv("WOW_SERVER_DIFFICULTY", defaultDifficulty, exitCodeWrongParam)
	challengeLength := env.MustReadIntEnv("WOW_SERVER_CHALLENGE_LENGTH", defaultChallengeLength, exitCodeWrongParam)
	solutionLength := env.MustReadIntEnv("WOW_SERVER_SOLUTION_LENGTH", defaultSolutionLength, exitCodeWrongParam)
	connTimeoutSeconds := env.MustReadIntEnv("WOW_SERVER_CONN_TIMEOUT_SECONDS", defaultConnTimeoutSeconds, exitCodeWrongParam)

	log.Printf("starting the server with difficulty '%v', challenge length '%v' and solution length '%v'",
		difficulty, challengeLength, solutionLength)

	// Set up TCP listener.

	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Println("failed to create tcp listener:", err)
		os.Exit(exitCodeFailedToListen)
	}
	defer listener.Close()

	log.Println("listening on", listener.Addr())
	defer log.Println("terminating")

	// Create execution context.

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Handle system signals.

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sig
		listener.Close()
		cancel()
	}()

	// Handle incoming connections.

	h := handler{
		difficulty:      difficulty,
		challengeLength: challengeLength,
		solutionLength:  solutionLength,
		connTimeout:     time.Duration(connTimeoutSeconds) * time.Second,
	}

	go func() {
		for connNum := 0; ctx.Err() == nil; connNum++ {

			conn, err := listener.Accept()
			if err != nil {
				if errors.Is(err, net.ErrClosed) {
					// Listener closed.
					return
				}

				logf(connNum, "failed to accept new tcp connection: %v", err)
				continue
			}

			go h.handleConnection(ctx, connNum, conn)
		}
	}()

	// Wait for termination.
	<-ctx.Done()
}

type handler struct {
	difficulty      int
	challengeLength int
	solutionLength  int
	connTimeout     time.Duration
}

// handleConnection handles given client connection.
// In case of any error it logs error message and closes the connection.
func (h handler) handleConnection(ctx context.Context, connNum int, conn net.Conn) {
	ctx, cancel := context.WithTimeout(ctx, h.connTimeout)
	defer cancel()

	// Close connection in case of timeout.

	go func() {
		<-ctx.Done()
		conn.Close()
	}()

	// Generate challenge for client.

	challenge, err := pow.GenerateRandomString(h.challengeLength)
	if err != nil {
		logf(connNum, "failed to generate callenge: %v", err)
		return
	}
	logf(connNum, "generated challenge '%v'", challenge)

	// Send challenge with difficulty to client.

	message := protocol.FormatChallengeForClient(challenge, h.difficulty, h.solutionLength)
	if _, err = fmt.Fprintln(conn, message); err != nil {
		logf(connNum, "failed to send challenge: %v", err)
		return
	}
	logf(connNum, "challenge sent, waiting for solution")

	// Read solution from client.

	solutionBytes := make([]byte, h.solutionLength+1)
	if _, err := conn.Read(solutionBytes); err != nil {
		logf(connNum, "failed to read solution: %v", err)
		return
	}
	if solutionBytes[h.solutionLength] != '\n' {
		logf(connNum, "solution length exceeds limit '%v'", h.solutionLength)
		return
	}
	solution := string(solutionBytes[:h.solutionLength])
	logf(connNum, "received solution '%v'", solution)

	// Verify the solution.

	if pow.VerifySolution(challenge, h.difficulty, solution) {
		// Send random quote to the client if solution is correct.

		logf(connNum, "solution is correct, generating a quote")
		quote, err := getQuote(ctx)
		if err != nil {
			logf(connNum, "failed to generate response: %v", err)
			return
		}

		logf(connNum, "sending the quote to the client: %v", quote)

		if _, err = fmt.Fprintln(conn, quote); err != nil {
			logf(connNum, "failed to send response: %v", err)
			return
		}
		logf(connNum, "quote sent")
	} else {
		logf(connNum, "solution is not correct")
	}
}

// logf logs formatted message to STDERR.
func logf(connNum int, format string, args ...any) {
	msg := fmt.Sprintf(format, args...)
	log.Printf("[%d] %s\n", connNum, msg)
}

// getQuote returns a random quote from the `fortune` program.
func getQuote(ctx context.Context) (string, error) {
	cmd := exec.CommandContext(ctx, "fortune")

	out := new(bytes.Buffer)
	cmd.Stdout = out

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("failed to get random quote from `fortune`: %w", err)
	}

	return strings.TrimSpace(strings.ReplaceAll(out.String(), "\n", " ")), nil
}
