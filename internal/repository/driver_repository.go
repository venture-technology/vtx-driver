package repository

import (
	"context"
	"database/sql"

	"github.com/venture-technology/vtx-driver/models"
)

type IDriverRepository interface {
	CreateDriver(ctx context.Context, driver *models.Driver) error
	DeleteDriver(ctx context.Context, cnh *string) error
	UpdateDriver(ctx context.Context, driver *models.Driver) error
	GetDriver(ctx context.Context, cnh *string) (*models.Driver, error)
	FindAllDrivers(ctx context.Context) ([]models.Driver, error)
	AuthDriver(ctx context.Context, driver *models.Driver) (*models.Driver, error)
}

type DriverRepository struct {
	db *sql.DB
}

func NewDriverRepository(conn *sql.DB) *DriverRepository {
	return &DriverRepository{
		db: conn,
	}
}

func (d *DriverRepository) CreateDriver(ctx context.Context, driver *models.Driver) error {
	sqlQuery := `INSERT INTO drivers (qrcode, name, email, password, cpf, cnh, street, number, zip, complement) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`
	_, err := d.db.Exec(sqlQuery, driver.QrCode, driver.Name, driver.Email, driver.Password, driver.CPF, driver.CNH, driver.Street, driver.Number, driver.ZIP, driver.Complement)
	return err
}

func (d *DriverRepository) DeleteDriver(ctx context.Context, cnh *string) error {
	return nil
}

func (d *DriverRepository) UpdateDriver(ctx context.Context, driver *models.Driver) error {
	return nil
}

func (d *DriverRepository) GetDriver(ctx context.Context, cnh *string) (*models.Driver, error) {
	return nil, nil
}

func (d *DriverRepository) FindAllDrivers(ctx context.Context) ([]models.Driver, error) {
	return nil, nil
}

func (d *DriverRepository) AuthDriver(ctx context.Context, driver *models.Driver) (*models.Driver, error) {
	return nil, nil
}
