FROM golang:1.17

WORKDIR /app

COPY go.* ./

RUN go mod download

COPY *.go ./

RUN go build -o api

EXPOSE 8080

ENTRYPOINT ["./api"]