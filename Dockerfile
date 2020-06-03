FROM golang:1.14.3 as builder
RUN apt-get update
RUN apt-get install -y make build-essential
WORKDIR /usr/src/app
COPY . .
RUN go build -o ./ws2redis ./src/main.go
ENTRYPOINT ["./ws2redis"]
