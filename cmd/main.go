package main

import (
	"go-burp/internal/pkg/proxy"
)

func main() {
	service := proxy.NewProxy(nil)
	panic(service.Server.ListenAndServe())
}
