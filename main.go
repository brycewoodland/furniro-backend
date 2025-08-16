package main

import (
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/gin-gonic/gin"
)

type Product struct {
	ID          string  `json:"id"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Discount    string  `json:"discount,omitempty"`
	IsNew       bool    `json:"isNew,omitempty"`
}

var products = []Product{
	{ID: "1", Title: "Syltherine", Description: "Stylish cafe chair", Price: 2500000.00, Discount: "-30%", IsNew: false},
	{ID: "2", Title: "Leviosa", Description: "Minimalist sofa", Price: 2500000.00, Discount: "", IsNew: false},
	{ID: "3", Title: "Lolito", Description: "Luxury couch", Price: 7000000.00, Discount: "-50%", IsNew: false},
	{ID: "4", Title: "Respira", Description: "Outdoor table & stools", Price: 500000.00, Discount: "", IsNew: true},
	{ID: "5", Title: "Grifo", Description: "Night lamp", Price: 1500000.00, Discount: "", IsNew: false},
	{ID: "6", Title: "Muggo", Description: "Small mug", Price: 100000.00, Discount: "", IsNew: true},
	{ID: "7", Title: "Pingky", Description: "Bedroom set", Price: 7000000.00, Discount: "-50%", IsNew: false},
	{ID: "8", Title: "Potty", Description: "Flower pot", Price: 50000.00, Discount: "", IsNew: true},
}

var productsMutex = &sync.RWMutex{}

func getProducts(c *gin.Context) {
	productsMutex.RLock()
	defer productsMutex.RUnlock()
	c.JSON(http.StatusOK, products)
}

func getProductByID(c *gin.Context) {
	id := c.Param("id")

	productsMutex.RLock()
	defer productsMutex.RUnlock()

	for _, p := range products {
		if p.ID == id {
			c.JSON(http.StatusOK, p)
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"message": "Product not found"})
}

func addProduct(c *gin.Context) {
	var newProduct Product

	if err := c.BindJSON(&newProduct); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	productsMutex.Lock()
	defer productsMutex.Unlock()

	for _, p := range products {
		if p.ID == newProduct.ID {
			c.JSON(http.StatusConflict, gin.H{"message": "Product with this ID already exists"})
			return
		}
	}

	products = append(products, newProduct)
	c.JSON(http.StatusCreated, newProduct)
}

func updateProduct(c *gin.Context) {
	id := c.Param("id")
	var updatedProduct Product

	if err := c.BindJSON(&updatedProduct); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	productsMutex.Lock()
	defer productsMutex.Unlock()

	found := false
	for i, p := range products {
		if p.ID == id {
			updatedProduct.ID = id
			products[i] = updatedProduct
			found = true
			break
		}
	}

	if !found {
		c.JSON(http.StatusNotFound, gin.H{"message": "Product not found"})
		return
	}
	c.JSON(http.StatusOK, updatedProduct)
}

func deleteProduct(c *gin.Context) {
	id := c.Param("id")

	productsMutex.Lock()
	defer productsMutex.Unlock()

	initialLen := len(products)
	for i, p := range products {
		if p.ID == id {
			products = append(products[:i], products[i+1:]...)
			break
		}
	}

	if len(products) == initialLen {
		c.JSON(http.StatusNotFound, gin.H{"message": "Product not found"})
		return
	}
	c.Status(http.StatusNoContent)
}

func main() {
	router := gin.Default()
	router.Use(func(c *gin.Context) {
		// c.Writer.Header().Set("Access-Control-Allow-Origin", "http://127.0.0.1:5500") // Replace with your frontend's actual origin
		c.Writer.Header().Set("Access-Control-Allow-Origin", "https://furniro-project-bryce.netlify.app")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// API routes for products
	api := router.Group("/api")
	{
		api.GET("/products", getProducts)
		api.GET("/products/:id", getProductByID)
		api.POST("/products", addProduct)
		api.PUT("/products/:id", updateProduct)
		api.DELETE("/products/:id", deleteProduct)
	}

	// Start the server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Println("Server started on port", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
