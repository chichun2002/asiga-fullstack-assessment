FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/main

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/main .
CMD ["./main"]