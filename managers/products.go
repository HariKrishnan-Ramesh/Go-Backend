package managers

import (
	"fmt"
	"main/common"
	"main/database"
	"main/models"
)

type ProductManager interface{
	Create(productData *common.ProductCreationInput) (*models.Product,error)
	List() ([]models.Product,error)
	Get(id string) (*models.Product, error)
}


type productManager struct{
	//dbclient
}

func NewProductManager() ProductManager {
	return &productManager{}
}

func (productManager *productManager) Create(productData *common.ProductCreationInput) (*models.Product,error){
	newProduct := &models.Product{
		SKU : productData.SKU,
		Name: productData.Name,
		Description: productData.Description,
		Price: productData.Price,
		Image: productData.Image,
		CategoryID: productData.CategoryID,
	}

	result := database.DB.Create(newProduct)
	if result.Error != nil {
		return nil,fmt.Errorf("failed to create new product %w",result.Error)
	}

return newProduct,nil
}

func (productManager *productManager) List() ([]models.Product,error){
	var products []models.Product
	result := database.DB.Find(&products)
	if result.Error != nil {
		return nil,fmt.Errorf("failed to list the products %w",result.Error)
	}

	return products,nil
}

//Get a single product by ID
func (productManager *productManager) Get(id string) (*models.Product, error) {
	var product models.Product

	result := database.DB.First(&product, id)
	if result.Error != nil {
		return &models.Product{}, fmt.Errorf("failed to get product %w",result.Error)
	}

	return &product, nil
}


