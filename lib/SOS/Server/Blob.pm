
=head1 NAME

SOS::Server::Blob - Blob-Server implementation.

=head1 DESCRIPTION

This module is the backbone of the SOS system, as it provides
the internal API upon which the public store/download server
is built.

This service exports a simple API comprising of three end-points

=over 8

=item GET /blob/:id

Retrieve the data with the specified ID.

=item GET /blob/

Retrieve a list of ID

=item POST /blob/:id

Store some data with the given ID.

=cut

=back

=cut

=head1 SYNOPSIS

=for example begin

    use SOS::Server::Blob;

    use Dancer;

    dance;

=for example end

=cut

=head1 METHODS

Now follows documentation on the available methods.

=cut


package SOS::Server::Blob;

use File::Basename qw! basename !;

use Dancer;



#
#  Is this server still alive?
#
get '/alive' => sub {
    return "alive";
};


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
    if ($@)
    {
        my %res;
        $res{ 'status' } = "FAILED";
        $res{ 'error' }  = $@;
        to_json( \%res );
    }
    else
    {
        #
        #  Otherwise all is OK - but return some data just to be sure.
        #
        my ( $dev,   $ino,     $mode, $nlink, $uid,
             $gid,   $rdev,    $size, $atime, $mtime,
             $ctime, $blksize, $blocks
           );

        (  $dev,  $ino,   $mode,  $nlink, $uid,     $gid, $rdev,
           $size, $atime, $mtime, $ctime, $blksize, $blocks )
          = stat($path);

        my %res;
        $res{ 'status' } = "OK";
        $res{ 'id' }     = $id;
        $res{ 'size' }   = $size;
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
