FROM golang:1.16-alpine

ENV GOPATH=/
RUN apk add --no-cache build-base

COPY proto /usr/local/go/src/proto

COPY ./ ./

RUN go mod download

RUN go build -race -ldflags "-s -w" -o server main.go

ENTRYPOINT ["./server"]

EXPOSE 4040
