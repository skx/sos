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
