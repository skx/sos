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
    my $server = $servers[rand @servers];

    my $uri = "$server/blob/$digest";
    my $req = HTTP::Request->new( 'POST', $uri );
    $req->content($data);

    my $lwp = LWP::UserAgent->new;
    my $res = $lwp->request($req);

    # If the data succeed then we're good.
    if ( $res->is_success() )
    {
        my %res;
        $res{'status'} = "OK";
        $res{'id'} = $digest;
        $res{'size'} = length($data);
        to_json( \%res );
    }
    else
    {
        my %res;
        $res{'status'} = "FAILED";
        to_json( \%res );
    }
};


1;
