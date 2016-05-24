# Scaling with SOS

SOS is conceptually simple, and is designed to scale to host millions of files.  There are two limiting factors to the size of the storage you can expect:

* Many filesystems suffer if you place thousands of files in one directory.
* To scale appropriately you need to have the right number of blob servers.



## Getting Started

The simple case is easy to setup and use.  If __all your objects will fit upon a single server__ then you need to only launch two blob servers:

* One to hold all the objects.
* A second to hold a replica/copy of your objects.

If you're paranoid, and I would recommend this myself, then you'd run three servers:

* One to hold all the objects.
* Two replicas to guard against data-loss.

At the point that your first server starts to become full though there will need to be some thought put into how you deploy things.

Configuring this deployment just requires you populate your list of servers in the configuration file `/etc/sos.conf`, or `~/.sos.conf`:

     # Comments are fine
     http://node1.example.com:1234
     http://node2.example.com:1234
     http://node3.example.com:1234
     http://node4.example.com:1234

**NOTE** Don't forget to schedule the `bin/replicate` script to ensure that you do indeed have replicas of your content!



## When one store isn't enough

When you have sufficiently many objects that they will no longer fit upon a single blob-server you can add more `blob-server` nodes:

* The original store, along with a pair of mirrors.
   * `blob-server1`
   * `blob-server2`
   * `blob-server3`
* A second store, along with a pair of mirrors.
   * `blob-server4`
   * `blob-server5`
   * `blob-server6`

The naive implementation of the `api-server` means that operations would work as you expect, but operations would be slower and would take longer:

* Any upload operation would be tried upon each server in turn.
   * If it fails it will be repeated on the next server.
* A download will be tried against each server in turn.
   * When a download operation succeeds the data is returned to the caller.

For example an upload operation would be tried upon each of these nodes:

* `blob-server1`
   * This would fail; the server is full.
* `blob-server2`
   * Because this is a replica-node it has the same data as `blob-server1`, so it is full, and the upload would fail.
* `blob-server3`
   * Because this is a replica-node it has the same data as `blob-server1`, so it is full, and the upload would fail.
* `blob-server4`
   * This is the new node, designed to add capacity.
   * It is not full, so the upload succeeds.

The case of downloading an object is similar, although statistically because the first three blob-servers are full it is more likely that they hold the content you're requesting.  Eventually though there will come a point in time where you're trying to download something from the second set and it takes "too long", because three irrelevant servers have been tried in turn.

There are two solutions that can be employed here:

* We can realize download operations are more common, and maintain an index:
    * Object `foo` lives on `blob-server1`
    * Object `bar` lives on `blob-server4`
    * This doesn't solve the problem of uploads being tried too many times, but certainly speeds up downloads.
* We can realize that the first three servers belong to a logical group.
    * So rather than iterating over blob-servers we iterate over groups.

Of the two solutions the second is cleaner, because it doesn't rely upon maintaining state, and it solves the problem of both uploads and downloads.

To define your relationships you merely need to update your configuration file `/etc/sos.conf` (or `~/.sos.conf`) to list "groups" of hosts.  When that is done storage and retrieval will operate in terms of groups rather than in terms of hosts.

In our example above we defined two groups:

* `blob-server1` - First server.
   * `blob-server2` - Replicated copy of the previous server.
* `blob-server3` - Second server.
   * `blob-server4` - Replicated copy of the previous server.

We'd define that by updating our configuration file to read:

     [1]
     -: http://blob-server1.example.com:1234
     -: http://blob-server2.example.com:1234

     [2]
     -: http://blob-server3.example.com:1234
     -: http://blob-server4.example.com:1234

With this in place an upload operation will try the first server from the first group, then the first server from the second group.

This allows efficient scaling, since the potential number of attempts is bounded by the number of _groups_, and not the number of _servers_.


## Real World Usage

In my personal deployment I have five sets of three servers, hosting in excess of 5 million objects.  Things work well.

(I have so many servers because the disk-space allocated to each is pretty small.)
