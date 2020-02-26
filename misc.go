// Proxy over HTTP. 2020 Arctic Kona. No rights reserved.
package main
//import "strings"
import "github.com/asaskevich/govalidator"
import "net/http"

// Checks that the request object has a valid "Proxy-Host" header
func valid_host( request * http.Request ) bool {
	if len( request.Header[ "Proxy-Host" ] ) == 0 {
		return false }

	return govalidator.IsDNSName( request.Header[ "Proxy-Host" ][ 0 ] ) || govalidator.IsIP( request.Header[ "Proxy-Host" ][ 0 ] )
}

// Like above
func valid_port( request * http.Request ) bool {
	if len( request.Header[ "Proxy-Port" ] ) == 0 {
		return false }

	return govalidator.IsPort( request.Header[ "Proxy-Port" ][ 0 ] )
}


