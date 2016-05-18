#!/usr/bin/perl -I../blib/lib/ -Iblib/lib/

use strict;
use warnings;

use Test::More tests => 4;

BEGIN
{

    #
    #  Helpers
    #
    use_ok( "SOS::Util", "Loaded module" );


    #
    #  Servers
    #
    use_ok( "SOS::Server::Blob",     "Loaded module" );
    use_ok( "SOS::Server::Download", "Loaded module" );
    use_ok( "SOS::Server::Upload",   "Loaded module" );
}
