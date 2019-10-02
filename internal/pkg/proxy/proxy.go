package proxy

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"net/http/httputil"
)

type Config struct {
	Port int
}

func copyHeader(dst, src http.Header) {
	for k, vv := range src {
		for _, v := range vv {
			dst.Add(k, v)
		}
	}
}

type Proxy struct {
	Config *Config

	Server *http.Server
	Client *http.Client

	Log *logrus.Logger
}

func NewProxy(config *Config) (p *Proxy) {
	p = &Proxy{Config: config}
	p.Server = &http.Server{
		Addr:    ":8080",
		Handler: http.HandlerFunc(p.handleHTTP),
		// Disable HTTP/2.
		TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler)),
	}

	p.Client = &http.Client{
		Transport:     &http.Transport{},
		CheckRedirect: nil,
		Timeout:       0,
	}

	p.Log = logrus.StandardLogger()
	p.Log.Log(logrus.InfoLevel, "hello world")
	return p
}

func (p *Proxy) handleHTTP(w http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()
	dump, err := httputil.DumpRequest(req, true)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	r := bufio.NewReader(bytes.NewReader(dump))
	new_req, err := http.ReadRequest(r)
	new_req.RequestURI = ""
	defer new_req.Body.Close()

	resp, err := p.Client.Do(new_req)
	defer new_req.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
	defer resp.Body.Close()
	copyHeader(w.Header(), resp.Header)
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}
