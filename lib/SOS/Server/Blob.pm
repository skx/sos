package SOS::Server::Blob;

use File::Basename qw! basename !;

use Dancer;


#
#  Store a new blob, with the given ID.
#
post '/blob/:id' => sub {
    my $id   = params->{ id };
    my $data = request()->body();

    #
    #  Write the data to our store, based on the ID.
    #
    my $path = $ENV{ 'STORAGE' } . "/$id";

    #
    #  Ensure we catch all errors.
    #
    eval {
        open( my $handle, ">", $path );
        print $handle $data;
        close($handle);
    };


    #
    #  If there were errors then we should report them.
    #
    if ( $@ )
    {
        my %res;
        $res{'status'} = "FAILED";
        $res{'error'} = $@;
        to_json( \%res );
    }
    else
    {
        #
        #  Otherwise all is OK - but return some data just to be sure.
        #
        my ($dev,$ino,$mode,$nlink,$uid,$gid,$rdev,$size,
            $atime,$mtime,$ctime,$blksize,$blocks);

           ($dev,$ino,$mode,$nlink,$uid,$gid,$rdev,$size,
            $atime,$mtime,$ctime,$blksize,$blocks)
          = stat($path);

        my %res;
        $res{'status'} = "OK";
        $res{'id'} = $id;
        $res{'size'} = $size;
        to_json( \%res );
    }
};


#
#  Retrieve a previously uploaded blob, via the ID.
#
get '/blob/:id' => sub {
    my $id   = params->{ id };
    my $path = $ENV{ 'STORAGE' } . "/$id";

    if ( -e $path )
    {
        send_file( $path, system_path => 1 );
    }
    else
    {
        status 'not_found';
    }
};



#
#  Get the list of blob IDs we know about.
#
get '/blob/?' => sub {
    my @ids;

    foreach my $file ( sort( glob( $ENV{ 'STORAGE' } . "/*" ) ) )
    {
        push( @ids, File::Basename::basename($file) );
    }

    require JSON;
    to_json( \@ids );
};



1;
