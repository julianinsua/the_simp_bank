# Builder step
FROM golang:1.21 AS builder
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o main main.go

# Run step
FROM alpine:3.19
WORKDIR /app
COPY --from=builder /app/main .
COPY app.env .

EXPOSE 8080

CMD ["/app/main"]
