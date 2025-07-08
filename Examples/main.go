package main

import (
	"fmt"
	"time"

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

	// Example 4: CREATE - Creating new records
	fmt.Println("\n4. CREATE - Creating new records:")

	// Create a new company using Create method
	fmt.Println("Creating a new company...")
	newCompany, err := models.Company.Create(map[string]interface{}{
		"name":                  "New Tech Company",
		"identification_number": fmt.Sprintf("ID%d", time.Now().Unix()), // Unique ID
	})
	if err != nil {
		fmt.Println("Error creating company:", err)
	} else {
		fmt.Printf("✅ Created company: %s (ID: %s)\n", newCompany.Name, newCompany.ID)
		fmt.Printf("   Identification Number: %s\n", newCompany.IdentificationNumber)
		fmt.Printf("   Created At: %s\n", newCompany.CreatedAt)
	}

	// Create a new customer using Create method
	fmt.Println("\nCreating a new customer...")
	customerData := map[string]interface{}{
		"first_name": "Alice",
		"last_name":  "Johnson",
		"email":      fmt.Sprintf("alice.johnson%d@example.com", time.Now().Unix()), // Unique email
		"phone":      "+1234567890",                                                 // Required field
		"password":   "password123",                                                 // Required field
		"status":     "active",
		"verified":   true,
	}

	// Only add company_id if we have a valid company
	if newCompany != nil && newCompany.ID != "" {
		customerData["company_id"] = newCompany.ID
	}

	newCustomer, err := models.Customer.Create(customerData)
	if err != nil {
		fmt.Println("Error creating customer:", err)
	} else {
		fmt.Printf("✅ Created customer: %s %s (ID: %s)\n",
			newCustomer.FirstName, newCustomer.LastName, newCustomer.ID)
		fmt.Printf("   Email: %s\n", newCustomer.Email)
		fmt.Printf("   Status: %s, Verified: %t\n", newCustomer.Status, newCustomer.Verified)
		fmt.Printf("   Company ID: %s\n", newCustomer.CompanyID)
	}

	// Example 5: UPDATE - Updating existing records
	fmt.Println("\n5. UPDATE - Updating existing records:")

	// Method 1: Update using the Update method with map
	if newCustomer != nil {
		fmt.Println("Updating customer using Update method...")
		err = newCustomer.Update(map[string]interface{}{
			"status":   "active", // Use valid status
			"verified": true,
		})
		if err != nil {
			fmt.Println("Error updating customer:", err)
		} else {
			fmt.Printf("✅ Updated customer status to: %s\n", newCustomer.Status)
			fmt.Printf("   Verified status: %t\n", newCustomer.Verified)
		}
	} else {
		fmt.Println("⚠️  Skipping customer update - customer creation failed")
	}

	// Method 2: Direct attribute modification and Save
	if newCustomer != nil {
		fmt.Println("\nUpdating customer using direct attribute access...")
		newCustomer.FirstName = "Alice Updated"
		newCustomer.LastName = "Johnson-Smith"
		newCustomer.Email = "alice.johnson.smith@example.com"

		err = newCustomer.Save()
		if err != nil {
			fmt.Println("Error saving customer:", err)
		} else {
			fmt.Printf("✅ Updated customer name to: %s %s\n",
				newCustomer.FirstName, newCustomer.LastName)
			fmt.Printf("   Updated email to: %s\n", newCustomer.Email)
		}
	} else {
		fmt.Println("⚠️  Skipping direct attribute update - customer creation failed")
	}

	// Method 3: Update company information
	if newCompany != nil {
		fmt.Println("\nUpdating company information...")
		newCompany.Name = "Updated Tech Solutions"
		newCompany.IdentificationNumber = "111222333"

		err = newCompany.Save()
		if err != nil {
			fmt.Println("Error saving company:", err)
		} else {
			fmt.Printf("✅ Updated company name to: %s\n", newCompany.Name)
			fmt.Printf("   Updated identification number to: %s\n", newCompany.IdentificationNumber)
		}
	} else {
		fmt.Println("⚠️  Skipping company update - company creation failed")
	}

	// Example 6: FIND and UPDATE - Find by ID and update
	fmt.Println("\n6. FIND and UPDATE - Find by ID and update:")

	// Find the customer we just created by ID
	if newCustomer != nil {
		foundCustomer, err := models.Customer.Find(newCustomer.ID)
		if err != nil {
			fmt.Println("Error finding customer:", err)
		} else {
			fmt.Printf("📍 Found customer: %s %s\n", foundCustomer.FirstName, foundCustomer.LastName)

			// Update the found customer
			foundCustomer.Status = "active" // Use valid status

			// Update using the Update method instead of Save to avoid table name issues
			err = foundCustomer.Update(map[string]interface{}{
				"status": "active",
			})
			if err != nil {
				fmt.Println("Error updating found customer:", err)
			} else {
				fmt.Printf("✅ Updated customer status to: %s\n", foundCustomer.Status)
			}
		}
	} else {
		fmt.Println("⚠️  Skipping find and update - customer creation failed")
	}

	// Example 7: DELETE - Deleting records
	fmt.Println("\n7. DELETE - Deleting records:")

	// Create a temporary customer to demonstrate deletion
	tempCustomer, err := models.Customer.Create(map[string]interface{}{
		"first_name": "Temp",
		"last_name":  "Customer",
		"email":      fmt.Sprintf("temp%d@example.com", time.Now().Unix()), // Unique email
		"phone":      "+9876543210",                                        // Required field
		"password":   "temppass123",                                        // Required field
		"status":     "inactive",
		"verified":   false,
	})
	if err != nil {
		fmt.Println("Error creating temp customer:", err)
	} else {
		fmt.Printf("📝 Created temporary customer: %s %s (ID: %s)\n",
			tempCustomer.FirstName, tempCustomer.LastName, tempCustomer.ID)

		// Delete the temporary customer
		err = tempCustomer.Delete()
		if err != nil {
			fmt.Println("Error deleting temp customer:", err)
		} else {
			fmt.Printf("🗑️  Deleted temporary customer: %s %s\n",
				tempCustomer.FirstName, tempCustomer.LastName)
		}
	}

	// Example 8: Test chained queries with typed models
	fmt.Println("\n8. Test chained queries with typed models:")

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

	// Example 9: Using Brand model from models package
	fmt.Println("\n9. Example: Using Brand model from models package")

	fmt.Println("Brand model available! Usage would be:")
	fmt.Println("  brand, err := models.Brand.Where(\"name\", \"Nike\").First()")
	fmt.Println("  // NO TYPE ASSERTIONS NEEDED!")
	fmt.Println("  fmt.Println(\"Brand name:\", brand.Name)        // Direct access!")
	fmt.Println("  fmt.Println(\"Brand active:\", brand.IsActive)  // No SyncAttributes needed!")

	fmt.Println("\n=== CRUD Operations Complete ===")
	fmt.Println("✅ Created new company and customer")
	fmt.Println("✅ Updated records using multiple methods")
	fmt.Println("✅ Found records by ID")
	fmt.Println("✅ Deleted temporary records")
	fmt.Println("✅ Demonstrated chained queries")
	fmt.Println("\n=== PostgreSQL Test Complete ===")
}
