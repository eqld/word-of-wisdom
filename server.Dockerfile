FROM golang:1.20-alpine

ENV WOW_SERVER_DIFFICULTY=2
ENV WOW_SERVER_CHALLENGE_LENGTH=16
ENV WOW_SERVER_SOLUTION_LENGTH=8
ENV WOW_SERVER_CONN_HANDLE_TIMEOUT_SECONDS=15

RUN apk add fortune

WORKDIR /app

COPY . ./
RUN go mod download
RUN go test ./...
RUN go build -o word-of-wisdom-server ./cmd/server

EXPOSE 8080

ENTRYPOINT ["./word-of-wisdom-server"]
