// Package entity defines main entities for business logic (services), data base mapping and
// HTTP response objects if suitable. Each logic group entities in own file.
package entity

// Signa -.
type Signa struct {
	ID          string `json:"id"           example:"1789dfb3-e060-433a-9db4-cb93b42768d8"`
	Signa       string `json:"signa"        example:"a.c"`
	Description string `json:"description"  example:"Before eat"`
}
