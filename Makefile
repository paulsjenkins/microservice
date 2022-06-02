all: run

run:
	go run cmd/main.go


# replace IMAGE_NAME with the name to use for your image
docker:
	docker build -t IMAGE_NAME .