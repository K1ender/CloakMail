FROM golang:1.23.4-alpine3.21 AS builder

WORKDIR /app

COPY . .

RUN go mod download

RUN go build -ldflags="-s -w" -o main ./cmd/main.go

FROM alpine:3

WORKDIR /app

COPY --from=builder /app/main .

CMD [ "./main" ]