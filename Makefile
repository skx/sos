
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


#
#  Targets solely used for testing.
#

d2000:
	./blob-server/blob-server -store ./data.2000 -port 2000

d2001:
	./blob-server/blob-server -store ./data.2001 -port 2001

api:
	./sos-server/sos-server -blob-server http://localhost:2000,http://localhost:2001

rep:
	./sos-replicator/sos-replicator -blob-server http://localhost:2000,http://localhost:2001 -verbose
