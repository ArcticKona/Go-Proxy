#!/bin/bash
rm client.go
mv _server.go server.go
exec go run server.go --server=:$PORT ssh://localhost:22 pptp://localhost:1723
exit $?
