#!/usr/bin/perl -Iblib/lib -I../blib/lib/ -Ilib/ -I../lib/

use strict;
use warnings;

use Test::More;

BEGIN
{
    my $str = "use Test::Strict;";

    ## no critic (Eval)
    eval($str);
    ## use critic

    plan skip_all => "Skipping as Test::Strict isn't installed"
      if ($@);
}

all_perl_files_ok();
