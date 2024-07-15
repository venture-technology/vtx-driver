package models

import "github.com/venture-technology/vtx-driver/utils"

type Driver struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	Email      string `json:"email"`
	Password   string `json:"password"`
	CPF        string `json:"cpf"`
	CNH        string `json:"cnh"`
	QrCode     string `json:"qrcode"`
	Street     string `json:"street"`
	Number     string `json:"number"`
	ZIP        string `json:"zip"`
	Complement string `json:"complement"`
}

func (d *Driver) ValidateCnh() bool {

	return utils.IsCNH(d.CNH)

}
