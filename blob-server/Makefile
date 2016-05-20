default: blob_server

clean:
	go clean

fmt:
	go fmt *.go

libs:
	test -d $$GOPATH/github.com/gorilla/mux || go get github.com/gorilla/mux


blob_server: libs 
	go build .
