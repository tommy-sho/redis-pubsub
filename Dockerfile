FROM golang:1.10.1

WORKDIR /go/src/github.com/tommy-sho/redis-pubsub
COPY . .

CMD ["go","run","."]