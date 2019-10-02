package gui

import (
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

type UserInterface struct {
	List              *widgets.List
	RequestParagraph  *widgets.Paragraph
	ResponseParagraph *widgets.Paragraph
	EventBus          chan interface{}
	Grid              *ui.Grid
}

func NewUserInterface() *UserInterface {
	gui := &UserInterface{
		List:              widgets.NewList(),
		RequestParagraph:  widgets.NewParagraph(),
		ResponseParagraph: widgets.NewParagraph(),
		EventBus:          make(chan interface{}, 128),
		Grid:              ui.NewGrid(),
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
		e := <-uiEvents
		switch e.ID {
		case "q", "<C-c>":
			return
		}
	}
}
