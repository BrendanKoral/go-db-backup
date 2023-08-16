BINARY_NAME=backup-all-dbs
.DEFAULT_GOAL := build

build:
	GOARCH=amd64 GOOS=linux go build -o ./bin/${BINARY_NAME}-linux main.go