default: blob-server

clean:
	go clean

fmt:
	go fmt *.go

blob-server: $(wildcard *.go)
	go build .
