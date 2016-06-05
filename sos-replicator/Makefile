default: sos-replicator

clean:
	go clean

fmt:
	go fmt *.go

sos-replicator: $(wildcard *.go)
	go build .
