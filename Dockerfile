FROM golang:1.21-alpine as builder

WORKDIR /app

COPY . .

RUN go get github.com/swaggo/swag/gen@v1.16.1
RUN go install github.com/swaggo/swag/cmd/swag@v1.16.1
RUN swag init ./main.go
RUN go build -o main.bin main.go

FROM alpine as release

WORKDIR /app

RUN apk add --no-cache bash
COPY --from=builder /app/main.bin /app/main.bin
COPY wait-for-it.sh /app/wait-for-it.sh
RUN chmod +x /app/wait-for-it.sh

EXPOSE 8080

ENTRYPOINT ["/app/wait-for-it.sh", "db", "3306", "--", "./main.bin"]