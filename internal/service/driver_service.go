package service

import (
	"context"

	"github.com/venture-technology/vtx-driver/internal/repository"
	"github.com/venture-technology/vtx-driver/models"
	"github.com/venture-technology/vtx-driver/utils"
)

type DriverService struct {
	driverrepository repository.IDriverRepository
}

func NewDriverService(driverrepository repository.IDriverRepository) *DriverService {
	return &DriverService{
		driverrepository: driverrepository,
	}
}

func (d *DriverService) CreateDriver(ctx context.Context, driver *models.Driver) error {

	driver.Password = utils.HashPassword(driver.Password)

	return d.driverrepository.CreateDriver(ctx, driver)
}

func (d *DriverService) GetDriver(ctx context.Context, cnh *string) (*models.Driver, error) {
	return d.driverrepository.GetDriver(ctx, cnh)
}

func (d *DriverService) FindAllDrivers(ctx context.Context) ([]models.Driver, error) {
	return d.driverrepository.FindAllDrivers(ctx)
}

func (d *DriverService) UpdateDriver(ctx context.Context, driver *models.Driver) error {
	return d.driverrepository.UpdateDriver(ctx, driver)
}

func (d *DriverService) DeleteDriver(ctx context.Context, cnh *string) error {
	return d.driverrepository.DeleteDriver(ctx, cnh)
}

func (d *DriverService) AuthDriver(ctx context.Context, driver *models.Driver) (*models.Driver, error) {
	return d.driverrepository.AuthDriver(ctx, driver)
}
