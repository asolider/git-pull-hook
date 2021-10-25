
BINARY_NAME=git-pull


.PHONY:build
build:
	go build -o $(BINARY_NAME) -v

build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o $(BINARY_NAME) -v

build-mac:
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o $(BINARY_NAME) -v

clean:
	rm -f $(BINARY_NAME)