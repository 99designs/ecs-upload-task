VERSION=$(shell git describe --tags --candidates=1 --dirty 2>/dev/null || echo "dev")
FLAGS=-X main.Version=$(VERSION)

release-multi-arch: clean
	mkdir -p bin
	gox -os="linux darwin windows" -arch="amd64 arm64" -ldflags="$(FLAGS)" -output="./bin/{{.Dir}}_{{.OS}}_{{.Arch}}"      

clean:
	rm -rf bin
	