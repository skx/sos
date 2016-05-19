Simple Object Storage
---------------------

The Simple Object Storage (SOS) is a HTTP-based object-storage system
which allows files to be uploaded, and later retrieved by ID.

Files can be replicated across a number hosts to ensure redundancy,
and despite the naive implementation it scales easily to thousands of files.


Dependencies
------------

The code is 100% pure-perl, and requires only a minimal set of dependencies, which are available to the Stable release of Debian GNU/Linux:

    apt-get install libdancer-perl libjson-perl

For production use I'd recommend the use of plack:

    apt-get install libplack-perl twiggy

As a proof of concept there is a golang version of the blob-server, included beneath `golang/` this is not yet 100% functional, but will be shortly.  If you have performance concerns you might consider using it.


Overview
--------

The implementation is split into three distinct parts:

* A (public) API which you interactive with.
    * Allowing you to upload/store files.
    * Allowing you to download/retrieve files.

* A blob-server.
    * The blob-servers are the things that actually store the data.
    * Multiple blob-servers can be executed on multiple hosts, and data will be replicated between them.

* A replication utility.

**NOTE**: In the past the public API-service was provided by one server, but in the wild it is expected you'd have security goals that would disallow that:

* The expectation is that everybody can download all stored-objects.
* But uploads should only be permitted from one or more hosts.
    * Splitting the api-server into a pair of distinct upload/download services allows a firewall to be applied more useful.
    * The only cost is the requirement to run two daemons instead of one.



Quick Start
-----------

In an ideal deployment at least two servers would be used:

* One server would run the upload & download service, allowing files to stored and retrieved.
* Each of the two servers would run a blob-service, allowing a single uploaded object to be replicated upon both hosts.

We can replicate this upon a single host though, for the purposes of testing.  You'll just need to make sure you have four terminals open to run the appropriate daemons.

First of all you'll want to launch a pair of blob-servers:

    ./bin/blob-server --port 4040 ./data1
    ./bin/blob-server --port 4041 ./dat2

  Record their names in a configuration file, such that the upload/download daemons know where their storage is located at:

    $ cat >> ~/.sos.conf<<EOF
    http://localhost:4040
    http://localhost:4041
    EOF

Now you can start an upload-server.  This server is what your code will interact with to upload content, and it will talk to the blob-servers to actually store your uploads on-disk:

    ./bin/upload-server

Finally you'll want to launch a download-server, which is what clients will connect to in order to retrieve previously-uploaded content:

    ./bin/download-server


By default the following ports will be used:

|service          | port |
|---------------- | ---- |
| upload-server   | 9991 |
| download-server | 9992 |

To upload a file to your server, to test it out run:

    $ curl -X POST --data-binary @/etc/passwd  http://localhost:9991/upload
    {"id":"cd5bd649c4dc46b0bbdf8c94ee53c1198780e430","size":2306,"status":"OK"}

To fetch this uploaded-file:

    $ curl http://localhost:9992/fetch/cd5bd649c4dc46b0bbdf8c94ee53c1198780e430
    ..
    $

At the point you run the upload the contents will only be present on one of the blob-servers, to ensure that it is mirrored you'll want to replicate the contents:

    $ ./bin/replicate --verbose

The default is to replicate all files into two servers, if you were running three blob-servers you could ensure that each one has all the files:

    $ ./bin/replicate --verbose --min-copies=3




Production Usage
----------------

Launching the various services under `plackup` with the `Twiggy` driver will give better performance and more scalability, and can be done like so:

**NOTE**: We're launching the upload server only listening upon the loopback adapter:

   $ plackup -Ilib/ -s Twiggy --workers=4 -0 127.0.0.1 -p 9991 -a bin/upload-server --access-log logs/uploads.log -E production

The download-server is listening upon all interfaces, as it should be publicly accessible:

   $ plackup -Ilib/  -s Twiggy --workers=4 -0 0.0.0.0 -p 9992 -a bin/download-server --access-log logs/downloads.log -E production

Otherwise the usage is the same as previously, recording the server-names in the configuration file `/etc/sos.conf`, or `~/.sos.conf`, and launching the appropriate blob-servers:

    $ bin/blob-server --port=2001  ./data1
    $ bin/blob-server --port=2002  ./data2

Finally on the upload/download-server host(s) define the list of servers in /etc/sos.conf:

    $ cat > /etc/sos.conf <<EOF
    server1.storage.lan:2001
    server2.storage.lan:2002
    EOF

**NOTE**: The blob-servers should be firewalled; they explicitly do not need to be publicly accessible.


Questions?
----------

Questions/Changes are most welcome; just report an issue.

Steve
-- 
