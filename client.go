package main
import "flag"
import "fmt"
import "io"
import "net"
import "net/http"
import "os"
import "strings"

const help = "Help not yet ready.\r\n"
var server = "example.com:80"
var tls = false

// Main loop
func main( ) {
	// Gets arguments
	flag.Usage = func( ) { fmt.Printf( help ) }
	flag.StringVar( & server , "server" , server , "" )
	flag.Parse( )

	// Simple autodetect SSL
//        _ , err := tls.Dial( "tcp" , server , & tls.Config{
//                InsecureSkipVerify: true,
//        } )
//        if ( err == nil ) {
//                tls = true }

	// Parse targets
	for _ , argument := range flag.Args( ) {
		argument := strings.Split( argument , ":" )
		if len( argument ) != 3 {
			fmt.Fprintf( os.Stderr , help )
			os.Exit( 30 ) }

		// Serve
		go func( ){
			listen , err := net.Listen( "tcp" , argument[ 1 ] + ":" + argument[ 2 ] )
			if err != nil { 
				fmt.Fprintf( os.Stderr , "%w" , err ) }
			for {
				tomuck , err := listen.Accept( )
				if err != nil { 
					fmt.Fprintf( os.Stderr , "%w" , err ) }

				// Make connection to server
				request , err := http.NewRequest( "GET" , server , nil )
				request.Header.Set( "Connection" , "upgrade" )
				request.Header.Set( "Upgrade" , argument[ 0 ] )
				frommuck := & http.Transport{ }
				response , err := ( & http.Client{ Transport: frommuck } ).Do( request )
				if err != nil {
					fmt.Fprintf( os.Stderr , "%w" , err ) }
				if response.StatusCode != 101 {
					fmt.Fprintf( os.Stderr , "ERROR: Server did not upgrade: %d" , response.StatusCode ) }

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
		}( )
	}

	os.Exit( 20 )
}

