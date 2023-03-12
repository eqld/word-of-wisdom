package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"github.com/eqld/word-of-wisdom/internal/env"
	"github.com/eqld/word-of-wisdom/internal/pow"
	"github.com/eqld/word-of-wisdom/internal/protocol"
)

const (
	exitCodeWrongUsage = iota + 1
	exitCodeWrongMessageFormat
	exitCodeFailedToConnect
	exitCodeFailedToReadFromConn
	exitCodeFailedToWriteToConn
	exitCodeFailedToSolveChallenge
)

const (
	defaultConnTimeoutSeconds = 3
)

func main() {

	connTimeoutSeconds := env.MustReadIntEnv("WOW_CLIENT_CONN_TIMEOUT_SECONDS", defaultConnTimeoutSeconds, exitCodeWrongUsage)

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Duration(connTimeoutSeconds)*time.Second)
	defer cancel()

	if len(os.Args) < 2 {
		log.Println("usage: word-of-wisdom-client <host:port>")
		os.Exit(exitCodeWrongUsage)
	}
	addr := os.Args[1]

	conn, err := net.Dial("tcp", addr)
	if err != nil {
		log.Println("failed to connect to server:", err)
		os.Exit(exitCodeFailedToConnect)
	}
	defer conn.Close()

	go func() {
		<-ctx.Done()
		conn.Close()
	}()

	message, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		log.Println("failed to receive challenge from server:", err)
		os.Exit(exitCodeFailedToReadFromConn)
	}

	challenge, difficulty, solutionLength, err := protocol.ParseChallengeForClient(message)
	if err != nil {
		log.Println("failed to parse challenge with difficulty:", err)
		os.Exit(exitCodeWrongMessageFormat)
	}

	solution, err := pow.SolveChallenge(ctx, challenge, difficulty, solutionLength)
	if err != nil {
		log.Println("failed to solve challenge:", err)
		os.Exit(exitCodeFailedToSolveChallenge)
	}

	if _, err = fmt.Fprintln(conn, solution); err != nil {
		log.Println("failed to send a solution to server:", err)
		os.Exit(exitCodeFailedToWriteToConn)
	}

	quote, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		log.Println("failed to receive quote from server:", err)
		os.Exit(exitCodeFailedToReadFromConn)
	}

	fmt.Println()
	fmt.Println("WOW QUOTE >>>", quote)
	fmt.Println()
}
