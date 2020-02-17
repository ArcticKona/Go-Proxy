// Proxy over HTTP. 2020 Arctic Kona. No rights reserved.
package main
import "io"
import "github.com/asaskevich/govalidator"
import "net/http"
import "net"

// Implements a forward TCP proxy, the most basic subprotocol
func proxy_tcpstream( response http.ResponseWriter , request * http.Request ) {
	// Check target host and port
	if len( request.Header[ "Proxy-Host" ] ) == 0 || ! ( govalidator.IsDNSName( request.Header[ "Proxy-Host" ][ 0 ] ) || govalidator.IsIP( request.Header[ "Proxy-Host" ][ 0 ] ) ) || len( request.Header[ "Proxy-Port" ] ) == 0 || ! govalidator.IsPort( request.Header[ "Proxy-Port" ][ 0 ] ) {	// ... lol
		http.Error( response , "" , 400 )
		return
	}

	// make TCP connection
	tomuck , err := net.Dial( "tcp" , request.Header[ "Proxy-Host" ][ 0 ] + ":" + request.Header[ "Proxy-Port" ][ 0 ] )
	if err != nil {
		http.Error( response , "" , 403 )
		return
	}

	// Hijack HTTP connection
	frommuck , _ , err := response.( http.Hijacker ).Hijack( )
	if err != nil {
		http.Error( response , "" , 500 )
		return
	}

	// Success!
	frommuck.Write( []byte( request.Proto + " 101 Public Access Proxy Upgrade OK\r\nServer: PA-Proxy/0.1\r\nConnection: Upgrade\r\nUpgrade: TCPStream\r\n\r\n" ) )

	// Shuttle data
	// TODO: Caveat: Does not shuttle TCP flags like URG and PSH which may be required for some protocols
	go func( ){
		io.Copy( tomuck , frommuck )
		tomuck.Close( )
	}( )
	go func( ) {
		io.Copy( frommuck , tomuck )
		frommuck.Close( )
	}( )

	return
}
