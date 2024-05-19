build:
	@echo "Building Go Lambda function"
	@docker run --rm -v $(PWD):/app -w /app -e GOOS=linux -e GOARCH=amd64 golang:latest go build -o avolta

deploy-staging:build
	rsync -a avolta ametory@103.172.205.9:/home/ametory/araya -v --stats --progress


