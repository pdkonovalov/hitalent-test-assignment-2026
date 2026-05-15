.PHONY: all swagger test

all: swagger test

swagger:
	go install github.com/swaggo/swag/cmd/swag@latest
	swag fmt -g cmd/server/main.go -d cmd/server,internal
	swag init -g cmd/server/main.go -o docs/

test:
	go test -v ./tests/e2e
