package main

import (
	ui "github.com/gizak/termui/v3"
	gui2 "go-burp/internal/pkg/gui"
	"log"
)

func main() {
	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer ui.Close()

	gui := gui2.NewUserInterface()
	gui.RunEventLoop()
}
