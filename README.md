# CVE-2021-3449 PoC exploit

Usage: `go run -host hostname:port`

This program connects to a TLSv1.2 server.
After the initial handshake, it sends a maliciously-crafted `ClientHello`
per RFC 5746 that causes the server to crash.

The exploit code was injected into a bundled version of the Go 1.14.15 `encoding/tls` package.

```
// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
```

Vulnerable OpenSSL servers: pre-1.1.1k,
if TLSv1.2 secure renegotiation is accepted.

– terorie
