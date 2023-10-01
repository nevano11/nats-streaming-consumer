package entity

import "fmt"

type Delivery struct {
	Name    string `json:"name"`
	Phone   string `json:"phone"`
	Zip     string `json:"zip"`
	City    string `json:"city"`
	Address string `json:"address"`
	Region  string `json:"region"`
	Email   string `json:"email"`
}

func (d *Delivery) String() string {
	return fmt.Sprintf(
		"Delivery: {Name: %s, Phone: %s, Zip: %s, City: %s, Address: %s, Region: %s, Email: %s}",
		d.Name,
		d.Phone,
		d.Zip,
		d.City,
		d.Address,
		d.Region,
		d.Email)
}
