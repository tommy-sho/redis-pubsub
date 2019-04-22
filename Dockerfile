FROM golang:1.12.0-alpine AS builder

COPY . /go/src/github.com/CyberAgentHack/19fresh_m

WORKDIR /go/src/github.com/CyberAgentHack/19fresh_m/src/gateway

ENV GO111MODULE=on

RUN apk add --no-cache \
        alpine-sdk \
        git

RUN GOOS=linux GOARCH=amd64 go build -o gateway cmd/server/main.go

From alpine:latest
COPY --from=builder /go/src/github.com/CyberAgentHack/19fresh_m/src/gateway .

CMD ["./gateway"]