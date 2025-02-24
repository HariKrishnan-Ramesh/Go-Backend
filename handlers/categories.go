package handlers

import (
	"main/common"
	"main/managers"

	"github.com/gin-gonic/gin"
)

type CategoryHandler struct {
	groupName      string
	categoryManager managers.CategoryManager
}

func NewCategoryHandler(categoryManager managers.CategoryManager) *CategoryHandler {
	return &CategoryHandler{
		"api/categories",
		categoryManager,
	}
}

func (categoryHandler *CategoryHandler) RegisterCategoryApis(router *gin.Engine) {
	categoryGroup := router.Group(categoryHandler.groupName)
	categoryGroup.GET("", categoryHandler.List)
	categoryGroup.GET(":categoryid/", categoryHandler.Get)
	categoryGroup.POST("", categoryHandler.Create)
	categoryGroup.PATCH(":categoryid/", categoryHandler.Update)
	categoryGroup.DELETE(":categoryid/", categoryHandler.Delete)
}

func (categoryHandler *CategoryHandler) List(ctx *gin.Context) {
	categories, err := categoryHandler.categoryManager.List()
	if err != nil {
		common.InternalServerErrorResponse(ctx, "Failed to list categories")
		return
	}
	common.SuccessResponseWithData(ctx, "Categories retrieved successfully", categories)
}

func (categoryHandler *CategoryHandler) Get(ctx *gin.Context) {
	categoryID, ok := ctx.Params.Get("categoryid")
	if !ok {
		common.BadResponse(ctx, "Category ID is required")
		return
	}

	category, err := categoryHandler.categoryManager.Get(categoryID)
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

func (categoryHandler *CategoryHandler) Create(ctx *gin.Context) {
	categoryData := common.NewCategoryCreationInput()  // Assuming you have this
	if err := ctx.BindJSON(&categoryData); err != nil {
		common.BadResponse(ctx, "Failed to bind category data")
		return
	}

	newCategory, err := categoryHandler.categoryManager.Create(categoryData)
	if err != nil {
		common.InternalServerErrorResponse(ctx, "Failed to create category")
		return
	}
	common.SuccessResponseWithData(ctx, "Category created successfully", newCategory)
}

func (categoryHandler *CategoryHandler) Update(ctx *gin.Context) {
	categoryID, ok := ctx.Params.Get("categoryid")
	if !ok {
		common.BadResponse(ctx, "Category ID is required")
		return
	}

	categoryUpdateData := common.NewCategoryUpdationInput()  // Assuming you have this
	if err := ctx.BindJSON(&categoryUpdateData); err != nil {
		common.BadResponse(ctx, "Failed to bind update data")
		return
	}

	updatedCategory, err := categoryHandler.categoryManager.Update(categoryID, categoryUpdateData)
	if err != nil {
		common.InternalServerErrorResponse(ctx, "Failed to update category")
		return
	}

	common.SuccessResponseWithData(ctx, "Category updated successfully", updatedCategory)
}

func (categoryHandler *CategoryHandler) Delete(ctx *gin.Context) {
	categoryID, ok := ctx.Params.Get("categoryid")
	if !ok {
		common.BadResponse(ctx, "Category ID is required")
		return
	}

	err := categoryHandler.categoryManager.Delete(categoryID)
	if err != nil {
		common.InternalServerErrorResponse(ctx, "Failed to delete category")
		return
	}

	common.SuccessResponse(ctx, "Category deleted successfully")
}