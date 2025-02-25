package handlers

import (
	"main/common"
	"main/managers"

	"github.com/gin-gonic/gin"
)

type WishlistHandler struct {
	groupName       string
	wishlistManager managers.WishlistManager
}

func NewWishlistHandler(wishlistManager managers.WishlistManager) *WishlistHandler {
	return &WishlistHandler{
		"api/wishlists",
		wishlistManager,
	}
}

func (wishlisthandler *WishlistHandler) RegisterWishlistApis(router *gin.Engine) {
	wishlistGroup := router.Group(wishlisthandler.groupName)
	wishlistGroup.POST("",wishlisthandler.Add)
}


func (wishlisthandler *WishlistHandler) Add(ctx *gin.Context) {
	wishlistData := common.NewWishlistCreationInput()
	if err := ctx.BindJSON(&wishlistData) ; err != nil {
		common.BadResponse(ctx,"Failed to bind the wishlist data")
		return
	}

	newWishlist , err := wishlisthandler.wishlistManager.Add(wishlistData)
	if err != nil {
		common.InternalServerErrorResponse(ctx, "Failed to add product to wishlist")
		return
	}

	common.SuccessResponseWithData(ctx,"Product added to wishlist successfully",newWishlist)
}