package main
import "akona.me/http-upgrade-proxy/lib"
import "flag"
import "fmt"
import "log"
import "io"
import "net"
import "net/url"
import "os"
import "strings"
import "time"

const help = "Help not yet ready.\r\n"
var server = ""

// Main loop
func main( ) {
	// Gets arguments
	flag.Usage = func( ) { fmt.Printf( help ) }
	flag.StringVar( & server , "server" , server , "" )
	flag.Parse( )
	if server == "" {
		fmt.Fprintf( os.Stderr , help )
		os.Exit( 30 ) }

	// Create object
	url , err := url.Parse( server )
	if err != nil { 
		log.Fatalf( "%w" , err ) }
	proxy , err := upgradeclient.New( url , & upgradeclient.Config{ } )
	if err != nil { 
		log.Fatalf( "%w" , err ) }

	// Parse targets
	for _ , argument := range flag.Args( ) {
		argument := strings.Split( argument , ":" )
		if len( argument ) != 3 {
			fmt.Fprintf( os.Stderr , help )
			os.Exit( 30 ) }
		if argument[ 1 ][ 0 : 2 ] == "//" {
			argument[ 1 ] = argument[ 1 ][ 2 : ] } 

		// Serve
		go func( ){
			listen , err := net.Listen( "tcp" , argument[ 1 ] + ":" + argument[ 2 ] )
			if err != nil { 
				log.Fatalf( "%w" , err ) }
			fmt.Printf( "Listening on %s\r\n" , argument[ 1 ] + ":" + argument[ 2 ] )

			for {
				// Recive connection
				frommuck , err := listen.Accept( )
				if err != nil {
					log.Printf( "%w" , err )
					continue }

				// Create stream connection
				tomuck , err := proxy.New( argument[ 0 ] )
				if err != nil { 
					log.Printf( "%w" , err )
					continue }

				// FIXME: Caveat: Does not shuttle TCP flags like URG and PSH which may be required for some protocols
				go func( ){
					io.Copy( tomuck , frommuck )
					tomuck.Close( )
				}( )
				go func( ) {
					io.Copy( frommuck , tomuck )
					frommuck.Close( )
				}( )

			}
		}( )
	}

	for ;; {
		time.Sleep( 100000 * time.Millisecond )
	}
}

