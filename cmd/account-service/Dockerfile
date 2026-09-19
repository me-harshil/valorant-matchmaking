FROM golang:1.27.0 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/account-service ./cmd/account-service

FROM alpine:latest
COPY --from=builder /bin/account-service /bin/account-service
EXPOSE 8081
ENTRYPOINT ["/bin/account-service"]