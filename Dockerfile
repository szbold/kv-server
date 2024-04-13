FROM golang:latest AS builder

WORKDIR /go/src/kv-server

COPY . .

RUN go build .


FROM builder AS tester

RUN go test ./...


FROM alpine:3.14 AS kv-server

COPY --from=builder /go/src/kv-server/kv-server /usr/local/bin/kv-server 

RUN apk add libc6-compat

EXPOSE 6379

CMD ["kv-server"]
