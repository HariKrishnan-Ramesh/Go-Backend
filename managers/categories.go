package managers

import (
	"fmt"
	"main/common"
	"main/database"
	"main/models"
	"strconv"
)

type CategoryManager interface {
	Create(categoryData *common.CategoryCreationInput) (*models.Category, error)
	List() ([]models.Category, error)
	Get(id string) (*models.Category, error)
	Update(categoryID string, categoryData *common.CategoryUpdationInput) (*models.Category, error)
	Delete(id string) error
	GetCategoryCount() (int, error)
}

type categoryManager struct {
}

func NewCategoryManager() CategoryManager {
	return &categoryManager{}
}

func (categoryManager *categoryManager) Create(categoryData *common.CategoryCreationInput) (*models.Category, error) {
	newCategory := &models.Category{
		Name:        categoryData.Name,
		Description: categoryData.Description,
		ParentID:    categoryData.ParentID,
	}

	result := database.DB.Create(newCategory)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to create category: %w", result.Error)
	}

	return newCategory, nil
}

func (categoryManager *categoryManager) List() ([]models.Category, error) {
	var categories []models.Category

	result := database.DB.Preload("Children").Preload("Products").Find(&categories)

	if result.Error != nil {
		return nil, fmt.Errorf("failed to list categories: %w", result.Error)
	}
	return categories, nil
}

func (categoryManager *categoryManager) Get(id string) (*models.Category, error) {
	var category models.Category
	result := database.DB.Preload("Children").Preload("Products").First(&category, id)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to get category: %w", result.Error)
	}
	return &category, nil
}

func (categoryManager *categoryManager) Update(categoryID string, categoryData *common.CategoryUpdationInput) (*models.Category, error) {

	id, err := strconv.Atoi(categoryID)
	if err != nil {
		return nil, fmt.Errorf("invalid category id: %w", err)
	}

	var category models.Category
	result := database.DB.First(&category, id)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to find category: %w", result.Error)
	}
	if categoryData.Name != "" {
		category.Name = categoryData.Name 
	}
	if categoryData.Description != "" {
		category.Description = categoryData.Description
	}
	if categoryData.ParentID != nil {
		category.ParentID = categoryData.ParentID
	}

	result = database.DB.Save(&category)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to update category: %w", result.Error)
	}
	return &category, nil
}

func (categoryManager *categoryManager) Delete(id string) error {
	var category models.Category
	result := database.DB.Delete(&category, id)
	if result.Error != nil {
		return fmt.Errorf("failed to delete category: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("category with id %s not found", id)
	}
	return nil
}

func (categoryManager *categoryManager) GetCategoryCount() (int, error) {
	var count int64
	result := database.DB.Model(&models.Category{}).Count(&count)
	if result.Error != nil {
		return 0, fmt.Errorf("failed to get category count: %w", result.Error)
	}
	return int(count), nil
}
