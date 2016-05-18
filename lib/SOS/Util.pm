
=head1 NAME

SOS::Util - Utility object for the Simple Object Storage system.

=head1 DESCRIPTION

This module contains a simple utility method which is designed to
retrieve the list of know blob-servers, implemented via the
L<SOS::Server::Blog> module.

Servers are listed in the file C</etc/sos.conf> or C<~/.sos.conf>
file, and we use this module to return only the listed servers which
are reachable.

=cut

=head1 SYNOPSIS

=for example begin

     use SOS::Util;

     my $util = SOS::Util->new();

     foreach my $server ($util->servers() )
     {
        print "blob-server: $server\n";
     }

=for example end

=cut

=head1 METHODS

Now follows documentation on the available methods.

=cut




#
#  Standard modules which we require.
#
use strict;
use warnings;


package SOS::Util;
require LWP::UserAgent;

our $VERSION = "0.2";


=head2 new

Constructor.

=cut

sub new
{
    my ( $proto, %supplied ) = (@_);
    my $class = ref($proto) || $proto;

    my $self = {};
    bless( $self, $class );
    return $self;

}


=head2 servers

Retrieve the servers from our configuration file(s) which are actually
alive, and reachable.

=cut

sub servers
{
    my ($self) = (@_);

    if ( !scalar( $self->{ 'servers' } ) )
    {
        $self->read_servers( $ENV{ 'HOME' } . "/.sos.conf" )
          if ( $ENV{ 'HOME' } );
        $self->read_servers("/etc/sos.conf");
    }

    #
    #  We'll test each server is alive before using it.
    #
    my @result;

    my $ua = LWP::UserAgent->new;
    $ua->timeout(10);
    $ua->env_proxy;

    foreach my $s ( @{ $self->{ 'servers' } } )
    {
        my $response = $ua->get( $s . '/' );
        my $code     = $response->code();

        if ( $code && ( $code =~ /^(404|200)$/ ) )
        {
            push( @result, $s );
        }
        else
        {
            warn "Server offline $s";
        }
    }

    return (@result);

}


=head2 read_servers

Read the server-names/ports from the supplied files and return
them to the caller.

=cut

sub read_servers
{
    my ( $self, $file ) = (@_);

    return unless -e $file;

    open( my $handle, "<", $file ) or
      return;
    while ( my $line = <$handle> )
    {
        next unless $line;
        chomp($line);

        next if ( $line =~ /^\#/ );

        push( @{ $self->{ 'servers' } }, $line );
    }
    close($handle);

}


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
