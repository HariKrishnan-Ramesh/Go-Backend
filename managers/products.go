package managers

import (
	"fmt"
	"main/common"
	"main/database"
	"main/models"
	"math/rand"
	"strconv"
	"time"
)

type ProductManager interface {
	Create(productData *common.ProductCreationInput) (*models.Product, error)
	List() ([]models.Product, error)
	Get(id string) (*models.Product, error)
	Update(productID string, productData *common.ProductUpdationInput) (*models.Product, error)
	Delete(id string) error
	GenerateSKU() string
	SeedProducts(count int) error
	GetLastProductID() (int, error)
}

type productManager struct {
	//dbclient
}

func NewProductManager() ProductManager {
	return &productManager{}
}

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

	result := database.DB.First(&product, id)
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

func (productManager *productManager) GenerateSKU() string {
	rand.Seed(time.Now().UnixNano())

	now := time.Now()
	dateString := now.Format("20250102")

	randomNumber := rand.Intn(100000)
	randomNumberString := fmt.Sprintf("%04d",randomNumber)

	sku := fmt.Sprintf("PROD-%s-%s",dateString,randomNumberString)
	return sku
}


func (productManager *productManager) SeedProducts(count int) error {

	lastProductID, err := productManager.GetLastProductID()
	if err != nil {
		fmt.Printf("Error getting last product ID: %v\n", err)
		return err //Handle appropriately
	}


	for i:=0; i<count; i++ {
		nextProductID := lastProductID + i +1
		sku := productManager.GenerateSKU()
		productData := &common.ProductCreationInput{
			SKU : sku,
			Name : fmt.Sprintf("Product %d",nextProductID ),
			Description: "Generated Product",
			Price : strconv.FormatFloat(float64(rand.Intn(100000))/100.0,'f',2,64),
			CategoryID: uint(rand.Intn(10)+1),
		}

		_,err := productManager.Create(productData)
		if err != nil {
			fmt.Printf("Error creating product %d: %v\n", i+1, err)
			return err
		}

		fmt.Printf("Product %d created with SKU: %s\n",i+1,sku)
	}
	return nil
}


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