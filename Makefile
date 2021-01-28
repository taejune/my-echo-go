

all: dev

dev:
	@echo "Build docker image..."
	docker build -t azssi/my-echo-go .
