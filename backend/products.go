package main

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CreateProductRequest struct {
	Name  string  `json:"name" binding:"required"`
	Price float64 `json:"price" binding:"required,gt=0"`
}

type UpdateProductRequest struct {
	Name  *string  `json:"name"`
	Price *float64 `json:"price"`
}

// Creates a product assigning it the next avaliable id json body for name and price
func CreateProduct(ctx *gin.Context) {
	var input CreateProductRequest

	if err := ctx.ShouldBindBodyWithJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	product := Product{
		Name:  input.Name,
		Price: input.Price,
	}

	result := db.Create(&product)
	if result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create product"})
		return
	}

	ctx.JSON(http.StatusCreated, product)
}

// Returns a product given the id
func GetProductByID(ctx *gin.Context) {
	id := ctx.Param("id")

	var product Product

	result := db.First(&product, id)
	if result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	ctx.JSON(http.StatusOK, product)
}

// Gets Products from Databse Given Parameters including paging size, search filter and sorts
func GetProducts(ctx *gin.Context) {
	// Default values
	limit := 10
	page := 1
	sortBy := "created_at"
	sortOrder := "desc"
	search := ""

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

	if searchParam := ctx.Query("search"); searchParam != "" {
		search = searchParam
	}

	offset := (page - 1) * limit
	orderClause := sortBy + " " + sortOrder

	var products []Product
	var count int64

	query := db.Model(&Product{})

	if search != "" {
		searchTerm := "%" + search + "%"
		query = query.Where("name ILIKE ?", searchTerm)
	}

	query.Count(&count)

	result := query.Order(orderClause).Limit(limit).Offset(offset).Find(&products)

	if result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve products"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"products": products,
		"pagination": gin.H{
			"total": count,
			"page":  page,
			"limit": limit,
			"pages": (count + int64(limit) - 1) / int64(limit),
		},
	})
}

// Deletes Product from Database
func DeleteProduct(ctx *gin.Context) {
	id := ctx.Param("id")

	result := db.Delete(&Product{}, id)
	if result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete product"})
		return
	}
	if result.RowsAffected == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Product deleted successfully"})
}

// Modifies Existing Product json body for name and price
func UpdateProduct(ctx *gin.Context) {
	id := ctx.Param("id")

	var product Product
	result := db.First(&product, id)
	if result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	var input UpdateProductRequest
	if err := ctx.ShouldBindBodyWithJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updates := map[string]interface{}{}

	if input.Name != nil {
		updates["name"] = *input.Name
	}
	if input.Price != nil {
		updates["price"] = *input.Price
	}

	if len(updates) == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "No valid fields to update"})
		return
	}

	result = db.Model(&product).Updates(updates)
	if result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update product"})
		return
	}

	db.First(&product, id)

	ctx.JSON(http.StatusOK, product)
}
