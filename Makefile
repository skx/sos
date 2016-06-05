
all:
	cd sos-server     && make
	cd sos-replicator && make
	cd blob-server    && make

clean:
	cd sos-server     && make clean
	cd sos-replicator && make clean
	cd blob-server    && make clean

fmt:
	cd sos-server      && make fmt
	cd sos-replicator  && make fmt
	cd blob-server     && make fmt
