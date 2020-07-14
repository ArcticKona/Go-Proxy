package main
import "akona.me/http-upgrade-proxy/lib"
import "crypto/tls"
import "flag"
import "fmt"
import "io"
import "net"
import "net/url"
import "os"
import "strings"
import "time"

const help = "Help not yet ready.\r\n"
var server = "example.com:80"

// Main loop
func main( ) {
	// Gets arguments
	flag.Usage = func( ) { fmt.Printf( help ) }
	flag.StringVar( & server , "server" , server , "" )
	flag.Parse( )

	// Create object
	server = libclient.New( url.Parse( server ) )

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

				// Recive connection
				frommuck , err := listen.Accept( )
				if err != nil {
					fmt.Fprintf( os.Stderr , "%w" , err ) }

				// Create stream connection
				tomuck , err := server.New( argument[ 0 ] )
				if err != nil { 
					fmt.Fprintf( os.Stderr , "%w" , err ) }

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

	for ;; {
		time.Sleep( 100000 * time.Millisecond )
	}
}

