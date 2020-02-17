// Proxy over HTTP. 2020 Arctic Kona. No rights reserved.
package main
import "fmt"
import "os"
import "github.com/asaskevich/govalidator"
import "net/http"

const help_text = "./TCPStream PORT\r\n\tSimple TCPStream over HTTP proxy server. Visit https://akona.me/public-acccess for more information.\r\n\r\n2020 Arctic Kona. No rights reserved.\r\n"

// Identifies the subprotocol to upgrade to.
func proxy_handle( response http.ResponseWriter , request * http.Request ) {
	if ( len( request.Header[ "Upgrade" ] ) == 0 ) {
		http.Redirect( response , request , "https://akona.me/public-access" , 307 )
		return
	}

	switch request.Header[ "Upgrade" ][ 0 ] {
		case "TCPStream":
			proxy_tcpstream( response , request )
		default:
			http.Error( response , "" , 501 )
	}

	return
}

// Main loop. Kinda obvious.
func main( ) {
	if len( os.Args ) == 1 || os.Args[ 1 ] == "--help" || os.Args[ 1 ] == "-h" || os.Args[ 1 ] == "/?" {
		fmt.Print( help_text )
		return
	}
	if ! govalidator.IsPort( os.Args[ 1 ] ) {
		fmt.Fprint( os.Stderr , os.Args[ 1 ] + ": not a port number\r\n" )
		os.Exit( 1 )
	}

	http.HandleFunc( "/" , proxy_handle )
	fmt.Print( "Online at " + os.Args[ 1 ] + "\r\n" )
	err := http.ListenAndServe( ":" + os.Args[ 1 ] , nil )
	if ( err != nil ) {
		fmt.Fprint( os.Stderr , err ) }

	return
}

