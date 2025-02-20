package handlers

import (
	"main/common"
	"main/managers"

	"github.com/gin-gonic/gin"
)

type ProductHandler struct {
	groupName      string
	productManager managers.ProductManager
}


func NewProductHandler(productManager managers.ProductManager) *ProductHandler {
	return &ProductHandler{
		"api/product" ,
		productManager,
	}
}


func (productHandler *ProductHandler) RegisterUserApis(router *gin.Engine){
	productGroup := router.Group(productHandler.groupName)
	productGroup.POST("/signup",productHandler.Create)
}

func (productHandler *ProductHandler) Create(ctx *gin.Context) {

	productData := common.NewProductCreationInput()

	err := ctx.BindJSON(&productData)
	if err!=nil{
		common.BadResponse(ctx, "Failed to bind data for product")
	}

	
}
