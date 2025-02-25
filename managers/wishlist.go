package managers

import (
	"fmt"
	"main/common"
	"main/database"
	"main/models"
)

type WishlistManager interface {
	Add(wishlistData *common.WishlistCreationInput) (*models.Wishlist, error)
}

type wishlistManager struct {
}

func NewWishlistManager() WishlistManager {
	return &wishlistManager{}
}

func (wishlistmanager *wishlistManager) Add(wishlistData *common.WishlistCreationInput) (*models.Wishlist, error) {
	newWishlist := &models.Wishlist{
		UserID: wishlistData.UserID,
		ProductID: wishlistData.ProductID,
	}

	result := database.DB.Create(newWishlist)
	if result.Error != nil {
		return nil,fmt.Errorf("failed to add product to wishlist: %w",result.Error)
	}

	return newWishlist,nil
}