// Proxy over HTTP. 2020 Arctic Kona. No rights reserved.
package main
import "io"
import "net/http"
import "net"

// Implements a forward TCP proxy
func proxy_tcpstream( response http.ResponseWriter , request * http.Request ) {
	if ! valid_host( request ) || ! valid_port( request ) {
		http.Error( response , "" , 400 )
		return
	}

	tomuck , err := net.Dial( "tcp" , request.Header[ "Proxy-Host" ][ 0 ] + ":" + request.Header[ "Proxy-Port" ][ 0 ] )
	if err != nil {
		http.Error( response , "" , 503 )
		return
	}

	frommuck , _ , err := response.( http.Hijacker ).Hijack( )
	if err != nil {
		http.Error( response , "" , 500 )
		return
	}

	frommuck.Write( []byte( request.Proto + " 101 Public Access Proxy Upgrade OK\r\nServer: " + proxy_name + "\r\nConnection: Upgrade\r\nUpgrade: TCPStream\r\n\r\n" ) )

	// FIXME: Caveat: Does not shuttle TCP flags like URG and PSH which may be required for some protocols
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

