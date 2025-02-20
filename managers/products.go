package managers

import (
	"fmt"
	"main/common"
	"main/database"
	"main/models"
)

type ProductManager interface {
	Create(productData *common.ProductCreationInput) (*models.Product, error)
	List() ([]models.Product, error)
	Get(id string) (*models.Product, error)
	Update(productID string, productData *common.ProductUpdationInput) (*models.Product, error)
	Delete(id string) error
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
