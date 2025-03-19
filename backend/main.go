package main

import (
	"net/http"
	"strconv"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func main() {
	// Connect to PostgreSQL
	dsn := "host=db user=admin password=secret dbname=products_db port=5432 sslmode=disable"
	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to database")
	}

	db.AutoMigrate(&Product{}, &Review{})

	r := gin.Default()
	r.Use(cors.Default())

	// Routes
	// Products
	r.POST("/products", CreateProduct)
	r.GET("/products/:id", GetProductByID)
	r.GET("/products", GetProducts)
	// Reviews
	r.POST("/reviews", CreateReview)
	r.GET("/products/:id/reviews", GetProductReviews)

	r.Run(":8080")
}

type Product struct {
	gorm.Model
	Name  string  `gorm:"not null" json:"name"`
	Price float64 `gorm:"not null" json:"price"`
}

type Review struct {
	gorm.Model
	Content   string `json:"content"`
	ProductID uint   `json:"product_id"`
}

type CreateProductRequest struct {
	Name  string  `json:"name" binding:"required"`
	Price float64 `json:"price" binding:"required,gt=0"`
}

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

func GetProductByID(ctx *gin.Context) {
	id := ctx.Param("id")

	var product Product

	result := db.First(&product, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve product"})
		return
	}

	ctx.JSON(http.StatusOK, product)
}

func GetProducts(ctx *gin.Context) {
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

	var products []Product
	var count int64

	db.Model(&Product{}).Count(&count)

	result := db.Order(orderClause).Limit(limit).Offset(offset).Find(&products)
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

type CreateReviewRequest struct {
	Content   string `json:"content" binding:"required"`
	ProductID uint   `json:"product_id" binding:"required"`
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

	result := db.Where("product_id = ?", product_id).Order(orderClause).Limit(limit).Offset(offset).Find(&reviews)
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
