package common

type CartCreationInput struct {
	UserID    uint `json:"userID" binding:"required"`
	ProductID uint `json:"productID" binding:"required"`
	Quantity  uint `json:"quantity" binding:"required,min=1"` // min quantity is 1
}

func NewCartCreationInput() *CartCreationInput {
	return &CartCreationInput{}
}

type CartUpdateInput struct {
	Quantity uint `json:"quantity" binding:"required,min=1"`
}

func NewCartUpdateInput() *CartUpdateInput {
	return &CartUpdateInput{}
}
