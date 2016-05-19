Simple Object Storage, in golang
--------------------------------

The Simple Object Storage (SOS) is a HTTP-based object-storage system
which allows files to be uploaded, and later retrieved by ID.

Files can be replicated across a number hosts to ensure redundancy,
and despite the naive implementation it does scale to millions of files.

The code written in [golang](http://golang.com/), easing deployment.

Building the code should be a simple as:

     cd ./golang
     make


Simple Design
-------------

The implementation of the object-store is built upon the primitive of a "blob server".  A blob server is a dumb service which provides three simple operations:

* Store a particular chunk of binary data with a specific name.
* Given a name retrieve the chunk of binary data associated with it.
* Return a list of all known names.

These primitives are sufficient to provide a robust replicating storage system, because it is possible to easily mirror their contents, providing we assume that the IDs only ever hold a particular set of data (i.e. data is immutable).

To replicate the contents of `blob-server-a` to `blob-server-b` the algorithm is obvious:

* Get the list of known-names of the blobs stored on `blob-server-a`.
* For each name, fetch the data associated with that name.
    * Now store that data, with the same name, on `blob-server-b`.

In real world situations the replication might become more complex over time, as different blob-servers might be constrained by differing amounts of disk-space, etc.  But the core-operation is both obvious and simple to implement.

(In the future you could imagine switching to from the HTTP-based blob-server to using something else: [redis](http://redis.io/), [memcached](https://memcached.org/), or [postgresql](http://postgresql.org/) would be obvious candidates!)

Ultimately the blob-servers provide the storage for the object-store, and the upload/download service just needs to mediate between them.  There isn't fancy logic or state to maintain, beyond that local to each node, so it is possible to run multiple blob-servers and multiple API-servers if required.

The important thing is to ensure that a replication-job is launched regularly, to ensure that blob-servers __are__ replicated:

    ./bin/replicate -v


Quick Start
-----------

In an ideal deployment at least two servers would be used:

* One server would run the API-server, which allows uploads to be made, and later retrieved.
* Each of the two servers would run a blob-service, allowing a single uploaded object to be replicated upon both hosts.

We can replicate this upon a single host though, for the purposes of testing.  You'll just need to make sure you have four terminals open to run the appropriate daemons.

First of all you'll want to launch a pair of blob-servers:

    ./golang/blob_server -store data1 -port 4001
    ./golang/blob_server -store data2 -port 4002

Record the names of the server in a configuration file, such that the upload/download daemons know where their storage is located at:

    $ cat >> ~/.sos.conf<<EOF
    http://localhost:4001
    http://localhost:4002
    EOF

Now you can start an API-server.  This server is what your code will interact with to upload content, and it will talk to the blob-servers to actually store your uploads on-disk:

    ./golang/api_server
    Launching API-server
    ..


By default the following ports will be used by the `api_server`:

|service          | port |
|---------------- | ---- |
| upload service   | 9991 |
| download service | 9992 |

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

* The API service must be visible to clients, to allow downloads to be made.

* It is assumed you might wish to restrict uploads to particular clients, rather than allow the world to make uploads.  The simplest way of doing this is to use a local firewall.

* The blob-servers should be reachable by the hosts running the API-service, but they do not need to be publicly visible, these should be firewalled.

* None of the servers need to be launched as root, because they don't bind to privileged ports, or require special access.
    * **NOTE**: Once [issue #6](https://github.com/skx/sos/issues/6) is implemented root privileges will be required for a successful `chroot()`, however if that fails things are not terrible.



Future Changes?
---------------

There are two specific changes which would be useful to see in the future:

* Marking particular blob-servers as preferred, or read-only.
     * If you have 10 severs, 8 of which are full, then it becomes useful to know that explicitly, rather than learning at runtime when many operations have to be retried, repeated, or canceled.
* Caching the association between object-IDs and the blob-server(s) upon which it is stored.
     * This would become more useful as the number of the blob-servers rises.

It would be possible switch to using _chunked_ storage, for example breaking up each file that is uploaded into 128Mb sections and treating them as distinct.  The reason that is not done at the moment is because it relies upon state:

* The public server needs to be able to know that the file with ID "NNNNABCDEF1234" is comprised of chunks "11111", "222222", "AAAAAA", "BBBBBB", & etc.
* That data must be always kept up to date and accessible.

At the moment the API-server is stateless.  You could even run 20 of them, behind a load-balancer, with no concerns about locking or sharing!  Adding state spoils that, and the complexity has not yet been judged worthwhile.


Questions?
----------

Questions/Changes are most welcome; just report an issue.

Steve
--
