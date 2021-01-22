FROM golang:1.14

EXPOSE 8080

WORKDIR /opt/app

COPY . .

RUN go build 

ENTRYPOINT ["/opt/app/echo-server-go"]
