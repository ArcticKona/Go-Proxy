Proxy over HTTP (server)
========================
What is?
--------
Proxy over HTTP allows proxying TCP and UDP connections over HTTP.

But why?
--------
HTTP CONNECTs cannot be pushed though reverse proxies, CDNs, and other infrastructure reliability.
This can be very important in censorship circumvention and other uses.

Proxy over HTTP can, while still retaining excellent performance.

TCPStream
---------
This opens a TCP tunnel to a remote server, just like the HTTP CONNECT method.

Example request to google.com via akona.me, and response.
```
GET / HTTP/1.1
Host: akona.me
Connection: upgrade
Upgrade: TCPStream
Proxy-Host: google.com
Proxy-Port: 443

HTTP/1.1 101 Switching Protocol
Server: PA-proxy/0.1
Connection: upgrade
Upgrade: TCPStream

```

Caveats:
*	TCP flags and pointers like URG and PSH are not proxied. This can be problematic for protocols like FTP


UDPStream
---------
todo: not implemented

