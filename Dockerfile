FROM golang:1.19 AS builder

WORKDIR /home
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o app main.go
RUN go env

FROM alpine
RUN apk add bash
WORKDIR /app
COPY --from=builder /home/app .
CMD ["./app"]
