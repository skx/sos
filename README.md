[![Travis CI](https://img.shields.io/travis/skx/sos/master.svg?style=flat-square)](https://travis-ci.org/skx/sos)
[![Go Report Card](https://goreportcard.com/badge/github.com/skx/sos)](https://goreportcard.com/report/github.com/skx/sos)
[![license](https://img.shields.io/github/license/skx/sos.svg)](https://github.com/skx/sos/blob/master/LICENSE)
[![Release](https://github-release-version.herokuapp.com/github/skx/sos/release.svg?style=flat)](https://github.com/skx/sos/releases/latest)

Simple Object Storage, in golang
--------------------------------

The Simple Object Storage (SOS) is a HTTP-based object-storage system which allows files to be uploaded, and later retrieved.

Files can be replicated across a number of hosts to ensure redundancy, and increased availability in the event of hardware failure.

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

Once built you'll find a single binary, `sos`, which implements a number
of sub-commands to provide functionality.



Overview
--------

You can read the [design overview](DESIGN.md) for more details, but the
SOS server relies upon the primitive of a "blob server" - which is a very
dumb service which provides three simple operations:

* Store a particular chunk of binary data with a specific name.
* Given a name retrieve the chunk of binary data associated with it.
* Return a list of all known names.

The public API is built upon the top of that primitive, and both are
launched via the same command `sos`, by specifying the sub-command
to use:

     $ ./sos blob-server ...
     $ ./sos api-server ...

Here the first command launches a blob-server, which is the back-end for
storage, and the second command launches the public API server - which is
what your code/users should operate against.

If you launch `sos` with no arguments you'll see brief details of the
available subcommands.



Quick Start
-----------

In an ideal deployment at least two hosts would be used:

* One host would run the public-server.
   * This allows uploads to be made, and later retrieved.
* Each of the two hosts would also run a blob-server.
   * The blob-servers provide the actual storage of the uploaded-objects.
   * The contents of these are replicated out of band.

We can simulate a deployment upon a single host for the purposes of testing.  You'll just need to make sure you have four terminals open to run the appropriate daemons.

First of all you'll want to launch a pair of blob-servers:

    $ sos blob-server -store data1 -port 4001
    $ sos blob-server -store data2 -port 4002

> **NOTE**: The storage-paths (`./data1` and `./data2` in the example above) is where the uploaded-content will be stored.  These directories will be created if missing.

In production usage you'd generally record the names of the blob-servers in a configuration file, either `/etc/sos.conf`, or `~/.sos.conf`, however they may also be specified upon the command line.

We'll then start the public/API-server ensuring that it knows about the blob-servers to store content in:

    $ sos api-server -blob-server http://localhost:4001,http://localhost:4002
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

If all goes well you'll receive a JSON-response as shown, and you can use the ID which is returned to retrieve your object:

    $ curl http://localhost:9992/fetch/cd5bd649c4dc46b0bbdf8c94ee53c1198780e430
    ..
    $

> **NOTE**: The download service runs on a different port.  This is so that you can make policy decisions about uploads/downloads via your local firewall.

At the point you run the upload the contents will only be present on one of the blob-servers, chosen at random.  To ensure your data is replicated you need to (regularly) launch the replication utility:

    $ sos replicate -blob-server http://localhost:4001,http://localhost:4002 --verbose
	group - server
	   default - http://localhost:4001
	   default - http://localhost:4002
    Syncing group: default
       Group member: http://localhost:4001
       Group member: http://localhost:4002
       Object cd5bd649c4dc46b0bbdf8c94ee53c1198780e430 is missing on http://localhost:4001
         Mirroring cd5bd649c4dc46b0bbdf8c94ee53c1198780e430 from http://localhost:4002 to http://localhost:4001
            Fetching :http://localhost:4002/blob/cd5bd649c4dc46b0bbdf8c94ee53c1198780e430
            Uploading :http://localhost:4001/blob/cd5bd649c4dc46b0bbdf8c94ee53c1198780e430


Meta-Data
---------

When uploading objects it is often useful to store meta-data, such as the original name of the uploaded object, the owner, or some similar data.  For that reason any header you add to your upload with an `X-`prefix will be stored and returned on download.

As a special case the header `X-Mime-Type` can be used to set the returned `Content-Type` header too.

For example uploading an image might look like this:

    $ curl -X POST -H "X-Orig-Filename: steve.jpg" \
                   -H "X-MIME-Type: image/jpeg" \
                   --data-binary @/home/skx/Images/tmp/steve.jpg \
            http://localhost:9991/upload
    {"id":"20b30df22469e6d7617c7da6a457d4e384945a06","status":"OK","size":17599}

Downloading will result in the headers being set:

    $ curl -v http://localhost:9992/fetch/20b30df22469e6d7617c7da6a457d4e384945a06 >/dev/null
    ..
    < HTTP/1.1 200 OK
    < X-Orig-Filename: steve.jpg
    < Date: Fri, 27 May 2016 06:17:39 GMT
    < Content-Type: image/jpeg
    < Transfer-Encoding: chunked
    <
    { [data not shown]




Production Usage
----------------

* The API service must be visible to clients, to allow downloads to be made.
    * Because the download service runs on port `9992` it is assumed that corporate firewalls would deny access.
    * We assume you'll configure an Apache/nginx/similar reverse-proxy to access the files via a host like `http://objects.example.com/`.

* It is assumed you might wish to restrict uploads to particular clients, rather than allow the world to make uploads.  The simplest way of doing this is to use your firewall to filter access to port `9991`.

* The blob-servers must be reachable by the host(s) running the API-service, but they should not be publicly visible.
    * If your blob-servers are exposed to the internet remote users could [use the API](API.md) to spider and download all your content.

* None of the servers need to be launched as root, because they don't bind to privileged ports, or require special access.
    * **NOTE**: [issue #6](https://github.com/skx/sos/issues/6) improved the security of the `blob-server` by invoking `chroot()`.  However `chroot()` will fail if the server is not launched as root, which is harmless.

* You can also read about scaling when your data is too large to fit upon a single `blob-server`:
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
