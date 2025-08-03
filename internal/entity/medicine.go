// Package entity defines main entities for business logic (services), data base mapping and
// HTTP response objects if suitable. Each logic group entities in own file.
package entity

// Medicine -.
type Medicine struct {
	ID       string `json:"id"       example:"1789dfb3-e060-433a-9db4-cb93b42768d8"`
	Name     string `json:"name"     example:"Panadol"`
	Quantity int64  `json:"quantity" example:"100"`
	Price    int64  `json:"price"    example:"10000"`
}
