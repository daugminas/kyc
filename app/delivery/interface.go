package delivery

import (
	"github.com/daugminas/kyc/app/adapter"
	"github.com/labstack/echo/v4"
)

type server struct {
	e *echo.Echo
	a adapter.UserAdapter
}

func NewServer(e *echo.Echo, a adapter.UserAdapter) *server {
	if e == nil {
		return nil
	}
	return &server{e: e, a: a}
}

func (s *server) Start(uri string) {
	s.e.Logger.Fatal(s.e.Start(uri))
}

type Server interface {
	RegisterRouter()
	Start(uri string)
}
