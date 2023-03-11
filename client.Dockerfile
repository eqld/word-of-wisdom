FROM golang:1.20-alpine AS build

WORKDIR /app

COPY . ./
RUN go mod download
RUN go test ./...
RUN go build -o word-of-wisdom-client ./cmd/client

FROM alpine:latest

WORKDIR /app

COPY --from=build /app/word-of-wisdom-client .

ENTRYPOINT ["./word-of-wisdom-client"]
