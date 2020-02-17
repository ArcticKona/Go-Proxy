// Proxy over HTTP. 2020 Arctic Kona. No rights reserved.
package main
//import "fmt"
//import "io"
//import "github.com/asaskevich/govalidator"
import "net/http"
//import "net"

// Not implemented
func proxy_udpstream( response http.ResponseWriter , request * http.Request ) {
	http.Error( response , "" , 501 )
	return
}
