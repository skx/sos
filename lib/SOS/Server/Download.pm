
=head1 NAME

SOS::Server::Download - Download-Server implementation.

=head1 DESCRIPTION

This module allows remote (unauthenticated) clients to download
objects which are stored on one or more L<SOS::Server::Blob> servers.

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

package SOS::Server::Download;

use Dancer;
use Digest::SHA qw( sha1_hex );

use HTTP::Request::Common qw(POST);

use SOS::Util;
use LWP::UserAgent;

set port => '9992';

my $util = SOS::Util->new();


#
#  Retrieve a previously uploaded blob.
#
get '/fetch/:id' => sub {
    my $id = params->{ id };

    #
    #  Work out which server has the content.
    #
    my @servers = $util->servers();
    foreach my $server (@servers)
    {
        my $uri = "$server/blob/$id";

        my $ua       = new LWP::UserAgent;
        my $response = $ua->get($uri);
        if ( $response->is_success() )
        {
            return $response->decoded_content();
        }
    }


    #
    #  Failed to find it
    #
    status 'not_found';

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
