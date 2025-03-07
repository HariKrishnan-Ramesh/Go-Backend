package managers

import (
	"fmt"
	"main/common"
	"main/database"
	"main/models"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

type ProductManager interface {
	Create(productData *common.ProductCreationInput) (*models.Product, error)
	List() ([]models.Product, error)
	Get(id string) (*models.Product, error)
	Update(productID string, productData *common.ProductUpdationInput) (*models.Product, error)
	Delete(id string) error
	SearchProducts(searchTerm string) ([]models.Product, error)
	GenerateSKU() string
	SeedProducts(count int) error
	GetLastProductID() (int, error)
	SeedCategories() error
}

type productManager struct {
	//dbclient
}

func NewProductManager() ProductManager {
	return &productManager{}
}

// Create a client
func (productManager *productManager) Create(productData *common.ProductCreationInput) (*models.Product, error) {
	newProduct := &models.Product{
		SKU:         productData.SKU,
		Name:        productData.Name,
		Description: productData.Description,
		Price:       productData.Price,
		Image:       productData.Image,
		CategoryID:  productData.CategoryID,
	}

	result := database.DB.Create(newProduct)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to create new product %w", result.Error)
	}

	database.DB.Preload("Category").First(&newProduct, newProduct.Id)

	return newProduct, nil
}

func (productManager *productManager) List() ([]models.Product, error) {
	var products []models.Product
	result := database.DB.Find(&products)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to list the products %w", result.Error)
	}

	return products, nil
}

// Get a single product by ID
func (productManager *productManager) Get(id string) (*models.Product, error) {
	var product models.Product

	productID, err := strconv.Atoi(id)
	if err != nil {
		return &models.Product{}, fmt.Errorf("invalid product ID: %w", err)
	}

	fmt.Printf("Attempting to get product with ID: %d\n", productID)

	result := database.DB.Preload("Category").First(&product, productID)
	if result.Error != nil {
		return &models.Product{}, fmt.Errorf("failed to get product %w", result.Error)
	}

	return &product, nil
}

// Update a Product
func (productManager *productManager) Update(productID string, productData *common.ProductUpdationInput) (*models.Product, error) {
	var product models.Product

	result := database.DB.First(&product, productID)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to find the product: %w", result.Error)
	}

	// Update fields only if they are present in the update data.
	if productData.SKU != "" {
		product.SKU = productData.SKU
	}
	if productData.Name != "" {
		product.Name = productData.Name
	}
	if productData.Description != "" {
		product.Description = productData.Description
	}
	if productData.Price != "" {
		product.Price = productData.Price
	}
	if productData.Image != "" {
		product.Image = productData.Image
	}
	if productData.CategoryID != 0 { //Assuming CategoryID can't be 0 if it is not an empty value.
		product.CategoryID = productData.CategoryID
	}

	result = database.DB.Save(&product)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to update product: %w", result.Error)
	}

	return &product, nil
}

// Delete a product
func (productManager *productManager) Delete(id string) error {
	var product models.Product
	result := database.DB.Delete(&product, id)
	if result.Error != nil {
		return fmt.Errorf("failed to delete product %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("product with id %s not found", id)
	}

	return nil
}

func (productManager *productManager) SearchProducts(searchTerm string) ([]models.Product, error) {
	var products []models.Product

	query := database.DB.Preload("Category").Where("name LIKE ? OR description LIKE ?", "%"+searchTerm+"%", "%"+searchTerm+"%")

	result := query.Find(&products)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to search products: %w", result.Error)
	}

	return products, nil
}

func (productManager *productManager) GenerateSKU() string {
	rand.Seed(time.Now().UnixNano())

	now := time.Now()
	dateString := now.Format("20250102")

	randomNumber := rand.Intn(100000)
	randomNumberString := fmt.Sprintf("%04d", randomNumber)

	sku := fmt.Sprintf("PROD-%s-%s", dateString, randomNumberString)
	return sku
}

// Seeding products
func (productManager *productManager) SeedProducts(count int) error {

	lastProductID, err := productManager.GetLastProductID()
	if err != nil {
		fmt.Printf("Error getting last product ID: %v\n", err)
		return err //Handle appropriately
	}

	var categories []models.Category
	result := database.DB.Find(&categories)
	if result.Error != nil {
		fmt.Printf("Error getting categories for seeding %v\n", err)
		return err
	}

	if len(categories) == 0 {
		fmt.Println("No categories found.")
		return fmt.Errorf("no categories to associate with products")
	}

	for i := 0; i < count; i++ {
		nextProductID := lastProductID + i + 1

		randomCategory := categories[rand.Intn(len(categories))]

		productName, productDescription := generateProductNameAndDescription(randomCategory.Name, nextProductID)

		productData := &common.ProductCreationInput{
			//SKU:         sku,
			Name:        productName,
			Description: productDescription,
			Price:       strconv.FormatFloat(float64(rand.Intn(100000))/100.0, 'f', 2, 64),
			CategoryID:  randomCategory.Id,
		}

		var newProduct *models.Product
		maxRetries := 5
		delay := 10 * time.Millisecond //intial delay

		for retry := 0; retry <= maxRetries; retry++ {
			sku := productManager.GenerateSKU()
			productData.SKU = sku

			newProduct, err = productManager.Create(productData)

			if err == nil {
				fmt.Printf("Product %d created with SKU: %s\n", i+1, sku)
				break
			}

			if isDuplicateKeyError(err) {
				fmt.Printf("Duplicate SKU error for product %d, retry %d: %v\n", i+1, retry+1, err)
				time.Sleep(delay)
				delay *= 2
			} else {
				fmt.Printf("Non duplicate error for product %d: %v\n", i+1, err)
				return err
			}
		}

		if newProduct == nil {
			return fmt.Errorf("failes to create product %d after %d retries", i+1, maxRetries)
		}
	}
	return nil
}

// product id generation
func (productManager *productManager) GetLastProductID() (int, error) {
	var product models.Product
	result := database.DB.Last(&product)
	if result.Error != nil {
		if result.Error.Error() == "record not found" {
			return 0, nil // No products exist yet. Start from 1.
		}
		return 0, fmt.Errorf("failed to get last product: %w", result.Error)
	}
	return int(product.Id), nil
}

func (productManager *productManager) SeedCategories() error {
	// Check if categories already exist
	var existingCategories []models.Category
	result := database.DB.Find(&existingCategories)
	if result.Error != nil {
		return fmt.Errorf("failed to check existing categories: %w", result.Error)
	}

	if len(existingCategories) > 0 {
		fmt.Println("Categories already exist. Skipping seeding.")
		return nil // Categories already exist, skip seeding
	}

	// Create some default categories
	categories := []models.Category{
		{Name: "Electronics", Description: "Electronic gadgets and devices"},
		{Name: "Clothing", Description: "Apparel and fashion items"},
		{Name: "Home & Kitchen", Description: "Products for home and kitchen"},
	}

	for _, category := range categories {
		result := database.DB.Create(&category)
		if result.Error != nil {
			return fmt.Errorf("failed to create category %s: %w", category.Name, result.Error)
		}
		fmt.Printf("Created category: %s\n", category.Name)
	}

	//Create subcategories example
	var electronicsCategory models.Category
	database.DB.Where("name = ?", "Electronics").First(&electronicsCategory)

	subcategories := []models.Category{
		{Name: "Mobile phones", Description: "Mobile phones", ParentID: &electronicsCategory.Id},
		{Name: "Laptops", Description: "Laptops", ParentID: &electronicsCategory.Id},
	}

	for _, category := range subcategories {
		result := database.DB.Create(&category)
		if result.Error != nil {
			return fmt.Errorf("failed to create category %s: %w", category.Name, result.Error)
		}
		fmt.Printf("Created category: %s\n", category.Name)
	}

	return nil
}

func isDuplicateKeyError(err error) bool {
	if err == nil {
		return false
	}

	errString := fmt.Sprintf("%v", err) // Convert error to string

	return strings.Contains(errString, "Error 1062") || strings.Contains(errString, "duplicate entry")
}

func generateProductNameAndDescription(categoryName string, productID int) (string, string) {
	switch strings.ToLower(categoryName) {
	case "electronics":
		return fmt.Sprintf("Electronic Gadget #%d", productID), fmt.Sprintf("A cutting-edge electronic gadget from our premium electronics line.  Product ID: %d.", productID)
	case "clothing":
		return fmt.Sprintf("Stylish Apparel Item #%d", productID), fmt.Sprintf("A trendy and comfortable clothing item.  Part of our latest fashion collection. Product ID: %d.", productID)
	case "home & kitchen":
		return fmt.Sprintf("Home Essential #%d", productID), fmt.Sprintf("A must-have item for your home and kitchen.  High-quality and durable. Product ID: %d.", productID)
	case "mobile phones":
		return fmt.Sprintf("Smartphone Model #%d", productID), fmt.Sprintf("Next gen mobile phones to meet every customers needs. Product ID: %d", productID)
	case "laptops":
		return fmt.Sprintf("Laptop Model #%d", productID), fmt.Sprintf("High performance laptops for work and gaming. Product ID: %d", productID)
	default:
		return fmt.Sprintf("Generic Product #%d", productID), fmt.Sprintf("A general-purpose product.  Product ID: %d.", productID)
	}
}
