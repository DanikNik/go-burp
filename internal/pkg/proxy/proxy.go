package proxy

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"go-burp/internal/pkg/config"
	"go-burp/internal/pkg/dumper"
	"go-burp/internal/pkg/request"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func copyHeader(dst, src http.Header) {
	for k, vv := range src {
		for _, v := range vv {
			dst.Add(k, v)
		}
	}
}

type Proxy struct {
	Config   *config.Config
	EventBus chan request.Message

	Server *http.Server
	Client *http.Client

	Dumper *dumper.Dumper

	Log *log.Logger
}

func NewProxy(config *config.Config) (p *Proxy) {
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

	writer, _ := os.OpenFile("go-burp.log", os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
	p.Log = log.New(writer, "[GO BURP]", log.LstdFlags)

	p.Dumper = dumper.NewDumper(p.Log)
	p.EventBus = make(chan request.Message, 128)

	p.Log.Println("Proxy inited successfully...")
	return p
}

func (p *Proxy) handleHTTP(w http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()
	dump, err := p.Dumper.DumpRequest(req, true)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	r := bufio.NewReader(bytes.NewReader([]byte(dump.Dump)))
	newReq, err := http.ReadRequest(r)
	newReq.RequestURI = ""
	defer newReq.Body.Close()

	resp, err := p.Client.Do(newReq)
	defer newReq.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
	defer resp.Body.Close()
	copyHeader(w.Header(), resp.Header)
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)

	builder := strings.Builder{}
	builder.WriteString(strconv.FormatInt(dump.Id, 10))
	builder.WriteString(dump.Host)

	p.EventBus <- request.Message{
		Request:  req,
		Response: resp,
		ListRepr: builder.String(),
	}
}

func (p *Proxy) RepeatRequest(id int64) (resp *http.Response) {
	req := p.Dumper.GetRequest(id)
	resp, err := p.Client.Do(req)
	if err != nil {
		panic(err)
	}
	return resp
}
