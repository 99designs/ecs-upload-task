VERSION=$(shell git describe --tags --candidates=1 --dirty 2>/dev/null || echo "dev")
FLAGS=-X main.Version=$(VERSION)

all: darwin linux windows

darwin:
	GOOS=darwin GOARCH=amd64 go build -ldflags="$(FLAGS)" -o bin/ecs-upload-task-osx-$(VERSION) .

linux:
	GOOS=linux GOARCH=amd64 go build -ldflags="$(FLAGS)" -o bin/ecs-upload-task-linux-$(VERSION) .

windows:
	GOOS=windows GOARCH=amd64 go build -ldflags="$(FLAGS)" -o bin/ecs-upload-task-win-$(VERSION) .

release-multi-arch:
	mkdir -p bin
	gox -os="linux darwin windows" -arch="amd64 arm64" -output="./bin/{{.Dir}}_{{.OS}}_{{.Arch}}"      

clean:
	rm -f bin/*