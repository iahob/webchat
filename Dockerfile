FROM golang:1.19 AS builder

WORKDIR /home
ADD . .
RUN go build -o app main.go

FROM alpine
WORKDIR /app
COPY --from=builder /home/app .
CMD ["./app"]
