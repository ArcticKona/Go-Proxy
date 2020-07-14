package libclient
import "crypto/tls"
import "fmt"
import "io"
import "net"
import "net/url"

type Client interface {
	New( string ) net.Conn
}

type client struct {
	url.URL
	Config
}

type Config struct {
	tls.Config
}

// Create new stream
func ( self client ) New( proto string ) ( net.Conn , error ) {

	// Create connection to server
	var tomuck net.Conn
	var err error
	if self.Scheme == "http" {
		tomuck , err := net.Dial( "tcp" , self.Host )
	} else if self.Scheme == "https" {
		tomuck , err := tls.Dial( "tcp" , self.Host , &self.Config ) }
	if err != nil {
		return nil , err }

	// Make HTTP request
	tomuck.Write( []byte( fmt.Sprintf( "GET %s HTTP/1.0\r\nHost: %s\r\nConnection: upgrade\r\nUpgrade: %s\r\n\r\n" , self.Path , self.Host , proto ) ) )

	var buff = make( []byte , 0 , 12 )
	_ , err = io.ReadFull( tomuck , buff )
	if err != nil {
		return nil , err }
	if string( buff )[ 0 : 5 ] != "HTTP/" || string( buff )[ 9 : 12 ] != "101" {
		return nil , fmt.Errorf( "Server did not upgrade: %d" , buff[ 9 : 12 ] ) }

	// Okie!
	return tomuck , nil
}

// Create new proxy server
func New( server url.URL , config Config ) ( Client , error ) {
	if server.Scheme != "http" && server.Scheme != "https" {
		return nil , fmt.Errorf( "I dont know this scheme: %s" , server.Scheme ) }
	var proxy client
	proxy.URL = server
	proxy.Config = config
	return proxy , nil
}


