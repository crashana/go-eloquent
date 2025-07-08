package models

import (
	"time"

	"github.com/crashana/go-eloquent"
)

// CompanyModel - Laravel style model
type CompanyModel struct {
	*eloquent.BaseModel

	// Struct fields for direct access - like Laravel
	ID                   string    `json:"id" db:"id"`
	Name                 string    `json:"name" db:"name"`
	IdentificationNumber string    `json:"identification_number" db:"identification_number"`
	CreatedAt            time.Time `json:"created_at" db:"created_at"`
	UpdatedAt            time.Time `json:"updated_at" db:"updated_at"`
}

// Define relationships for CompanyModel
func (c *CompanyModel) Customers() *eloquent.Relationship {
	rb := eloquent.NewRelationshipBuilder(c)
	return rb.HasMany("customers", "CustomerModel")
}

// NewCompany creates a new CompanyModel instance
func NewCompany() *CompanyModel {
	company := &CompanyModel{
		BaseModel: eloquent.NewBaseModel(),
	}

	// Configure the model based on actual table structure
	company.Table("companies").
		PrimaryKey("id").
		Fillable("name", "identification_number").
		Casts(map[string]string{
			"created_at": "datetime",
			"updated_at": "datetime",
		})

	return company
}

// Global static instance for Company model - Laravel style
var Company = eloquent.NewModelStatic(func() *CompanyModel {
	return NewCompany()
})
