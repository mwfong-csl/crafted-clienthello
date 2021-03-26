# CVE-2021-3449 PoC exploit

Usage: `go run . -host hostname:port`

This program implements a proof-of-concept exploit of CVE-2021-3449
affecting OpenSSL servers pre-1.1.1k if TLSv1.2 secure renegotiation is accepted.

It connects to a TLSv1.2 server and immediately initiates an RFC 5746 "secure renegotiation".
The attack involves a maliciously-crafted `ClientHello` that causes the server to crash
by causing a NULL pointer dereference (Denial-of-Service).

The exploit code was injected into a bundled version of the Go 1.14.15 `encoding/tls` package.

```
// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
```

– terorie
