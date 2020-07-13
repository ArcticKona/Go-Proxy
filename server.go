// Proxy over HTTP. 2020 Arctic Kona. Some rights reserved.
// TODO: Make use of HTTP/2
package main
import "flag"
import "fmt"
import "io"
import "net"
import "net/http"
import "os"
import "strings"

const help = "Help not yet ready.\r\n"
var target = make( map[string]string )
var listen = ":80"

// Main loop
func main( ) {
	// Gets arguments
	flag.Usage = func( ) { fmt.Printf( help ) }
	flag.StringVar( & listen , "server" , listen, "" )
	flag.Parse( )

	// Parse targets
	for _ , argument := range flag.Args( ) {
		argument := strings.Split( argument , ":" )
		if len( argument ) != 3 {
			fmt.Fprintf( os.Stderr , help )
			os.Exit( 30 ) }
		if argument[ 1 ][ 0 : 2 ] == "//" {
			argument[ 1 ] = argument[ 1 ][ 2 : ]
		} 
		target[ argument[ 0 ] ] = argument[ 1 ] + ":" + argument[ 2 ]
	}

	// Serve
	http.HandleFunc( "/" , func( response http.ResponseWriter , request * http.Request ) {

		// Client didnt request upgrade :(
		if ( len( request.Header[ "Connection" ] ) == 0 || len( request.Header[ "Upgrade" ] ) == 0 ) {
			http.Error( response , "" , 204 )
			return
		}

		for _ , header := range request.Header[ "Upgrade" ] {
			for _ , protocol := range strings.Split( header , "," ) {
				protocol = strings.TrimSpace( protocol )
				if _ , err := target[ protocol ] ; err {

					// Connects on match
					tomuck , err := net.Dial( "tcp" , target[ protocol ] )
					if err != nil {
						http.Error( response , "" , 503 )
						return }

					frommuck , _ , err := response.( http.Hijacker ).Hijack( )
					if err != nil {
						http.Error( response , "" , 500 )
						return }

					frommuck.Write( []byte( request.Proto + " 101 Upgrade OK\r\nConnection: Upgrade\r\nUpgrade: " + protocol + "\r\n\r\n" ) )

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
			}
		}

		// Oops
		http.Error( response , "" , 501 )
		return

	} )
	err := http.ListenAndServe( listen , nil )
	fmt.Fprintf( os.Stderr , "%s\r\n" , err )

	os.Exit( 20 )
}


