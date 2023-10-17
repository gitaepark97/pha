package loader

import (
	"database/sql"

	"github.com/gitaepark/pha/controller"
	"github.com/gitaepark/pha/repository"
	"github.com/gitaepark/pha/service"
	"github.com/gitaepark/pha/util"
)

type Server struct {
	config     util.Config
	controller *controller.Controller
}

func NewServer(config util.Config, conn *sql.DB) (*Server, error) {
	repository := repository.New(conn)
	service := service.NewService(config, repository)
	controller := controller.NewController(config, service)

	server := &Server{
		config:     config,
		controller: controller,
	}

	return server, nil
}

func (server *Server) Start(address string) error {
	return server.controller.Run(address)
}
