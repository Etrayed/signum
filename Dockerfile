# syntax=docker/dockerfile:1

FROM golang:1.20-alpine

WORKDIR /app

ADD config/ config/
ADD processor/ processor/

COPY go.mod go.sum ./
RUN go mod download

COPY *.go .

RUN CGO_ENABLED=0 GO111MODULE=on GOOS=linux GOARCH=amd64 go build -o /signum -ldflags="-X main.Version=1.0.0"

CMD [ "/signum" ]