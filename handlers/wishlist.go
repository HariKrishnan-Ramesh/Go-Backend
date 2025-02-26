package handlers

import (
	"main/common"
	"main/managers"
	"strconv"

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
	wishlistGroup.POST("", wishlisthandler.Add)
	wishlistGroup.GET(":userid", wishlisthandler.View)
	wishlistGroup.GET("", wishlisthandler.ViewAll)
	wishlistGroup.DELETE(":wishlistid", wishlisthandler.Delete)
}

func (wishlisthandler *WishlistHandler) Add(ctx *gin.Context) {
	wishlistData := common.NewWishlistCreationInput()
	if err := ctx.BindJSON(&wishlistData); err != nil {
		common.BadResponse(ctx, "Failed to bind the wishlist data")
		return
	}

	newWishlist, err := wishlisthandler.wishlistManager.Add(wishlistData)
	if err != nil {
		common.InternalServerErrorResponse(ctx, "Failed to add product to wishlist")
		return
	}

	common.SuccessResponseWithData(ctx, "Product added to wishlist successfully", newWishlist)
}

func (wishlisthandler *WishlistHandler) View(ctx *gin.Context) {
	userIDStr, ok := ctx.Params.Get("userid")
	if !ok {
		common.BadResponse(ctx, "UserID required")
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		common.BadResponse(ctx, "Invalid user id")
		return
	}

	wishlistItems, err := wishlisthandler.wishlistManager.View(uint(userID))
	if err != nil {
		common.InternalServerErrorResponse(ctx, "Failed to view wishlist")
		return
	}

	common.SuccessResponseWithData(ctx, "Wishlist Successfully retrieved", wishlistItems)
}

func (wishlisthandler *WishlistHandler) Delete(ctx *gin.Context) {
	wishlistIDStr, ok := ctx.Params.Get("wishlistid")
	if !ok {
		common.InternalServerErrorResponse(ctx, "WishlistID is required")
		return
	}

	wishlistID, err := strconv.Atoi(wishlistIDStr)
	if err != nil {
		common.BadResponse(ctx, "Invalid Wishlist ID")
		return
	}

	err = wishlisthandler.wishlistManager.Delete(uint(wishlistID))
	if err != nil {
		common.InternalServerErrorResponse(ctx, "Failed to delete wishlist item")
		return
	}

	common.SuccessResponse(ctx, "Wishlist item added successfully")
}

func (wishlisthandler *WishlistHandler) ViewAll(ctx *gin.Context) {
	wishlists, err := wishlisthandler.wishlistManager.ViewAll()
	if err != nil {
		common.InternalServerErrorResponse(ctx, "Failed to view all wishlists")
		return
	}

	common.SuccessResponseWithData(ctx, "All Wishlists retrieved successfully", wishlists)
}