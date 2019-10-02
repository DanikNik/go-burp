package gui

import (
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	"go-burp/internal/pkg/request"
)

type UserInterface struct {
	List              *widgets.List
	RequestParagraph  *widgets.Paragraph
	ResponseParagraph *widgets.Paragraph
	Grid              *ui.Grid

	EventBus chan request.Message
}

func NewUserInterface() *UserInterface {
	gui := &UserInterface{
		List:              widgets.NewList(),
		RequestParagraph:  widgets.NewParagraph(),
		ResponseParagraph: widgets.NewParagraph(),
		Grid:              ui.NewGrid(),
		EventBus:          make(chan request.Message, 128),
	}

	termWidth, termHeight := ui.TerminalDimensions()
	gui.Grid.SetRect(0, 0, termWidth, termHeight)

	gui.Grid.Set(
		ui.NewCol(1.0/2, gui.List),
		ui.NewCol(1.0/2,
			ui.NewRow(1.0/2, gui.RequestParagraph),
			ui.NewRow(1.0/2, gui.ResponseParagraph),
		),
	)

	gui.List.SelectedRowStyle = ui.NewStyle(ui.ColorBlack, ui.ColorGreen)
	gui.List.WrapText = false

	gui.List.Rows = []string{
		"1",
		"1",
		"1",
		"1",
		"1",
		"1",
	}
	return gui
}

func (g *UserInterface) RunEventLoop() {
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
				//case "i":
				//	g.EventBus <- request.Message{
				//		Request:  nil,
				//		ListRepr: "ZHOPE",
				//	}
			case "<Enter>":

			}
			ui.Render(g.List)
		case reqMessage := <-g.EventBus:
			g.List.Rows = append(g.List.Rows, reqMessage.ListRepr)
			ui.Render(g.List)
		}
	}
}
