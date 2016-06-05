default: sos-server

clean:
	go clean

fmt:
	go fmt *.go

sos-server: $(wildcard *.go)
	go build .
