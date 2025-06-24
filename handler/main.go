package handler

import (
	"golang-tutorial/contract"
	"log"

	"github.com/gin-gonic/gin"
)

type controller interface {
	// getPrefix returns the route prefix that the controller will use.
	getPrefix() string

	// initService initializes the necessary services for the controller.
	// This service typically contains the business logic required by the controller.
	initService(service *contract.Service)

	// initRoute sets up the routes for the controller within the given router group.
	initRoute(app *gin.RouterGroup)
}

func New(app *gin.Engine, service *contract.Service) {
	allController := []controller{
		&perpusController{},
	}

	// do not modify the code below there
	for _, c := range allController {
		c.initService(service)
		group := app.Group(c.getPrefix())
		c.initRoute(group)
		log.Printf("initiate route %s\n", c.getPrefix())
	}
}
