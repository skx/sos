Replication
-----------

Replication is a key feature of the SoS, ensuring that uploaded objects are copied sufficiently many times that data-loss should be unlikely.

Replication is implemented in terms of groups, which are defined in the configuration file `/etc/sos.conf`, or `~/.sos.conf`.  A sample configuration would look like this:

    #
    # This configuration file defines three groups
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
