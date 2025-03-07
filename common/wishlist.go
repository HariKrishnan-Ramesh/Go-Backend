package common

type WishlistCreationInput struct {
	UserID    uint `json:"userID" binding:"required"`
	ProductID uint `json:"productID" binding:"required"`
}

func NewWishlistCreationInput() *WishlistCreationInput {
	return &WishlistCreationInput{}
}
