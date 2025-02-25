package handlers

import (
	"main/common"
	"main/managers"

	"github.com/gin-gonic/gin"
)

type CartHandler struct {
	groupName    string
	cartManager managers.CartManager
}

func NewCartHandler(cartManager managers.CartManager) *CartHandler {
	return &CartHandler{
		"api/carts",
		cartManager,
	}
}

func (carthandler *CartHandler) RegisterCartApis(router *gin.Engine){
	cartGroup := router.Group(carthandler.groupName)
	cartGroup.POST("",carthandler.Add)
}

func (carthandler *CartHandler) Add(ctx *gin.Context) {
	cartData := common.NewCartCreationInput()
	if err := ctx.BindJSON(&cartData); err != nil {
		common.BadResponse(ctx, "Failed to bind cart data")
		return
	}

	newCartItem, err := carthandler.cartManager.Add(cartData)
	if err != nil {
		common.InternalServerErrorResponse(ctx," Failed to add product to cart")
		return
	}

	common.SuccessResponseWithData(ctx, "Product add to cart successfully", newCartItem)
}