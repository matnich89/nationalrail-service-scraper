FROM golang:1.22 AS builder

RUN mkdir /app
COPY . /app
WORKDIR /app

RUN CGO_ENABLED=0 GOOS=linux go build -o app main.go

EXPOSE 8080

FROM scratch AS production
COPY --from=builder /app .
CMD ["./app"]