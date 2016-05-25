Simple Object Storage, in golang
--------------------------------

The Simple Object Storage (SOS) is a HTTP-based object-storage system which allows files to be uploaded, and later retrieved by ID.

Files can be replicated across a number of hosts to ensure redundancy, and increase availability in the event of hardware failure.

* [The design of the system](DESIGN.md).
* [Scaling to large numbers of objects](SCALING.md).
* [How replication works](REPLICATION.md).
* [The APIs we present, both internal and private](API.md).


Installation
------------

Building the code should pretty idiomatic for a golang user:

     #
     # Download the code to $GOPATH/src
     # If already present is should be updated.
     #
     go get -u github.com/skx/sos/...

If you prefer to build manually:

     $ git clone https://github.com/skx/sos.git
     $ cd sos
     $ make

Once built you'll find three binaries:

| Path          | Purpose                                    |
|---------------|--------------------------------------------|
| bin/replicate | The Perl based replication utility.        |
| sos-server    | The public-facing upload/download service. |
| blob-server   | The blob-server used for storing content.  |




Quick Start
-----------

In an ideal deployment at least two servers would be used:

* One server would run the `sos-server`, which allows uploads to be made, and later retrieved.
* Each of the two servers would run a `blob-server`, allowing an object to be replicated upon both hosts.

We can simulate this upon a single host though, for the purposes of testing.  You'll just need to make sure you have four terminals open to run the appropriate daemons.

First of all you'll want to launch a pair of blob-servers:

    $ blob-server -store data1 -port 4001
    $ blob-server -store data2 -port 4002

> **NOTE**: The storage-paths (`./data1` and `./data2` in the example above) is where the uploaded-content will be stored.  These directories will be created if missing.

In production usage you'd generally record the names of the blob-servers in a configuration file, either `/etc/sos.conf`, or `~/.sos.conf`, however they may also be specified upon the command line.

We start the `sos-server` ensuring that it knows about the blob-servers to store content in:

    $ sos-server -blob-server http://localhost:4001,http://localhost:4002
    Launching API-server
    ..


Now you, or your code, can connect to the server and start uploading/downloading objects.  By default the following ports will be used by the `sos-server`:

|service           | port |
|----------------- | ---- |
| upload service   | 9991 |
| download service | 9992 |

Providing you've started all three daemons you can now perform a test upload with `curl`:

    $ curl -X POST --data-binary @/etc/passwd  http://localhost:9991/upload
    {"id":"cd5bd649c4dc46b0bbdf8c94ee53c1198780e430","size":2306,"status":"OK"}

If all goes well you'll receive a JSON-response as shown, and you can use the ID which is returned to retrieve your download:

    $ curl http://localhost:9992/fetch/cd5bd649c4dc46b0bbdf8c94ee53c1198780e430
    ..
    $

> **NOTE**: The download service runs on a different port.  This is so that you can make policy decisions about uploads/downloads via your local firewall.

At the point you run the upload the contents will only be present on one of the blob-servers, chosen at random, to ensure that it is mirrored you'll want to replicate the contents:

    $ ./bin/replicate --verbose





Production Usage
----------------

* The API service must be visible to clients, to allow downloads to be made.
    * Because the download service runs on port `9992` it is assumed that corporate firewalls would deny access.
    * We assume you'll configure an Apache/nginx/similar reverse-proxy to access the files via a host like http://objects.example.com/`.

* It is assumed you might wish to restrict uploads to particular clients, rather than allow the world to make uploads.  The simplest way of doing this is to use a local firewall.

* The blob-servers should be reachable by the hosts running the API-service, but they should not be publicly visible.
    * If your blob-servers are exposed to the internet remote users could [use the API](API.md) to download all your content.

* None of the servers need to be launched as root, because they don't bind to privileged ports, or require special access.
    * **NOTE**: [issue #6](https://github.com/skx/sos/issues/6) improved the security of the `blob-server` by invoking `chroot()`.  However `chroot()` will fail if the server is not launched as root, which is harmless.

* Scaling when your data is too large to fit upon a single `blob-server`:
   * [Read about scaling SoS](SCALING.md)


Future Changes?
---------------

It would be possible to switch to using _chunked_ storage, for example breaking up each file that is uploaded into 128Mb sections and treating them as distinct.  The reason that is not done at the moment is because it relies upon state:

* The public server needs to be able to know that the file with a given ID is comprised of the following chunks of data:
    * `a5d606958533634fed7e6d5a79d6a5617252021f`
    * `038deb6940db2d0e7b9ee9bba70f3501a0667989`
    * `a7914eb6ff984f97c5f6f365d3d93961be2e8617`
    * `...`
* That data must be always kept up to date and accessible.

At the moment the API-server is stateless, so tracking that data is not possible.  It possible to imagine using [redis](http://redis.io/), or some other external database to record the data, but that increases the complexity of deployment.


Questions?
----------

Questions/Changes are most welcome; just report an issue.

Steve
--
