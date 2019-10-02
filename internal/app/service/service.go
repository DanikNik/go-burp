package service

import (
	"fmt"
	"go-burp/internal/pkg/gui"
	"go-burp/internal/pkg/proxy"
	"go-burp/internal/pkg/request"
)

type Service struct {
	Ui          *gui.UserInterface
	Proxy       *proxy.Proxy
	EventBus    chan request.Message
	RequestList map[int64]request.Message
}

func NewService() *Service {
	s := &Service{
		Ui:       gui.NewUserInterface(),
		Proxy:    proxy.NewProxy(nil),
		EventBus: make(chan request.Message, 128),
	}

}

func (s *Service) Run() {
	go s.Ui.RunEventLoop()
	go s.Proxy.Server.ListenAndServe()

	for {
		select {
		case e := <-s.Proxy.EventBus:
			fmt.Println(e)

			//case e := <-s.Ui.EventBus:
			//	fmt.Println(e)

		}
	}
}

//func (s *Service) Prepare(config *config.Config) {
//
//}
