package main

import (
	ui "github.com/gizak/termui/v3"
	service2 "go-burp/internal/app/service"
	"log"
)

func main() {
	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer ui.Close()

	service := service2.NewService()
	service.Run()
}
