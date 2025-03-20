package main

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CreateReviewRequest struct {
	Content   string `json:"content" binding:"required"`
	ProductID uint   `json:"product_id" binding:"required"`
}

type UpdateReviewRequest struct {
	Content   *string `json:"content"`
	ProductID *uint   `json:"product_id"`
}

func CreateReview(ctx *gin.Context) {
	var input CreateReviewRequest

	if err := ctx.ShouldBindBodyWithJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var product Product
	if result := db.First(&product, input.ProductID); result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	review := Review{
		Content:   input.Content,
		ProductID: input.ProductID,
	}

	result := db.Create(&review)
	if result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create product"})
		return
	}

	ctx.JSON(http.StatusCreated, review)
}

func GetProductReviews(ctx *gin.Context) {
	product_id := ctx.Param("id")
	// Default values
	limit := 10
	page := 1
	sortBy := "created_at"
	sortOrder := "desc"

	// Param Parse
	if limitParam := ctx.Query("limit"); limitParam != "" {
		if parsedLimit, err := strconv.Atoi(limitParam); err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	if pageParam := ctx.Query("page"); pageParam != "" {
		if parsedPage, err := strconv.Atoi(pageParam); err == nil && parsedPage > 0 {
			page = parsedPage
		}
	}

	if sort := ctx.Query("sort"); sort != "" {
		validColumns := map[string]bool{
			"name": true, "price": true, "created_at": true,
		}
		if validColumns[sort] {
			sortBy = sort
		}
	}

	if order := ctx.Query("order"); order != "" {
		validColumns := map[string]bool{
			"asc": true, "desc": true,
		}
		if validColumns[order] {
			sortOrder = order
		}
	}

	offset := (page - 1) * limit
	orderClause := sortBy + " " + sortOrder

	var reviews []Review
	var count int64

	db.Model(&Review{}).Where("product_id = ?", product_id).Count(&count)

	result := db.Model(&Review{}).Where("product_id = ?", product_id).Order(orderClause).Limit(limit).Offset(offset).Find(&reviews)
	if result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve products"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"reviews": reviews,
		"pagination": gin.H{
			"total": count,
			"page":  page,
			"limit": limit,
			"pages": (count + int64(limit) - 1) / int64(limit),
		},
	})
}

func DeleteReview(ctx *gin.Context) {
	id := ctx.Param("id")

	result := db.Delete(&Review{}, id)
	if result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete review"})
		return
	}
	if result.RowsAffected == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Review not found"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Review deleted successfully"})
}

func UpdateReview(ctx *gin.Context) {
	id := ctx.Param("id")

	var review Review
	result := db.First(&review, id)
	if result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Review not found"})
		return
	}

	var input UpdateReviewRequest
	if err := ctx.ShouldBindBodyWithJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updates := map[string]interface{}{}

	if input.Content != nil {
		updates["content"] = *input.Content
	}
	if input.ProductID != nil {
		updates["product_id"] = *input.ProductID
	}

	if len(updates) == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "No valid fields to update"})
		return
	}

	result = db.Model(&review).Updates(updates)
	if result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update review"})
		return
	}

	db.First(&review, id)

	ctx.JSON(http.StatusOK, review)
}
