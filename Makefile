all: dev

dev:
	go build main.go

image:
	go mod tidy
	docker build -t azssi/my-echo-go .


push-image: image
	docker push azssi/my-echo-go

cert:
	@echo "Generate certificate for: $(CN)"
	mkdir -p tls
	openssl req -x509 -newkey rsa:4096 -sha256 -nodes -keyout tls/server.key -out tls/server.crt -subj "/CN=$(CN)" -days 3650
