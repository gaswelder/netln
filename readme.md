# netln

A program that listens on one TCP address and passes everything to
another TCP address. Also could be called a proxy or a pipe.


## Usage

	$ netln [-s <max-speed>] <remote-addr> <listen-addr>

The `s` flag specifies trasfer speed limit as string in form
`<number>[<unit>]` where possible units are "K" for Kbps and "M" for
Mbps. If no unit is given, bytes per second are assumed.

The remote address goes before the listen address (analogous to the
`ln` command).


## Examples

To "bring out" an HTTP server from the internal network:

	$ netln 192.168.1.1:80 :80

Now all connections on port 80 to this machine will be retransmitted to
the port 80 at 192.168.1.1.

To create a throttled version of a server:

	$ netln -s 33.6K :80 :8080

Now browsing at port 8080 will be the same as at the default port
except the visitor will be able to appreciate how bloated web pages
have become...
