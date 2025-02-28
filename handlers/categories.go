package handlers

import (
	"main/common"
	"main/managers"

	"github.com/gin-gonic/gin"
)

type CategoryHandler struct {
	groupName       string
	categoryManager managers.CategoryManager
}

func NewCategoryHandler(categoryManager managers.CategoryManager) *CategoryHandler {
	return &CategoryHandler{
		"api/categories",
		categoryManager,
	}
}

func (categoryhandler *CategoryHandler) RegisterCategoryApis(router *gin.Engine) {
	categoryGroup := router.Group(categoryhandler.groupName)
	categoryGroup.GET("", categoryhandler.List)
	categoryGroup.GET(":categoryid", categoryhandler.Get) 
	categoryGroup.POST("", categoryhandler.Create)
	categoryGroup.PATCH(":categoryid", categoryhandler.Update) 
	categoryGroup.DELETE(":categoryid", categoryhandler.Delete)
}

func (categoryhandler *CategoryHandler) List(ctx *gin.Context) {
	categories, err := categoryhandler.categoryManager.List()
	if err != nil {
		common.InternalServerErrorResponse(ctx, "Failed to list categories")
		return
	}
	common.SuccessResponseWithData(ctx, "Categories retrieved successfully", categories)
}

func (categoryhandler *CategoryHandler) Get(ctx *gin.Context) {
	categoryID := ctx.Param("categoryid") 
	if categoryID == "" {
		common.BadResponse(ctx, "Category ID is required")
		return
	}

	category, err := categoryhandler.categoryManager.Get(categoryID)
	if err != nil {
		common.BadResponse(ctx, "Failed to get category")
		return
	}

	if category.Id == 0 {
		common.BadResponse(ctx, "Category not found")
		return
	}

	common.SuccessResponseWithData(ctx, "Category retrieved successfully", category)
}

func (categoryhandler *CategoryHandler) Create(ctx *gin.Context) {
	categoryData := common.NewCategoryCreationInput() 
	if err := ctx.BindJSON(&categoryData); err != nil {
		common.BadResponse(ctx, "Failed to bind category data")
		return
	}

	newCategory, err := categoryhandler.categoryManager.Create(categoryData)
	if err != nil {
		common.InternalServerErrorResponse(ctx, "Failed to create category")
		return
	}
	common.SuccessResponseWithData(ctx, "Category created successfully", newCategory)
}

func (categoryhandler *CategoryHandler) Update(ctx *gin.Context) {
	categoryID := ctx.Param("categoryid") 
	if categoryID == "" {
		common.BadResponse(ctx, "Category ID is required")
		return
	}

	categoryUpdateData := common.NewCategoryUpdationInput() 
	if err := ctx.BindJSON(&categoryUpdateData); err != nil {
		common.BadResponse(ctx, "Failed to bind update data")
		return
	}

	updatedCategory, err := categoryhandler.categoryManager.Update(categoryID, categoryUpdateData)
	if err != nil {
		common.InternalServerErrorResponse(ctx, "Failed to update category")
		return
	}

	common.SuccessResponseWithData(ctx, "Category updated successfully", updatedCategory)
}

func (categoryhandler *CategoryHandler) Delete(ctx *gin.Context) {
	categoryID := ctx.Param("categoryid")
	if categoryID == "" {
		common.BadResponse(ctx, "Category ID is required")
		return
	}

	err := categoryhandler.categoryManager.Delete(categoryID)
	if err != nil {
		common.InternalServerErrorResponse(ctx, "Failed to delete category")
		return
	}

	common.SuccessResponse(ctx, "Category deleted successfully")
}
