package common




type ProductCreationInput struct {
	SKU         string `json:"sku"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Price       string `json:"price"`
	Image       string `json:"image,omitempty"`
	CategoryID  uint   `json:"categoryID"`
}


type ProductUpdationInput struct {
	SKU         string `json:"sku" gorm:"uniqueIndex"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Price       string `json:"price"`
	Image       string `json:"image,omitempty"`
	CategoryID  uint   `json:"categoryID"`
}


func NewProductCreationInput() *ProductCreationInput{
	return &ProductCreationInput{}
}

func NewProductUpdationInput() *ProductUpdationInput{
	return &ProductUpdationInput{}
}