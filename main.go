package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/venture-technology/vtx-driver/config"
	"github.com/venture-technology/vtx-driver/internal/controller"
	"github.com/venture-technology/vtx-driver/internal/repository"
	"github.com/venture-technology/vtx-driver/internal/service"

	_ "github.com/lib/pq"
)

func main() {

	config, err := config.Load("config/config.yaml")
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	db, err := sql.Open("postgres", newPostgres(config.Database))
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	err = migrate(db, config.Database.Schema)
	if err != nil {
		log.Fatalf("failed to execute migrations: %v", err)
	}

	router := gin.Default()

	driverRepository := repository.NewDriverRepository(db)
	driverService := service.NewDriverService(driverRepository)
	driverController := controller.NewDriverController(driverService)

	driverController.RegisterRoutes(router)

	fmt.Println(driverController)
	router.Run(fmt.Sprintf(":%d", config.Server.Port))
}

func newPostgres(dbconfig config.Database) string {
	return "user=" + dbconfig.User +
		" password=" + dbconfig.Password +
		" dbname=" + dbconfig.Name +
		" host=" + dbconfig.Host +
		" port=" + dbconfig.Port +
		" sslmode=disable"
}

func migrate(db *sql.DB, filepath string) error {
	schema, err := os.ReadFile(filepath)
	if err != nil {
		return err
	}

	_, err = db.Exec(string(schema))
	if err != nil {
		return err
	}

	return nil
}
