package handlers

import (
	"main/common"
	"main/managers"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CartHandler struct {
	groupName    string
	cartManager managers.CartManager
}

func NewCartHandler(cartManager managers.CartManager) *CartHandler {
	return &CartHandler{
		"api/cart",
		cartManager,
	}
}

func (carthandler *CartHandler) RegisterCartApis(router *gin.Engine) {
	cartGroup := router.Group(carthandler.groupName)
	cartGroup.POST("", carthandler.Add)
	cartGroup.GET(":userid", carthandler.View)
	cartGroup.PATCH(":cartid", carthandler.Update)
	cartGroup.DELETE(":cartid", carthandler.Delete) 
}

func (carthandler *CartHandler) Add(ctx *gin.Context) {
	cartData := common.NewCartCreationInput()
	if err := ctx.BindJSON(&cartData); err != nil {
		common.BadResponse(ctx, "Failed to bind cart data")
		return
	}

	newCartItem, err := carthandler.cartManager.Add(cartData)
	if err != nil {
		common.InternalServerErrorResponse(ctx, "Failed to add product to cart")
		return
	}

	common.SuccessResponseWithData(ctx, "Product added to cart successfully", newCartItem)
}

func (carthandler *CartHandler) View(ctx *gin.Context) {
	userIDStr := ctx.Param("userid")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		common.BadResponse(ctx, "Invalid User ID")
		return
	}

	cartItems, err := carthandler.cartManager.View(uint(userID))
	if err != nil {
		common.InternalServerErrorResponse(ctx, "Failed to view cart")
		return
	}

	common.SuccessResponseWithData(ctx, "Cart retrieved successfully", cartItems)
}

func (carthandler *CartHandler) Update(ctx *gin.Context) {
	cartIDStr := ctx.Param("cartid")
	cartID, err := strconv.Atoi(cartIDStr)
	if err != nil {
		common.BadResponse(ctx, "Invalid Cart ID")
		return
	}

	updateData := common.NewCartUpdateInput()
	if err := ctx.BindJSON(&updateData); err != nil {
		common.BadResponse(ctx, "Failed to bind update data")
		return
	}

	updatedCartItem, err := carthandler.cartManager.Update(uint(cartID), updateData)
	if err != nil {
		common.InternalServerErrorResponse(ctx, "Failed to update cart item")
		return
	}

	common.SuccessResponseWithData(ctx, "Cart item updated successfully", updatedCartItem)
}

func (carthandler *CartHandler) Delete(ctx *gin.Context) {
	cartIDStr := ctx.Param("cartid") 
	cartID, err := strconv.Atoi(cartIDStr)
	if err != nil {
		common.BadResponse(ctx, "Invalid Cart ID")
		return
	}

	err = carthandler.cartManager.Delete(uint(cartID))
	if err != nil {
		common.InternalServerErrorResponse(ctx, "Failed to delete the cart item")
		return
	}

	common.SuccessResponse(ctx, "Cart item deleted successfully")
}