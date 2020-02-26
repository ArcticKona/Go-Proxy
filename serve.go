// Proxy over HTTP. 2020 Arctic Kona. No rights reserved.
package main
import "fmt"
import "os"
import "strings"
import "strconv"
import "net/http"
import "flag"

const proxy_name = "PA-Proxy/0.1"
const serve_help = "./proxy [-port=PORT] [-tcpstream] [-tcpbind] [-udpbind]\r\n\tSimple TCP over HTTP proxy server. Visit https://akona.me/public-acccess/server for more information.\r\n\r\n2020 Arctic Kona. No rights reserved.\r\n"

// Default settings
var serve_port = 80
var serve_tcpstream = true
var serve_tcpbind = false
var serve_udpbind = false

// Identifies the protocol to upgrade to.
func proxy_handle( response http.ResponseWriter , request * http.Request ) {
	if ( len( request.Header[ "Connection" ] ) == 0 || len( request.Header[ "Upgrade" ] ) == 0 ) {
		http.Redirect( response , request , "https://akona.me/public-access" , 307 )
		return
	}

	// Standard calls me to try each protocol. Not sure why a client would want that.
	for _ , protocol := range strings.Split( request.Header[ "Upgrade" ][ 0 ] , "," ) {
		switch strings.TrimSpace( protocol ) {
			case "TCPStream":
				if serve_tcpstream {
					proxy_tcpstream( response , request )
					return
				}
			case "TCPBind":
				if serve_tcpbind {
					proxy_tcpbind( response , request )
					return
				}
			case "UDPBind":
				if serve_udpbind {
					proxy_udpbind( response , request )
					return
				}
		}
	}

	http.Error( response , "" , 501 )
	return
}

// Main loop. Kinda obvious.
func main( ) {
	flag.Usage = func( ) {
		fmt.Print( serve_help )
	}
	flag.IntVar( & serve_port , "port" , serve_port , "port to listen on" )
	flag.BoolVar( & serve_tcpstream , "tcpstream" , ! serve_tcpstream , "enable tcpstream protocol" )
	flag.BoolVar( & serve_tcpbind , "tcpbind" , ! serve_tcpbind , "enable tcpbind protocol" )
	flag.BoolVar( & serve_udpbind , "udpbind" , ! serve_udpbind , "enable udpbind protocol" )
	flag.Parse( )
	if serve_port < 0 || serve_port > 65535 {
		fmt.Fprintf( os.Stderr , "%d: not a port \r\n" , serve_port )
		os.Exit( 10 )
	}

	// Serve
	http.HandleFunc( "/" , proxy_handle )
	fmt.Printf( "Starting at %d \r\n" , serve_port )
	err := http.ListenAndServe( ":" + strconv.Itoa( serve_port ) , nil )
	if err != nil {
		fmt.Fprintf( os.Stderr , "%s \r\n" , err ) }

	os.Exit( 20 )
}

