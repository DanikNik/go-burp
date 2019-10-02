#!/usr/bin/env bash

case $(uname -s) in
Linux*) sslConfig=/etc/ssl/openssl.cnf ;;
Darwin*) sslConfig=/System/Library/OpenSSL/openssl.cnf ;;
esac
openssl req \
  -newkey rsa:2048 \
  -x509 \
  -nodes \
  -keyout server.key \
  -new \
  -out server.pem \
  -subj /CN=localhost \
  -reqexts SAN \
  -extensions SAN \
  -config <(cat $sslConfig \
    <(printf '[SAN]\nsubjectAltName=DNS:localhost')) \
  -sha256 \
  -days 3650


#GET http://www.elit-san.com/ HTTP/1.1
#Accept: text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8
#Accept-Encoding: gzip, deflate
#Accept-Language: en-US,en;q=0.5
#Connection: keep-alive
#Cookie: _ga=GA1.2.608787467.1569325588; _ym_uid=156932558813045714; _ym_d=1569325588
#Upgrade-Insecure-Requests: 1
#User-Agent: Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:69.0) Gecko/20100101 Firefox/69.0
#


#GET / HTTP/1.1
#Host: www.elit-san.com
#User-Agent: Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:69.0) Gecko/20100101 Firefox/69.0
#Accept: text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8
#Accept-Encoding: gzip, deflate
#Accept-Language: en-US,en;q=0.5
#Cache-Control: max-age=0
#Connection: keep-alive
#Cookie: _ga=GA1.2.608787467.1569325588; _ym_uid=156932558813045714; _ym_d=1569325588
#Upgrade-Insecure-Requests: 1