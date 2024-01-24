build:
	echo "Building ..."
	CGO_ENABLED=0 go build -o bin/server cmd/main.go

run:
	bin/server -env resources/.env
