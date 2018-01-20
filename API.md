
API
---

There are two APIs we present, one from each sub-command:

* `sos api-server`
   * This is the public-facing API which allows objects to be uploaded & downloaded
* `sos blob-server`
   * This is an implementation-detail.


## Blob Server

The blob-server is designed to store "data" with an "id".  The data may be any binary string of arbitrary length, whereas the ID is assumed to be an alphanumeric string.

> GET /blobs

* Return a JSON array of all known object-IDs.

> POST /blob/${id}

* Store the submitted HTTP body in the blob-server, with the given ID.
* Returns a JSON array on success.

> GET /blob/${id}

* Retrieve the data associated with the specified ID, if it exists.
* Return `HTTP 404` in the event of an ID not being found.

> HEAD /blob/${id}

* Determine whether content exists for the specified ID.
* Return `HTTP 200 OK` on success.
* Return `HTTP 404` if not found.


## SOS Server

The SOS server is the public-facing server which allows uploads and downloads to be made.

**NOTE** The upload & download services run on different ports to simplify any
access-restrictions you might wish to impose.


> GET /fetch/${id}

* Fetch the content with the specified ID.
* Return `HTTP 404` on error.

> POST /upload

* Store the submitted HTTP body in the SOS-server.
* Assuming success a JSON object is returned containing the following keys:
     * `id`: The ID of the uploaded content.
     * `size`: The number of bytes received.
