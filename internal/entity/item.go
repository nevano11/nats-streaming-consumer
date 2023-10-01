package entity

import "fmt"

type Item struct {
	ChrtId      int    `json:"chrt_id"`
	TrackNumber string `json:"track_number"`
	Price       int    `json:"price"`
	Rid         string `json:"rid"`
	Name        string `json:"name"`
	Sale        int    `json:"sale"`
	Size        string `json:"size"`
	TotalPrice  int    `json:"total_price"`
	NmId        int    `json:"nm_id"`
	Brand       string `json:"brand"`
	Status      int    `json:"status"`
}

func (i *Item) String() string {
	return fmt.Sprintf("Item: {ChrtId: %d, TrackNumber: %s, Price: %d, Rid: %s, Name: %s, Sale: %d, Size: %s, TotalPrice: %d, NmId: %d, Brand: %s, Status: %d}",
		i.ChrtId,
		i.TrackNumber,
		i.Price,
		i.Rid,
		i.Name,
		i.Sale,
		i.Size,
		i.TotalPrice,
		i.NmId,
		i.Brand,
		i.Status)
}
