package models

import (
	"time"

	"github.com/crashana/go-eloquent"
)

// CustomerModel - Eloquent style model
type CustomerModel struct {
	*eloquent.BaseModel

	// Struct fields for direct access - like Eloquent
	ID                   string    `json:"id" db:"id"`
	Status               string    `json:"status" db:"status"`
	Email                string    `json:"email" db:"email"`
	Phone                string    `json:"phone" db:"phone"`
	FirstName            string    `json:"first_name" db:"first_name"`
	LastName             string    `json:"last_name" db:"last_name"`
	IdentificationNumber string    `json:"identification_number" db:"identification_number"`
	IsOrganizationOwner  bool      `json:"is_organization_owner" db:"is_organization_owner"`
	Verified             bool      `json:"verified" db:"verified"`
	CompanyID            string    `json:"company_id" db:"company_id"`
	LegalForm            string    `json:"legal_form" db:"legal_form"`
	CreatedAt            time.Time `json:"created_at" db:"created_at"`
	UpdatedAt            time.Time `json:"updated_at" db:"updated_at"`
}

// Define relationships for CustomerModel
func (c *CustomerModel) Company() *eloquent.Relationship {
	rb := eloquent.NewRelationshipBuilder(c)
	return rb.BelongsTo("company", "CompanyModel")
}

// NewCustomer creates a new CustomerModel instance
func NewCustomer() *CustomerModel {
	customer := &CustomerModel{
		BaseModel: eloquent.NewBaseModel(),
	}

	// Configure the model based on actual table structure
	customer.Table("customers").
		PrimaryKey("id").
		Fillable("status", "email", "phone", "first_name", "last_name", "identification_number", "password", "is_organization_owner", "verified", "company_id", "legal_form").
		Hidden("password", "remember_token").
		Casts(map[string]string{
			"is_organization_owner": "boolean",
			"verified":              "boolean",
			"last_login_at":         "datetime",
			"created_at":            "datetime",
			"updated_at":            "datetime",
			"deleted_at":            "datetime",
		})

	return customer
}

// Global static instance for Customer model - Eloquent style
var Customer = eloquent.NewModelStatic(func() *CustomerModel {
	return NewCustomer()
})
