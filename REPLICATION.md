Replication
-----------

Replication is a key feature of the SoS, ensuring that uploaded objects are copied sufficiently many times that data-loss should be unlikely.

Replication is implemented in terms of groups, which are defined in the configuration file `/etc/sos.conf`, or `~/.sos.conf`.  A sample configuration would look like this:

    #
    # This configuration file defines three groups.
    #
    # Lines prefixed with "#" are comments, and are ignored.
    #
    [1]
    -: http://node1.example.com:3000
    -: http://node1-2.example.com:3000
    -: http://node1-3.example.com:3000
    [2]
    -: http://node2.example.com:3000
    -: http://node2-2.example.com:3000
    -: http://node2-3.example.com:3000
    [3]
    -: http://node3.example.com:3000
    -: http://node3-2.example.com:3000
    -: http://node3-3.example.com:3000

Replication is carried out __exclusively__ between members of the same group, so in the first group above we have this:

* Data from `http://node1.example.com:3000` is replicated
   * To `http://node1-2.example.com:3000`.
   * To `http://node1-3.example.com:3000`.
* Data from `http://node1-2.example.com:3000` is replicated
   * To `http://node1.example.com:3000`.
   * To `http://node1-3.example.com:3000`.
* Data from `http://node1-3.example.com:3000` is replicated
   * To `http://node1.example.com:3000`.
   * To `http://node1-3.example.com:3000`.

Similarly the second-group replicate between themselves, and the third-group
replicats between the members in that group.

There is __zero__ replication between nodes in different groups.

If you wish to have two copies of all uploaded objects then you need to define each group with two members:

    [1]
    -: http://node1.example.com:3000
    -: http://mirror1.example.com:3000
    [2]
    -: http://node2.example.com:3000
    -: http://mirror2.example.com:3000
    [3]
    -: http://node3.example.com:3000
    -: http://mirror3.example.com:3000

If you wish to have five replicas of all uploaded objects the principle is the same:

    [1]
    -: http://node1.example.com:3000
    -: http://mirror1.example.com:3000
    -: http://mirror2.example.com:3000
    -: http://mirror3.example.com:3000
    -: http://mirror4.example.com:3000
    [2]
    -: http://node2.example.com:3000
    -: http://mirror5.example.com:3000
    -: http://mirror6.example.com:3000
    -: http://mirror7.example.com:3000
    -: http://mirror8.example.com:3000


Triggering Replication
----------------------

Replication is __not__ triggered automatically, although in the future that is an ideal enhancement.  To trigger replication you must run the replication sub-command manually, and regularly:

    $ sos replicate [-verbose]
