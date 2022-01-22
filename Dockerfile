FROM golang:1.16-alpine

WORKDIR $GOPATH/github.com/airbenders/profile

COPY . .

RUN go get -d -v ./...

EXPOSE 8080

CMD ["go", "run", "main.go"]