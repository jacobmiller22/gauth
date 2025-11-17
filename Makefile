.DEFAULT: 
	build
.PHONY: fmt vet build test

fmt:
	go fmt ./...

vet: fmt
	# go vet ./...

build: vet
	go build -o ./bin/gauth ./cmd/gauth

test: 
	go test -v ./...
