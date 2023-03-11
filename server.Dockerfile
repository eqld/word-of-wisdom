FROM golang:1.20-alpine

RUN apk add fortune

WORKDIR /app

COPY . ./
RUN go mod download
RUN go test ./...
RUN go build -o word-of-wisdom-server ./cmd/server

EXPOSE 8080

ENTRYPOINT ["./word-of-wisdom-server"]
