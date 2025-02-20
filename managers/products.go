package managers

import (
	"main/common"
	"main/models"
)

type ProductManager interface{
	Create(productData *common.ProductCreationInput) (*models.Product,error)
}


type productManager struct{
	//dbclient
}

func NewProductManager() ProductManager {
	return &productManager{}
}

func (productManager *productManager) Create(productData *common.ProductCreationInput) (*models.Product,error){
return nil,nil
}