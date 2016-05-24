
all:
	cd sos-server  && make
	cd blob-server && make

clean:
	cd sos-server  && make clean
	cd blob-server && make clean

fmt:
	cd sos-server  && make fmt
	cd blob-server && make fmt
