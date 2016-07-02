# netln

A program that listens on one TCP address and passes everything to
another TCP address. Also could be called a proxy or a pipe.


## Examples

To "bring out" an HTTP server from the internal network:

	$ netln :80 192.168.1.1:80

Now all connections on port 80 to this machine will be retransmitted to
the port 80 at 192.168.1.1.

To create a throttled version of a server:

	$ netln -s 33.6K :80 :8080

Now browsing at port 8080 will be the same as at the default port
except the visitor will be able to appreciate how bloated web pages
have become...
