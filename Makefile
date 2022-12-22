.SILENT:

go-build:
	go mod download && CGO_ENABLED=0 GOOS=linux go build -o ./build/inHabrBot cmd/inHabrBot/main.go

docker-build: go-build
	docker build . -t in_habr_bot

run: docker-build
	docker-compose up