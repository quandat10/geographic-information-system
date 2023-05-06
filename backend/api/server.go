package api

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"github.com/rs/zerolog/log"
)

type Server struct {
	store  neo4j.Driver
	router *echo.Echo
}

func NewServer(store neo4j.Driver, router *echo.Echo) (*Server, error) {
	server := &Server{
		store:  store,
		router: router,
	}
	server.setupRouter()
	return server, nil
}

func (s *Server) setupRouter() {
	s.router.Use(middleware.Logger())
	s.router.Use(middleware.Recover())
	s.router.Use(middleware.Gzip())

	r := s.router.Group("/api")
	r.POST("/user", s.NewUser)
	r.PATCH("/user", s.UpdateUser)
}

func (s *Server) StartServer(address string) error {
	err := s.store.VerifyConnectivity()
	if err != nil {
		log.Fatal().Err(err).Msg(err.Error())
	}

	return s.router.Start(address)
}
