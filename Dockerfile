# Stage 1: Build the application
FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Stage 2: Create the final image
FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/main .
COPY --from=builder /app/server.crt .
COPY --from=builder /app/server.key .
COPY --from=builder /app/config.yaml .

EXPOSE 8999

CMD ["./main"]
