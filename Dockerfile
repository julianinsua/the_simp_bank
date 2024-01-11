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
COPY --from=builder /app/db/migrations ./migration
RUN apk add curl
RUN curl -fsSL https://raw.githubusercontent.com/pressly/goose/master/install.sh | GOOSE_INSTALL=/app/goose sh
COPY --from=builder /app/start.sh .
RUN ["chmod", "+x", "/app/start.sh"]
COPY --from=builder /app/wait-for.sh .


# Documentation
EXPOSE 8080

# Execution
CMD ["/app/main"]
ENTRYPOINT ["/app/start.sh"]
