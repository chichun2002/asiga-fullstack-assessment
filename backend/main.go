package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// Connect to PostgreSQL
	dsn := "host=db user=admin password=secret dbname=products_db port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to database")
	}

	db.AutoMigrate(&Product{}, &Review{})

	r := gin.Default()
	r.Use(cors.Default())

	// Routes
	// r.GET("/products", GetProducts)
	// r.POST("/products", CreateProduct)

	r.Run(":8080")
}

type Product struct {
	gorm.Model
	Name    string
	Price   float64
	Reviews []Review
}

type Review struct {
	gorm.Model
	Content   string
	ProductID uint
}
