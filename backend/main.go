package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

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
	r.DELETE("/products/:id", DeleteProduct)
	r.PATCH("/products/:id", UpdateProduct)

	// Reviews
	r.POST("/reviews", CreateReview)
	r.GET("/products/:id/reviews", GetProductReviews)
	r.DELETE("/reviews/:id", DeleteReview)
	r.PATCH("/reviews/:id", UpdateReview)

	r.Run(":8080")
}
