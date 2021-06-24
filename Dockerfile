FROM golang:latest

WORKDIR /app

COPY main.go /app

#RUN go mod download

ENTRYPOINT go run main.go

#FROM golang:1.12
#RUN mkdir short-links
#WORKDIR /short-links/
#ADD . /short-links
#RUN go build -race -ldflags "-s -w" -o short-links/bin/server short-links/server/main.go
#ENTRYPOINT ["/bin/server"]
#EXPOSE 8080
