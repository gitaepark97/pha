package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/gitaepark/pha/service"
	"github.com/gitaepark/pha/util"
	"github.com/gitaepark/pha/util/validator"
	"github.com/rs/zerolog/log"
)

type Controller struct {
	config  util.Config
	service service.Service
	router  *gin.Engine
}

func NewController(config util.Config, service service.Service) *Controller {
	controller := &Controller{
		config:  config,
		service: service,
	}

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		err := v.RegisterValidation("phone_number", validator.ValidatePhoneNumber)
		if err != nil {
			log.Fatal().Msg("cannot create validation")
		}
		err = v.RegisterValidation("date", validator.ValidateDate)
		if err != nil {
			log.Fatal().Msg("cannot create validation")
		}
		err = v.RegisterValidation("product_size", validator.ValidateProductSize)
		if err != nil {
			log.Fatal().Msg("cannot create validation")
		}
	}

	controller.setupRouter()

	return controller
}

func (controller *Controller) setupRouter() {
	controller.router = gin.Default()

	controller.setHealthCheck()

	controller.setAuthRouter()
	controller.setProductRouter()
}

func (controller *Controller) Run(address string) error {
	return controller.router.Run(address)
}

func (controller *Controller) setHealthCheck() {
	controller.router.GET("/api/health", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"message": "OK"})
	})
}
