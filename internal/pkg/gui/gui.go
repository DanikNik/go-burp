package gui

import (
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	"go-burp/internal/pkg/request"
	"log"
	"strconv"
	"strings"
)

type UserInterface struct {
	List              *widgets.List
	RequestParagraph  *widgets.List
	ResponseParagraph *widgets.List
	Grid              *ui.Grid

	EventBus        chan request.Message
	SelectedReqChan chan int64
}

func NewUserInterface() *UserInterface {
	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	gui := &UserInterface{
		List:              widgets.NewList(),
		RequestParagraph:  widgets.NewList(),
		ResponseParagraph: widgets.NewList(),
		Grid:              ui.NewGrid(),
		EventBus:          make(chan request.Message, 128),
		SelectedReqChan:   make(chan int64),
	}

	termWidth, termHeight := ui.TerminalDimensions()
	gui.Grid.SetRect(0, 0, termWidth, termHeight)

	gui.Grid.Set(
		ui.NewCol(1.0/8, gui.List),
		ui.NewCol(1-1.0/8,
			ui.NewRow(1.0/2, gui.RequestParagraph),
			ui.NewRow(1.0/2, gui.ResponseParagraph),
		),
	)

	gui.List.SelectedRowStyle = ui.NewStyle(ui.ColorBlack, ui.ColorGreen)
	gui.List.WrapText = false
	return gui
}

func (g *UserInterface) RunEventLoop() {
	defer ui.Close()
	ui.Render(g.Grid)
	uiEvents := ui.PollEvents()
	for {
		select {
		case e := <-uiEvents:
			switch e.ID {
			case "q", "<C-c>":
				return
			case "<Resize>":
				payload := e.Payload.(ui.Resize)
				g.Grid.SetRect(0, 0, payload.Width, payload.Height)
				ui.Clear()
				ui.Render(g.Grid)
			case "<Down>":
				g.List.ScrollDown()
			case "<Up>":
				g.List.ScrollUp()
			case "<Enter>":
				n := g.List.SelectedRow
				id, err := strconv.Atoi(strings.Split(g.List.Rows[n], " ")[0])
				if err != nil {
					panic(err)
				}
				g.SelectedReqChan <- int64(id)
			case "h":
				g.RequestParagraph.ScrollHalfPageUp()
			case "n":
				g.RequestParagraph.ScrollHalfPageDown()
			case "j":
				g.ResponseParagraph.ScrollHalfPageUp()
			case "m":
				g.ResponseParagraph.ScrollHalfPageDown()
			}
			ui.Render(g.List, g.RequestParagraph, g.ResponseParagraph)
		case reqMessage := <-g.EventBus:
			g.List.Rows = append(g.List.Rows, reqMessage.ListRepr)
			ui.Render(g.List)
		}
	}
}
