Design Overview
---------------

The implementation of the object-store is built upon the primitive of a "blob server".  A blob server is a dumb service which provides three simple operations:

* Store a particular chunk of binary data with a specific name.
* Given a name retrieve the chunk of binary data associated with it.
* Return a list of all known names.

Ultimately the blob-servers provide the storage for the object-store, and the upload/download service just needs to mediate between them there is no need to maintain any local state, or meta-data.

The `blob-server` primitives are sufficient to provide a robust replicating storage system, because it is possible to easily mirror their contents, providing we assume that the IDs only ever hold a particular set of data (i.e. data is immutable).

To replicate the contents of `blob-server-a` to `blob-server-b` the algorithm is obvious:

* Get the list of known-names of the blobs stored on `blob-server-a`.
* For each name, fetch the data associated with that name.
    * Now store that data, with the same name, on `blob-server-b`.

There is more [documentation on replication](REPLICATION.md) which explains how it works in terms of groups.



Upload Operation
----------------

An upload operation involves:

* Contacting every `blob-server` in turn, and attempting the upload.
   * If an upload succeeds return the data to the client.
* If every `blob-server` has been contacted, and the upload failed, then we return an HTTP 500 error-code to the caller.


Download Operation
------------------

A download operation is similar to an upload:

* For every known `blob-server` we request the specified object.
    * If the content is found it is returned to the caller.
* If the content has not been found and all the known blob-servers have been queried then the object is missing, and a HTTP 404 status-code is returned to the caller.
