package managers

import (
	"fmt"
	"main/common"
	"main/database"
	"main/models"
)

type WishlistManager interface {
	Add(wishlistData *common.WishlistCreationInput) (*models.Wishlist, error)
	View(userID uint) ([]models.Wishlist, error)
	ViewAll() ([]models.Wishlist, error)
	Delete(wishlistID uint) error
}

type wishlistManager struct {
}

func NewWishlistManager() WishlistManager {
	return &wishlistManager{}
}

func (wishlistmanager *wishlistManager) Add(wishlistData *common.WishlistCreationInput) (*models.Wishlist, error) {
	newWishlist := &models.Wishlist{
		UserID:    wishlistData.UserID,
		ProductID: wishlistData.ProductID,
	}

	result := database.DB.Create(newWishlist)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to add product to wishlist: %w", result.Error)
	}

	result = database.DB.Preload("User").Preload("Product.Category").First(&newWishlist, newWishlist.Id)
	if result != nil {
		fmt.Printf("Error preloading User/Product: %v\n", result.Error)
	}

	return newWishlist, nil
}

func (wishlistmanager *wishlistManager) View(userID uint) ([]models.Wishlist, error) {
	var wishlistitems []models.Wishlist

	result := database.DB.Preload("User").Preload("Product.Category").Where("user_id=?", userID).Find(&wishlistitems)

	if result.Error != nil {
		return nil, fmt.Errorf("failed to view wishlist: %w", result.Error)
	}

	return wishlistitems, nil
}

func (wishlistmanager *wishlistManager) Delete(wishlistID uint) error {
	var wishlist models.Wishlist
	result := database.DB.Delete(&wishlist, wishlistID)
	if result.Error != nil {
		return fmt.Errorf("failed to delete wishlist item: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("wishlist item with id %d not found", wishlistID)
	}

	return nil
}

func (wishlistmanager *wishlistManager) ViewAll() ([]models.Wishlist, error) {
	var wishlists []models.Wishlist

	result := database.DB.Preload("User").Preload("Product.Category").Find(&wishlists)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to view all wishlists: %w", result.Error)
	}

	return wishlists, nil
}
