package models

import (
	"time"

	"github.com/crashana/go-eloquent"
)

// BrandModel - Eloquent style model (example)
type BrandModel struct {
	*eloquent.BaseModel

	// Just define struct fields with db tags - that's it!
	ID          string    `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Description string    `json:"description" db:"description"`
	IsActive    bool      `json:"is_active" db:"is_active"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// NewBrand creates a new BrandModel instance
func NewBrand() *BrandModel {
	brand := &BrandModel{
		BaseModel: eloquent.NewBaseModel(),
	}

	// Configure the model
	brand.Table("brands").
		PrimaryKey("id").
		Fillable("name", "description", "is_active").
		Casts(map[string]string{
			"is_active":  "boolean",
			"created_at": "datetime",
			"updated_at": "datetime",
		})

	return brand
}

// Global static instance for Brand model - Eloquent style
var Brand = eloquent.NewModelStatic(func() *BrandModel {
	return NewBrand()
})
