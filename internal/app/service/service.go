package service

import (
	ui "github.com/gizak/termui/v3"
	"go-burp/internal/pkg/gui"
	"go-burp/internal/pkg/proxy"
	"go-burp/internal/pkg/request"
	"os"
	"strings"
)

type Service struct {
	Ui          *gui.UserInterface
	Proxy       *proxy.Proxy
	RequestList map[int64]request.Message
}

func NewService() *Service {
	s := &Service{
		Ui:          gui.NewUserInterface(),
		Proxy:       proxy.NewProxy(nil),
		RequestList: make(map[int64]request.Message),
	}
	return s
}

func (s *Service) Run() {
	go func() {
		s.Ui.RunEventLoop()
		os.Exit(0)
	}()
	go s.Proxy.Server.ListenAndServe()

	for {
		select {
		case e := <-s.Proxy.EventBus:
			s.RequestList[e.Id] = e
			s.Ui.EventBus <- e
		case id := <-s.Ui.SelectedReqChan:
			req := string(s.RequestList[id].RequestDump)
			resp := string(s.RequestList[id].ResponseDump)
			s.Ui.RequestParagraph.Rows = strings.Split(req, "\n")
			ui.Render(s.Ui.RequestParagraph)
			s.Ui.ResponseParagraph.Rows = strings.Split(resp, "\n")
			ui.Render(s.Ui.ResponseParagraph)
		}
	}
}
