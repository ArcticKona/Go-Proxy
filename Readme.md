DOESNT WORK YET

Forwards any TCP protocol over HTTP.

HTTP CONNECTs cannot be pushed though existing servers, CDNs, and other infrastructure reliability.

Caveats:
*	TCP flags and pointers like URG and PSH are not proxied. This can be problematic for protocols like FTP
