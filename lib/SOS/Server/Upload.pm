
=head1 NAME

SOS::Server::Upload - Upload-Server implementation.

=head1 DESCRIPTION

This module allows remote (unauthenticated) clients to upload new
objects to the object-store, which will be handled via a number
of L<SOS::Server::Blob> servers.

When an upload operation is carried out data submitted is sent to
B<one> blob-server, with the assumption that a replication client
will trigger this being mirrored to other host(s).

=cut

=head1 SYNOPSIS

=for example begin

    use SOS::Server::Download;

    use Dancer;

    dance;

=for example end

=cut

=head1 METHODS

Now follows documentation on the available methods.

=cut


package SOS::Server::Upload;

use Dancer;
use Digest::SHA qw( sha1_hex );

use HTTP::Request::Common qw(POST);
use SOS::Util;
use LWP::UserAgent;


set port => '9991';

my $util = SOS::Util->new();


#
#  Store a new blob.
#
post '/upload' => sub {
    my $data = request()->body();

    # Hash the data.
    my $digest = sha1_hex($data);

    # Pick a server at random and upload the data.
    my @servers = $util->servers();
    my $server  = $servers[rand @servers];

    my $uri = "$server/blob/$digest";
    my $req = HTTP::Request->new( 'POST', $uri );
    $req->content($data);

    my $lwp = LWP::UserAgent->new;
    my $res = $lwp->request($req);

    # If the data succeed then we're good.
    if ( $res->is_success() )
    {
        my %res;
        $res{ 'status' } = "OK";
        $res{ 'id' }     = $digest;
        $res{ 'size' }   = length($data);
        to_json( \%res );
    }
    else
    {
        my %res;
        $res{ 'status' } = "FAILED";
        to_json( \%res );
    }
};


1;


=head1 LICENSE

This module is free software; you can redistribute it and/or modify it
under the terms of either:

a) the GNU General Public License as published by the Free Software
Foundation; either version 2, or (at your option) any later version,
or

b) the Perl "Artistic License".

=cut

=head1 AUTHOR

Steve Kemp <steve@steve.org.uk>

=cut
