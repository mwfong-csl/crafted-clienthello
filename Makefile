crafted-clienthello: main.go $(wildcard tls/*.go)
	go build -o crafted-clienthello .

export UBUNTU_VERSION

.ONESHELL:
demo-openssl: crafted-clienthello
	$(MAKE) -C demo start-openssl
	sleep 2
	./crafted-clienthello -host localhost:4433
	$(MAKE) -C demo stop-openssl

.ONESHELL:
demo-apache: crafted-clienthello
	$(MAKE) -C demo start-apache
	sleep 10
	./crafted-clienthello -host localhost:443
	sleep 5
	$(MAKE) -C demo stop-apache

.ONESHELL:
demo-haproxy: crafted-clienthello
	$(MAKE) -C demo start-haproxy
	sleep 3
	./crafted-clienthello -host localhost:4433
	sleep 3
	$(MAKE) -C demo stop-haproxy

.ONESHELL:
demo-lighttpd: crafted-clienthello
	$(MAKE) -C demo start-lighttpd
	sleep 3
	./crafted-clienthello -host localhost:4433
	sleep 3
	$(MAKE) -C demo stop-lighttpd

.ONESHELL:
demo-nginx: crafted-clienthello
	$(MAKE) -C demo start-nginx
	sleep 3
	./crafted-clienthello -host localhost:4433
	sleep 3
	$(MAKE) -C demo stop-nginx

.ONESHELL:
demo-nodejs: crafted-clienthello
	$(MAKE) -C demo start-nodejs
	sleep 8
	./crafted-clienthello -host localhost:4433
	sleep 3
	$(MAKE) -C demo stop-nodejs

clean:
	rm -f crafted-clienthello
	$(MAKE) -C demo clean
