all: dev

dev:
	go build main.go

docker-image:
	docker build -t azssi/my-echo-go .


docker-push: docker-image
	docker push azssi/my-echo-go

cert:
	@echo "Generate certificate for: $(CN)"
	openssl req -x509 -newkey rsa:4096 -sha256 -nodes -keyout key.pem -out cert.pem -subj "/CN=$(CN)" -days 3650