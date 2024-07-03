package controller

import (
	"log"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/venture-technology/vtx-driver/internal/service"
	"github.com/venture-technology/vtx-driver/models"
)

type ClaimsDriver struct {
	CNH string `json:"cnh"`
	jwt.StandardClaims
}

type DriverController struct {
	driverservice *service.DriverService
}

func NewDriverController(driverservice *service.DriverService) *DriverController {
	return &DriverController{
		driverservice: driverservice,
	}
}

func (ct *DriverController) RegisterRoutes(router *gin.Engine) {

	// conf := config.Get()

	// driverMiddleware := func(c *gin.Context) {

	// 	secret := []byte(conf.Server.Secret)

	// 	tokenString, err := c.Cookie("token")
	// 	if err != nil {
	// 		c.JSON(http.StatusUnauthorized, gin.H{"error": "Sem cookie de sessão"})
	// 		c.Abort()
	// 		return
	// 	}

	// 	if tokenString == "" {
	// 		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token não fornecido"})
	// 		c.Abort()
	// 		return
	// 	}

	// 	token, err := jwt.ParseWithClaims(tokenString, &ClaimsDriver{}, func(token *jwt.Token) (interface{}, error) {
	// 		return secret, nil
	// 	})

	// 	if err != nil {
	// 		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token inválido"})
	// 		c.Abort()
	// 		return
	// 	}

	// 	claims, ok := token.Claims.(*ClaimsDriver)
	// 	if !ok || !token.Valid {
	// 		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token inválido"})
	// 		c.Abort()
	// 		return
	// 	}

	// 	c.Set("cnh", claims.CNH)
	// 	c.Set("isAuthenticated", true)
	// 	c.Next()

	// }

	api := router.Group("vtx-driver/api/v1")

	api.GET("/ping", ct.Ping) // pingar rota
	api.POST("/driver", ct.CreateDriver)

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

	err := ct.driverservice.CreateDriver(c, &input)

	if err != nil {
		log.Printf("error to create driver: %s", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "an error occurred when creating driver"})
		return
	}

	log.Print("driver create was successful")

	c.JSON(http.StatusCreated, input)

}
