VERSION=$(shell git describe --tags --candidates=1 --dirty 2>/dev/null || echo "dev")
FLAGS=-X main.Version=$(VERSION)

all: darwin linux windows

darwin:
	GOOS=darwin GOARCH=amd64 go build -ldflags="$(FLAGS)" -o bin/ecs-upload-task-osx-$(VERSION) .

linux:
	GOOS=linux GOARCH=amd64 go build -ldflags="$(FLAGS)" -o bin/ecs-upload-task-linux-$(VERSION) .

windows:
	GOOS=windows GOARCH=amd64 go build -ldflags="$(FLAGS)" -o bin/ecs-upload-task-win-$(VERSION) .