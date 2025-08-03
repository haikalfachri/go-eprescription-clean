// Package entity defines main entities for business logic (services), data base mapping and
// HTTP response objects if suitable. Each logic group entities in own file.
package entity

// Patient -.
type Patient struct {
	ID     string `json:"id"    example:"1789dfb3-e060-433a-9db4-cb93b42768d8"`
	Name   string `json:"name"  example:"John Doe"`
	Age    int    `json:"age"   example:"30"`
	Gender string `json:"gender" example:"male"`
}
