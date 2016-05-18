
use strict;
use warnings;

use 5.008;

use ExtUtils::MakeMaker 6.30;

use File::ShareDir::Install;



my %WriteMakefileArgs = (
    NAME         => 'SOS',
    EXE_FILES    => ['bin/blob-server','bin/download-server','bin/upload-server','bin/replicate'],
    VERSION_FROM => 'lib/SOS/Util.pm',
    PREREQ_PM    => {

        'Digest::SHA'               => 0,
        'File::Basename'             => 0,
                     'Getopt::Long'               => 0,
                     'Dancer' => 0,
        'Pod::Usage'                 => 0,

    },

    TEST_REQUIRES => {
        'Test::More'          => 0,
        'Test::NoTabs'        => 0,
        'Test::Pod'           => 0,
        'Test::Pod::Coverage' => 0,
        'Test::Strict'        => 0,
    },


    BUILD_REQUIRES => { 'Test::More'              => 0,
                        'Test::Pod'               => 0,
                        'Test::Strict'            => 0,
                        'Test::NoTabs'            => 0,
                      },

    ABSTRACT => 'A simple replicating object-storage system.',
    AUTHOR   => 'Steve Kemp <steve@steve.org.uk>',
    LICENSE  => "perl",

    MIN_PERL_VERSION => '5.008',

    META_MERGE => {
                   resources => {
                       license    => 'http://dev.perl.org/licenses/',
                       homepage   => 'https://github.com/skx/sos/',
                       bugtracker => 'https://github.com/skx/sos/issues',
                       repository => 'https://github.com/skx/sos.git',
                       GitMirror  => 'http://git.steve.org.uk/cpan/sos',
                   },
                  },
);



unless ( eval {ExtUtils::MakeMaker->VERSION(6.56)} )
{
    my $br = delete $WriteMakefileArgs{ BUILD_REQUIRES };
    my $pp = $WriteMakefileArgs{ PREREQ_PM };
    for my $mod ( keys %$br )
    {
        if ( exists $pp->{ $mod } )
        {
            $pp->{ $mod } = $br->{ $mod } if $br->{ $mod } > $pp->{ $mod };
        }
        else
        {
            $pp->{ $mod } = $br->{ $mod };
        }
    }
}

delete $WriteMakefileArgs{ CONFIGURE_REQUIRES }
  unless eval {ExtUtils::MakeMaker->VERSION(6.52)};

WriteMakefile(%WriteMakefileArgs);



package MY;
use File::ShareDir::Install 'postamble';
