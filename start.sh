#!/bin/bash
rm client.go
mv _server.go server.go
exec go run server.go --server=:$PORT ssh://localhost:22 pptp://localhost:1723 ssh-mobile://localhost:8022 ssh-kona://localhost:12003
exit $?
