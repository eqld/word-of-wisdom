package main

import (
	"bufio"
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

	"github.com/eqld/word-of-wisdom/internal/pow"
	"github.com/eqld/word-of-wisdom/internal/protocol"
)

const (
	difficulty      = 2
	challengeLength = 16
)

func main() {

	// Set up TCP listener.

	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Println("failed to create tcp listener:", err)
		os.Exit(1)
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

			go handleConnection(ctx, connNum, conn)
		}
	}()

	// Wait for termination.
	<-ctx.Done()
}

// handleConnection handles given client connection.
// In case of any error it logs error message and closes the connection.
func handleConnection(ctx context.Context, connNum int, conn net.Conn) {
	const connHandleTimeout = 15 * time.Second

	defer conn.Close()

	ctx, cancel := context.WithTimeout(ctx, connHandleTimeout)
	defer cancel()

	challenge, err := pow.GenerateRandomString(challengeLength)
	if err != nil {
		logf(connNum, "failed to generate callenge: %v", err)
		return
	}

	logf(connNum, "generated challenge '%v', difficulty is '%v'", challenge, difficulty)

	// Send challenge with difficulty to client.

	message := protocol.FormatChallengeForClient(challenge, difficulty)
	if _, err = fmt.Fprintln(conn, message); err != nil {
		logf(connNum, "failed to send challenge: %v", err)
		return
	}

	logf(connNum, "challenge sent, waiting for solution")

	// Read solution from client.

	solution, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		logf(connNum, "failed to read solution: %v", err)
		return
	}
	solution = strings.TrimRight(solution, "\n")
	logf(connNum, "received solution '%v'", solution)

	if pow.VerifySolution(challenge, difficulty, solution) {
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
