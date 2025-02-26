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

func (ch *CategoryHandler) RegisterCategoryApis(router *gin.Engine) {
	categoryGroup := router.Group(ch.groupName)
	categoryGroup.GET("", ch.List)
	categoryGroup.GET(":categoryid", ch.Get) 
	categoryGroup.POST("", ch.Create)
	categoryGroup.PATCH(":categoryid", ch.Update) 
	categoryGroup.DELETE(":categoryid", ch.Delete)
}

func (ch *CategoryHandler) List(ctx *gin.Context) {
	categories, err := ch.categoryManager.List()
	if err != nil {
		common.InternalServerErrorResponse(ctx, "Failed to list categories")
		return
	}
	common.SuccessResponseWithData(ctx, "Categories retrieved successfully", categories)
}

func (ch *CategoryHandler) Get(ctx *gin.Context) {
	categoryID := ctx.Param("categoryid") 
	if categoryID == "" {
		common.BadResponse(ctx, "Category ID is required")
		return
	}

	category, err := ch.categoryManager.Get(categoryID)
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

func (ch *CategoryHandler) Create(ctx *gin.Context) {
	categoryData := common.NewCategoryCreationInput() 
	if err := ctx.BindJSON(&categoryData); err != nil {
		common.BadResponse(ctx, "Failed to bind category data")
		return
	}

	newCategory, err := ch.categoryManager.Create(categoryData)
	if err != nil {
		common.InternalServerErrorResponse(ctx, "Failed to create category")
		return
	}
	common.SuccessResponseWithData(ctx, "Category created successfully", newCategory)
}

func (ch *CategoryHandler) Update(ctx *gin.Context) {
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

	updatedCategory, err := ch.categoryManager.Update(categoryID, categoryUpdateData)
	if err != nil {
		common.InternalServerErrorResponse(ctx, "Failed to update category")
		return
	}

	common.SuccessResponseWithData(ctx, "Category updated successfully", updatedCategory)
}

func (ch *CategoryHandler) Delete(ctx *gin.Context) {
	categoryID := ctx.Param("categoryid")
	if categoryID == "" {
		common.BadResponse(ctx, "Category ID is required")
		return
	}

	err := ch.categoryManager.Delete(categoryID)
	if err != nil {
		common.InternalServerErrorResponse(ctx, "Failed to delete category")
		return
	}

	common.SuccessResponse(ctx, "Category deleted successfully")
}
