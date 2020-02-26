// Proxy over HTTP. 2020 Arctic Kona. No rights reserved.
package main
import "net/http"

// Not implemented
func proxy_tcpbind( response http.ResponseWriter , request * http.Request ) {
	http.Error( response , "" , 501 )
	return
}
