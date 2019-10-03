package main

import (
	manager "go-burp/internal/app/service"
)

func main() {
	service := manager.NewService()
	service.Run()
}
