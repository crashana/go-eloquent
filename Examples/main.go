package main

import (
	"fmt"

	"postgres_test/models"

	"github.com/crashana/go-eloquent"
)

func main() {
	// Database connection is automatically initialized from .env file
	// No need to manually configure database connection!
	defer eloquent.GetManager().CloseAll()

	fmt.Println("=== PostgreSQL Real Database Test ===")

	// Example 1: Get all companies - NO TYPE ASSERTIONS NEEDED!
	fmt.Println("\n1. Get all companies:")
	company, err := models.Company.Where("name", "like", "%test%").First()
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		// Direct access to typed model - no type assertions!
		fmt.Println("company id:", company.ID)
		fmt.Println("company name:", company.Name)
		fmt.Println("company identification_number:", company.IdentificationNumber)
		fmt.Println("company created_at:", company.CreatedAt)
	}

	// Example 2: Get a customer and access attributes directly - NO TYPE ASSERTIONS!
	fmt.Println("\n2. Get customer by email:")
	customer, err := models.Customer.Where("email", "customer@yversy.ge").First()
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		// Direct access to typed model - no type assertions!
		fmt.Println("customer id:", customer.ID)
		fmt.Println("customer name:", customer.FirstName, customer.LastName)
		fmt.Println("customer email:", customer.Email)
		fmt.Println("customer status:", customer.Status)
		fmt.Println("customer verified:", customer.Verified)
		fmt.Println("customer company_id:", customer.CompanyID)
	}

	// Example 3: Get all customers and access attributes directly - NO TYPE ASSERTIONS!
	fmt.Println("\n3. Get all customers:")
	customers, err := models.Customer.All()
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Printf("Found %d customers:\n", len(customers))
		for i, customer := range customers {
			// Direct access to typed models - no type assertions!
			fmt.Printf("  %d. %s %s (%s) - %s\n", i+1,
				customer.FirstName,
				customer.LastName,
				customer.Email,
				customer.Status)
		}
	}

	// Example 4: Using Brand model from models package
	fmt.Println("\n4. Example: Using Brand model from models package")

	fmt.Println("Brand model available! Usage would be:")
	fmt.Println("  brand, err := models.Brand.Where(\"name\", \"Nike\").First()")
	fmt.Println("  // NO TYPE ASSERTIONS NEEDED!")
	fmt.Println("  fmt.Println(\"Brand name:\", brand.Name)        // Direct access!")
	fmt.Println("  fmt.Println(\"Brand active:\", brand.IsActive)  // No SyncAttributes needed!")

	// Example 5: Test chained queries with typed models
	fmt.Println("\n5. Test chained queries with typed models:")

	// Test chained Where clauses
	fmt.Println("Testing chained Where clauses:")
	customers2, err := models.Customer.Where("status", "active").Where("verified", false).Get()
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Printf("Found %d active unverified customers:\n", len(customers2))
		for i, customer := range customers2 {
			// Direct access to typed models - no type assertions!
			fmt.Printf("  %d. %s %s (%s) - verified: %t\n", i+1,
				customer.FirstName,
				customer.LastName,
				customer.Email,
				customer.Verified)
		}
	}

	fmt.Println("\n=== PostgreSQL Test Complete ===")
}
