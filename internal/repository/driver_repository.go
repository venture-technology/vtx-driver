package repository

import (
	"context"
	"database/sql"
	"fmt"

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

func (dr *DriverRepository) CreateDriver(ctx context.Context, driver *models.Driver) error {
	sqlQuery := `INSERT INTO drivers (qrcode, name, email, password, cpf, cnh, street, number, zip, complement) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`
	_, err := dr.db.Exec(sqlQuery, driver.QrCode, driver.Name, driver.Email, driver.Password, driver.CPF, driver.CNH, driver.Street, driver.Number, driver.ZIP, driver.Complement)
	return err
}

func (dr *DriverRepository) GetDriver(ctx context.Context, cnh *string) (*models.Driver, error) {
	sqlQuery := `SELECT id, name, cpf, cnh, qrcode, email, street, number, zip, complement FROM drivers WHERE cnh = $1 LIMIT 1`
	var driver models.Driver
	err := dr.db.QueryRow(sqlQuery, *cnh).Scan(
		&driver.ID,
		&driver.Name,
		&driver.CPF,
		&driver.CNH,
		&driver.QrCode,
		&driver.Email,
		&driver.Street,
		&driver.Number,
		&driver.ZIP,
		&driver.Complement,
	)
	if err != nil || err == sql.ErrNoRows {
		return nil, err
	}
	return &driver, nil
}

func (dr *DriverRepository) UpdateDriver(ctx context.Context, driver *models.Driver) error {
	sqlQuery := `SELECT name, email, street, number, zip, complement FROM drivers WHERE cnh = $1 LIMIT 1`
	var currentDriver models.Driver
	err := dr.db.QueryRow(sqlQuery, driver.CNH).Scan(
		&currentDriver.Name,
		&currentDriver.Email,
		&currentDriver.Street,
		&currentDriver.Number,
		&currentDriver.ZIP,
		&currentDriver.Complement,
	)
	if err != nil || err == sql.ErrNoRows {
		return err
	}

	if driver.Name != "" && driver.Name != currentDriver.Name {
		currentDriver.Name = driver.Name
	}
	if driver.Email != "" && driver.Email != currentDriver.Email {
		currentDriver.Email = driver.Email
	}
	if driver.Street != "" && driver.Street != currentDriver.Street {
		currentDriver.Street = driver.Street
	}
	if driver.Number != "" && driver.Number != currentDriver.Number {
		currentDriver.Number = driver.Number
	}
	if driver.ZIP != "" && driver.ZIP != currentDriver.ZIP {
		currentDriver.ZIP = driver.ZIP
	}
	if driver.Complement != "" && driver.Complement != currentDriver.Complement {
		currentDriver.Complement = driver.Complement
	}

	sqlQueryUpdate := `UPDATE drivers SET name = $1,  email = $2, street = $3, number = $4, zip = $5, complement = $6 WHERE cnh = $7`
	_, err = dr.db.ExecContext(ctx, sqlQueryUpdate, currentDriver.Name, currentDriver.Email, currentDriver.Street, currentDriver.Number, currentDriver.ZIP, currentDriver.Complement, driver.CNH)
	return err
}

func (dr *DriverRepository) DeleteDriver(ctx context.Context, cnh *string) error {
	tx, err := dr.db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		} else if err != nil {
			_ = tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()
	_, err = tx.Exec("DELETE FROM drivers WHERE cnh = $1", *cnh)
	return err
}

func (dr *DriverRepository) AuthDriver(ctx context.Context, driver *models.Driver) (*models.Driver, error) {
	sqlQuery := `SELECT id, name, cpf, cnh, email, qrcode, password FROM drivers WHERE email = $1 LIMIT 1`
	var driverData models.Driver
	err := dr.db.QueryRow(sqlQuery, driver.Email).Scan(
		&driverData.ID,
		&driverData.Name,
		&driverData.CPF,
		&driverData.CNH,
		&driverData.Email,
		&driverData.QrCode,
		&driverData.Password,
	)
	if err != nil || err == sql.ErrNoRows {
		return nil, err
	}
	match := driverData.Password == driver.Password
	if !match {
		return nil, fmt.Errorf("email or password wrong")
	}
	driverData.Password = ""
	return &driverData, nil
}
