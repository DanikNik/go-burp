package main

import (
	"go-burp/internal/pkg/proxy"
)


func main(){
	service := proxy.NewProxy(nil)
	_ = service.Server.ListenAndServe()
}