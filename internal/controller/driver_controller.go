package controller

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/venture-technology/vtx-driver/internal/exceptions"
	"github.com/venture-technology/vtx-driver/internal/middleware"
	"github.com/venture-technology/vtx-driver/internal/service"
	"github.com/venture-technology/vtx-driver/models"
	"github.com/venture-technology/vtx-driver/utils"
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

	api := router.Group("api/v1")

	api.GET("/ping", ct.Ping) // pingar rota
	api.POST("/driver", ct.CreateDriver)
	api.GET("/driver/:cnh")
	api.PATCH("/driver", middleware.DriverMiddleware())
	api.DELETE("/driver", middleware.DriverMiddleware())
	api.POST("/login/driver")

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
		c.JSON(http.StatusBadRequest, exceptions.InvalidBodyContentResponseError(err))
		return
	}

	urlQrCode, err := ct.driverservice.CreateAndSaveQrCode(c, input.CNH)

	if err != nil {
		log.Printf("error to create QrCode: %s", err.Error())
		c.JSON(http.StatusBadRequest, exceptions.InternalServerResponseError(err, "an error occured qwhen creating QrCode"))
		return
	}

	input.QrCode = urlQrCode

	err = ct.driverservice.CreateDriver(c, &input)

	if err != nil {
		log.Printf("error to create driver: %s", err.Error())
		c.JSON(http.StatusInternalServerError, exceptions.InternalServerResponseError(err, "an error occured qwhen creating driver"))
		return
	}

	log.Print("driver create was successful")

	c.JSON(http.StatusCreated, input)

}

func (ct *DriverController) GetDriver(c *gin.Context) {

	cnh := c.Param("cnh")

	driver, err := ct.driverservice.GetDriver(c, &cnh)

	if err != nil {
		log.Printf("error while found driver: %s", err.Error())
		c.JSON(http.StatusBadRequest, exceptions.InternalServerResponseError(err, "driver not found"))
		return
	}

	c.JSON(http.StatusOK, driver)

}

func (ct *DriverController) UpdateDriver(c *gin.Context) {

	cnhInteface, err := ct.driverservice.ParserJwtDriver(c)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "cnh of cookie don't found"})
		return
	}

	cnh, err := utils.InterfaceToString(cnhInteface)

	log.Printf("trying change your infos --> %v", cnh)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "the value isn't string"})
		return
	}

	var input models.Driver

	if err := c.BindJSON(&input); err != nil {
		log.Printf("error to parsed body: %s", err.Error())
		c.JSON(http.StatusBadRequest, exceptions.InvalidBodyContentResponseError(err))
		return
	}

	input.CNH = *cnh

	err = ct.driverservice.UpdateDriver(c, &input)

	if err != nil {
		c.JSON(http.StatusBadRequest, exceptions.InternalServerResponseError(err, "internal server error at update"))
		return
	}

	log.Print("infos updated")

	c.JSON(http.StatusOK, gin.H{"message": "updated w successfully"})

}

func (ct *DriverController) DeleteDriver(c *gin.Context) {

	cnhInteface, err := ct.driverservice.ParserJwtDriver(c)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "cnh of cookie don't found"})
		return
	}

	cnh, err := utils.InterfaceToString(cnhInteface)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "the value isn't string"})
		return
	}

	err = ct.driverservice.DeleteDriver(c, cnh)

	if err != nil {
		log.Printf("error whiling deleted school: %s", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "error to deleted school"})
		return
	}

	c.SetCookie("token", "", -1, "/", c.Request.Host, false, true)

	log.Printf("deleted your account --> %v", cnh)

	c.JSON(http.StatusOK, gin.H{"message": "driver deleted w successfully"})

}

func (ct *DriverController) AuthDriver(c *gin.Context) {

	var input models.Driver

	log.Printf("doing login --> %s", input.Email)

	if err := c.BindJSON(&input); err != nil {
		log.Printf("error to parsed body: %s", err.Error())
		c.JSON(http.StatusBadRequest, exceptions.InvalidBodyContentResponseError(err))
		return
	}

	driver, err := ct.driverservice.AuthDriver(c, &input)

	if err != nil {
		log.Printf("wrong email or password: %s", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "wrong email or password"})
		return
	}

	jwt, err := ct.driverservice.CreateTokenJWTDriver(c, driver)

	log.Printf("token returned --> %v", jwt)

	if err != nil {
		log.Printf("error to create jwt token: %s", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "error to create jwt token"})
		return
	}

	c.SetCookie("token", jwt, 3600, "/", c.Request.Host, false, true)

	c.JSON(http.StatusAccepted, gin.H{
		"driver": driver,
		"token":  jwt,
	})

}
