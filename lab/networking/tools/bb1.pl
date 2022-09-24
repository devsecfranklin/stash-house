#!/usr/bin/perl
#
# franklin@bitsmasher.net
#
# 1/1/2011

eval "use Expect; 1" or die "Expect module not installed.\n";
$timeout = 10;

# go to enable mode, and then conf t
$exp = Expect->spawn("telnet labgate");
# prevent the program's output from being shown on our STDOUT
#$command->log_stdout(0);
unless($exp->expect(10, '-re','.*Password:\s*'))
{
        if( $exp->before())
        {
                print $exp->before();
                print "\n";
        }
        else
        {
                print "Timed out\n";
        }
        next;
}
print "Sending password to the host\n";
$exp->send("boots\r");
print "\n";
$exp->expect($timeout,'-re', "labgate>");
$exp->send("bb1\n");
$exp->expect($timeout, '-re', 'net/ccie\s*');
$exp->send("\n\r");
$exp->expect($timeout, '-re', '>\s*');
$exp->send("enable\r");
$exp->expect($timeout, '-re', '#\s*');
$exp->interact();

