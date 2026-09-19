FROM golang:1.27.0 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/match-service ./cmd/match-service

FROM alpine:latest
COPY --from=builder /bin/match-service /bin/match-service
EXPOSE 8083
ENTRYPOINT ["/bin/match-service"]
