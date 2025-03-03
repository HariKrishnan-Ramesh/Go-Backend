package common

type CategoryCreationInput struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	ParentID    uint   `json:"parentID"`
}

func NewCategoryCreationInput() *CategoryCreationInput {
	return &CategoryCreationInput{}
}

type CategoryUpdationInput struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	ParentID    uint   `json:"parentID,omitempty"`
}

func NewCategoryUpdationInput() *CategoryUpdationInput {
	return &CategoryUpdationInput{}
}
