package controller

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/venture-technology/vtx-driver/internal/middleware"
	"github.com/venture-technology/vtx-driver/internal/service"
	"github.com/venture-technology/vtx-driver/models"
)

type DriverController struct {
	driverservice *service.DriverService
}

func NewDriverController(driverservice *service.DriverService) *DriverController {
	return &DriverController{
		driverservice: driverservice,
	}
}

func (ct *DriverController) RegisterRoutes(router *gin.Engine) {

	api := router.Group("vtx-driver/api/v1")

	api.GET("/ping", ct.Ping) // pingar rota
	api.POST("/driver", ct.CreateDriver)
	api.GET("/driver/:cnh")
	api.PATCH("/driver/:cnh", middleware.DriverMiddleware())
	api.DELETE("/driver/:cnh", middleware.DriverMiddleware())

}

func (ct *DriverController) Ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"ping": "pong",
	})
}

func (ct *DriverController) CreateDriver(c *gin.Context) {

	var input models.Driver

	if err := c.BindJSON(&input); err != nil {
		log.Printf("error to parsed body: %s", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body content"})
		return
	}

	urlQrCode, err := ct.driverservice.CreateAndSaveQrCode(c, input.CNH)

	if err != nil {
		log.Printf("error to create QrCode: %s", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "an error occured qwhen creating QrCode"})
		return
	}

	input.QrCode = urlQrCode

	err = ct.driverservice.CreateDriver(c, &input)

	if err != nil {
		log.Printf("error to create driver: %s", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "an error occurred when creating driver"})
		return
	}

	log.Print("driver create was successful")

	c.JSON(http.StatusCreated, input)

}
