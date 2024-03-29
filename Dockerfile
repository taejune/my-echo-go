FROM  golang:1.19

WORKDIR /usr/src/app

COPY go.mod go.sum ./
COPY pkg ./pkg
COPY main.go ./

RUN go mod download && go mod verify

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -ldflags '-s -w' -o /usr/local/bin/app main.go

EXPOSE 8080

CMD ["app"]
