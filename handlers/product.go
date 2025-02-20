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
	productGroup.POST("",productHandler.Create)
	productGroup.GET("",productHandler.List)
}

func (productHandler *ProductHandler) Create(ctx *gin.Context) {

	productData := common.NewProductCreationInput()

	err := ctx.BindJSON(&productData)
	if err!=nil{
		common.BadResponse(ctx, "Failed to bind data for product")
		return
	}

	newProduct, err := productHandler.productManager.Create(productData)
	if err!=nil{
		common.InternalServerErrorResponse(ctx, "Failed to create product")
		return
	}

	common.SuccessResponseWithData(ctx, "Product created successfully", newProduct)

}

//List all products
func (productHandler *ProductHandler) List(ctx *gin.Context){
	products,err := productHandler.productManager.List()

	if err != nil {
		common.InternalServerErrorResponse(ctx, "Failed to list products")
		return
	}

	common.SuccessResponseWithData(ctx, "Products retrieved Successfully",products)

}

//Get single Item by list
func (productHandler *ProductHandler) Get(ctx *gin.Context){
	productID, ok := ctx.Params.Get("productid")
	if !ok {
		common.BadResponse(ctx, "Product ID is required")
		return
	}

	product, err := productHandler.productManager.Get(productID)
	if err != nil {
		common.BadResponse(ctx, "Failed to get product")
		return
	}

	if product.ID == 0 {
		common.BadResponse(ctx, "Product Not Found")
		return
	}

	common.SuccessResponseWithData(ctx, "Product Retrieved Successfully", product)
}



