FROM golang:1.27.0 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/matchmaking-service ./cmd/matchmaking-service

FROM alpine:latest
COPY --from=builder /bin/matchmaking-service /bin/matchmaking-service
EXPOSE 8084
ENTRYPOINT ["/bin/matchmaking-service"]
