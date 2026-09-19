FROM golang:1.27.0 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/rating-service ./cmd/rating-service

FROM alpine:latest
COPY --from=builder /bin/rating-service /bin/rating-service
EXPOSE 8082
ENTRYPOINT ["/bin/rating-service"]
