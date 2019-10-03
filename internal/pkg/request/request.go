package request

import "net/http"

type Message struct {
	Request      *http.Request
	Response     *http.Response
	ListRepr     string
	Id           int64
	Host         string
	RequestDump  []byte
	ResponseDump []byte
}

type ResponseMessage struct {
	Response *http.Response
	Dump     []byte
}
