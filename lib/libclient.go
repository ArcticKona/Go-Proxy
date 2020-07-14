package upgradeclient
import "bufio"
import "crypto/tls"
import "fmt"
import "io"
import "net"
import "net/url"
import "strings"
//import "os"

type Client interface {
	New( string ) ( net.Conn , error )
}

type client struct {
	* url.URL
	* tls.Config
	RedirectMax int
	redirectDepth int
}

type Config struct {
	tls.Config
	RedirectMax int
}

// Create new stream
func ( self client ) New( proto string ) ( net.Conn , error ) {
	// Properly format URL 'cause golang wouldnt do it for me :(
	if self.URL.Port( ) == "" {
		if self.Scheme == "http" {
			self.Host = self.Host + ":80"
		} else if self.Scheme == "https" {
			self.Host = self.Host + ":443" }
	}
	if self.Path == "" {
		self.Path = "/" }

	// Create connection to server
	var tomuck net.Conn
	var err error
	if self.Scheme == "http" {
		tomuck , err = net.Dial( "tcp" , self.Host )
	} else if self.Scheme == "https" {
		tomuck , err = tls.Dial( "tcp" , self.Host , self.Config ) }
	if err != nil {
		return nil , err }

	// Make HTTP request	TODO: Timeouts
	tomuck.Write( []byte( fmt.Sprintf( "GET %s HTTP/1.0\r\nHost: %s\r\nConnection: upgrade\r\nUpgrade: %s\r\n\r\n" , self.Path , self.Host , proto ) ) )
//fmt.Fprintf( os.Stderr , "GET %s HTTP/1.0\r\nHost: %s\r\nConnection: upgrade\r\nUpgrade: %s\r\n\r\n" , self.Path , self.Host , proto )
	var buff = make( []byte , 12 )
	_ , err = io.ReadFull( tomuck , buff )
	if err != nil {
		return nil , err }

	// Is it HTTP?
	if string( buff )[ 0 : 5 ] != "HTTP/" {
		return nil , fmt.Errorf( "Protocol mismatch" ) }

	// Is it redirect?
	if string( buff )[ 9 : 11 ] == "30" {
		if self.redirectDepth == self.RedirectMax {
			return nil , fmt.Errorf( "Too many redirects" ) }
		scanner := bufio.NewScanner( tomuck )
		for {
			scanner.Scan( )
			header := strings.Split( scanner.Text( ) , ":" )
			if ( len( header ) < 2 ) {
				continue }
			if strings.ToLower( header[ 0 ] ) == "location" {
				self.URL , err = url.Parse( strings.TrimSpace( strings.Join( header[ 1 : ] , ":" ) ) )
				if err != nil {
						return nil , err }
				tomuck.Close( )
				self.redirectDepth ++
				return self.New( proto )
			}
		}
		if err := scanner.Err( ) ; err != nil {
			return nil , err }
		return nil , fmt.Errorf( "Server did something I cant understand" )
	}

	// No upgrade
	if string( buff )[ 9 : 12 ] != "101" {
		return nil , fmt.Errorf( "Server did not upgrade: %s\r\n" , string( buff[ 9 : 12 ]  ) ) }

	// Ignore headers
	scanner := bufio.NewScanner( tomuck )
	for {
		scanner.Scan( )
		if strings.TrimSpace( scanner.Text( ) ) == "" {
			break }
	}
	if err := scanner.Err( ) ; err != nil {
		return nil , err }

	// Okie!
	return tomuck , nil
}

// Create new proxy server
func New( server * url.URL , config * Config ) ( Client , error ) {
	if server.Scheme != "http" && server.Scheme != "https" {
		return nil , fmt.Errorf( "I dont know this scheme: %s\r\n" , server.Scheme ) }
	var proxy client
	proxy.URL = server
	proxy.Config = & config.Config
	proxy.RedirectMax = config.RedirectMax
	if proxy.RedirectMax == 0 {
		proxy.RedirectMax = 16 }
	return proxy , nil
}


