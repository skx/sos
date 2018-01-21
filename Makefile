#
# Trivial Makefile for sos
#


#
# Build our binary by default
#
all: sos


#
# Build our main binary
#
sos: $(wildcard *.go)
	go build .


#
# Run our tests
#
test:
	go test -coverprofile fmt

#
# Clean our build
#
clean:
	rm sos cover.out foo.html || true


cover:
	go test -coverprofile=cover.out

cover-html:
	go test -coverprofile=cover.out
	go tool cover -html=cover.out -o foo.html
	firefox foo.html
