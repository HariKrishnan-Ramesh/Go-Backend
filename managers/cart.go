package managers

import (
	"fmt"
	"main/common"
	"main/database"
	"main/models"
)

type CartManager interface {
	Add(cartData *common.CartCreationInput)(*models.Cart,error)
	View(userID uint) ([]models.Cart, error)
	Update(cartID uint, updateData *common.CartUpdateInput) (*models.Cart, error)
}

type cartManager struct{
	//dbclient
}

func NewCartManager() CartManager {
	return &cartManager{}
}

func (cartmanager *cartManager) Add(cartData *common.CartCreationInput)(*models.Cart,error) {
	var existingCartItem models.Cart
	result := database.DB.Where("user_id = ? AND product_id = ?",cartData.UserID,cartData.ProductID).First(&existingCartItem)

	if result.Error == nil {
		existingCartItem.Quantity += cartData.Quantity
		result = database.DB.Save(&existingCartItem)
		if result.Error != nil {
			return nil, fmt.Errorf("failed to update cart item quantity: %w", result.Error)
		}
		database.DB.Preload("User").Preload("Product.Category").First(&existingCartItem, existingCartItem.Id)
		return &existingCartItem, nil
	} else if result.Error != nil && result.Error.Error() != "record not found" {
		return nil, fmt.Errorf("failed to check existing cart item: %w",result.Error)
	}

	newCartItem := &models.Cart{
		UserID: cartData.UserID,
		ProductID: cartData.ProductID,
		Quantity: cartData.Quantity,
	}

	result = database.DB.Create(newCartItem)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to add product to cart: %w", result.Error)
	}

	database.DB.Preload("User").Preload("Product.Category").First(&newCartItem, newCartItem.Id)
	return newCartItem,nil
}


func (cartmanager *cartManager) View(userID uint) ([]models.Cart, error) {
	var cartItems []models.Cart
	result := database.DB.Preload("User").Preload("Product.Category").Where("user_id = ?").Find(&cartItems)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to view cart: %w",result.Error)
	}
	return cartItems, nil
}


func (cartmanager *cartManager) Update(cartID uint, updateData *common.CartUpdateInput) (*models.Cart, error) {
	var cartItem models.Cart
	result := database.DB.First(&cartItem, cartID)
	if result.Error != nil {
		return nil, fmt.Errorf("cart item not found: %w",result.Error)
	}

	cartItem.Quantity = updateData.Quantity
	result = database.DB.Save(&cartItem)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to update cart item: %w", result.Error)
	}

	database.DB.Preload("User").Preload("Product.Category").First(&cartItem, cartItem.Id)
	return &cartItem, nil

}

