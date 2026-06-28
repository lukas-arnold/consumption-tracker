FROM golang:1.26 AS builder

WORKDIR /app

COPY go.mod ./
COPY . .

RUN CGO_ENABLED=0 go build -o app ./cmd/consumption-tracker

FROM alpine:3.20

WORKDIR /app

COPY --from=builder /app/app ./consumption-tracker

CMD ["./consumption-tracker"]