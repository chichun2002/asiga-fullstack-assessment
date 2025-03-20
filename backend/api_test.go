package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	baseURL = "http://localhost:8080" // Your actual running server URL
)

// Test data structs that match your API models
type ProductTest struct {
	ID          uint      `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type ReviewTest struct {
	ID        uint      `json:"id"`
	Content   string    `json:"content"`
	ProductID uint      `json:"product_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Helper function to create a test product
func createTestProductLive(t *testing.T) ProductTest {
	// Create product payload
	productData := map[string]interface{}{
		"name":        fmt.Sprintf("Test Product %d", time.Now().Unix()),
		"description": "This is a test product created by automated tests",
		"price":       99.99,
	}

	// Convert to JSON
	payloadBytes, err := json.Marshal(productData)
	require.NoError(t, err)

	// Send request
	resp, err := http.Post(
		fmt.Sprintf("%s/products", baseURL),
		"application/json",
		bytes.NewBuffer(payloadBytes),
	)
	require.NoError(t, err)
	defer resp.Body.Close()

	// Check status code
	require.Equal(t, http.StatusCreated, resp.StatusCode)

	// Parse response
	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	var product ProductTest
	err = json.Unmarshal(body, &product)
	require.NoError(t, err)

	return product
}

// Helper function to delete a test product
func deleteTestProductLive(t *testing.T, productID uint) {
	// Create DELETE request
	req, err := http.NewRequest(
		http.MethodDelete,
		fmt.Sprintf("%s/products/%d", baseURL, productID),
		nil,
	)
	require.NoError(t, err)

	// Send request
	client := &http.Client{}
	resp, err := client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	// Check status code
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

// Helper to create a test review
func createTestReviewLive(t *testing.T, productID uint) ReviewTest {
	// Create review payload
	reviewData := map[string]interface{}{
		"content":    fmt.Sprintf("Test Review %d", time.Now().Unix()),
		"product_id": productID,
	}

	// Convert to JSON
	payloadBytes, err := json.Marshal(reviewData)
	require.NoError(t, err)

	// Send request
	resp, err := http.Post(
		fmt.Sprintf("%s/reviews", baseURL),
		"application/json",
		bytes.NewBuffer(payloadBytes),
	)
	require.NoError(t, err)
	defer resp.Body.Close()

	// Check status code
	require.Equal(t, http.StatusCreated, resp.StatusCode)

	// Parse response
	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	var review ReviewTest
	err = json.Unmarshal(body, &review)
	require.NoError(t, err)

	return review
}

// Helper to delete a test review
func deleteTestReviewLive(t *testing.T, reviewID uint) {
	// Create DELETE request
	req, err := http.NewRequest(
		http.MethodDelete,
		fmt.Sprintf("%s/reviews/%d", baseURL, reviewID),
		nil,
	)
	require.NoError(t, err)

	// Send request
	client := &http.Client{}
	resp, err := client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	// Check status code
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestCreateAndGetProduct(t *testing.T) {
	// Create a test product
	createdProduct := createTestProductLive(t)
	defer deleteTestProductLive(t, createdProduct.ID)

	// Verify the created product
	assert.NotZero(t, createdProduct.ID)
	assert.Contains(t, createdProduct.Name, "Test Product")
	assert.Equal(t, "This is a test product created by automated tests", createdProduct.Description)
	assert.Equal(t, 99.99, createdProduct.Price)

	// Get the product
	resp, err := http.Get(fmt.Sprintf("%s/products/%d", baseURL, createdProduct.ID))
	require.NoError(t, err)
	defer resp.Body.Close()

	// Check status code
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Parse response
	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	var product ProductTest
	err = json.Unmarshal(body, &product)
	require.NoError(t, err)

	// Verify returned product
	assert.Equal(t, createdProduct.ID, product.ID)
	assert.Equal(t, createdProduct.Name, product.Name)
	assert.Equal(t, createdProduct.Description, product.Description)
	assert.Equal(t, createdProduct.Price, product.Price)
}

func TestUpdateProduct(t *testing.T) {
	// Create a test product
	product := createTestProductLive(t)
	defer deleteTestProductLive(t, product.ID)

	// Update payload
	updateData := map[string]interface{}{
		"name":  "Updated Product Name",
		"price": 149.99,
	}

	// Convert to JSON
	payloadBytes, err := json.Marshal(updateData)
	require.NoError(t, err)

	// Create PATCH request
	req, err := http.NewRequest(
		http.MethodPatch,
		fmt.Sprintf("%s/products/%d", baseURL, product.ID),
		bytes.NewBuffer(payloadBytes),
	)
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	// Send request
	client := &http.Client{}
	resp, err := client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	// Check status code
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Get the updated product
	resp, err = http.Get(fmt.Sprintf("%s/products/%d", baseURL, product.ID))
	require.NoError(t, err)
	defer resp.Body.Close()

	// Parse response
	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	var updatedProduct ProductTest
	err = json.Unmarshal(body, &updatedProduct)
	require.NoError(t, err)

	// Verify updates
	assert.Equal(t, "Updated Product Name", updatedProduct.Name)
	assert.Equal(t, 149.99, updatedProduct.Price)
	assert.Equal(t, product.Description, updatedProduct.Description) // Should not have changed
}

func TestGetAllProducts(t *testing.T) {
	// Create 2 test products
	product1 := createTestProductLive(t)
	defer deleteTestProductLive(t, product1.ID)

	product2 := createTestProductLive(t)
	defer deleteTestProductLive(t, product2.ID)

	// Get all products
	resp, err := http.Get(fmt.Sprintf("%s/products", baseURL))
	require.NoError(t, err)
	defer resp.Body.Close()

	// Check status code
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Parse response
	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	// Create a struct that matches your API response format
	type ProductsResponse struct {
		Products []ProductTest `json:"products"`
	}

	var response ProductsResponse
	err = json.Unmarshal(body, &response)
	require.NoError(t, err)

	products := response.Products
	// Check if our created products are in the list
	found1, found2 := false, false
	for _, p := range products {
		if p.ID == product1.ID {
			found1 = true
		}
		if p.ID == product2.ID {
			found2 = true
		}
	}

	assert.True(t, found1, "First product not found in results")
	assert.True(t, found2, "Second product not found in results")
}

func TestReviewOperations(t *testing.T) {
	// Create a test product
	product := createTestProductLive(t)
	defer deleteTestProductLive(t, product.ID)

	// Create a review
	review := createTestReviewLive(t, product.ID)
	defer deleteTestReviewLive(t, review.ID)

	// Verify the review was created
	assert.NotZero(t, review.ID)
	assert.Contains(t, review.Content, "Test Review")
	assert.Equal(t, product.ID, review.ProductID)

	// Get reviews for product
	resp, err := http.Get(fmt.Sprintf("%s/products/%d/reviews", baseURL, product.ID))
	require.NoError(t, err)
	defer resp.Body.Close()

	// Check status code
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Parse response
	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	type ReviewResponse struct {
		Reviews []ReviewTest `json:"reviews"`
	}

	var response ReviewResponse
	err = json.Unmarshal(body, &response)
	require.NoError(t, err)

	reviews := response.Reviews
	// Should have at least one review
	assert.NotEmpty(t, reviews)

	// Find our review
	found := false
	for _, r := range reviews {
		if r.ID == review.ID {
			found = true
			assert.Equal(t, review.Content, r.Content)
			assert.Equal(t, product.ID, r.ProductID)
			break
		}
	}
	assert.True(t, found, "Review not found in results")

	// Update review
	updateData := map[string]interface{}{
		"content": "Updated review content",
	}

	// Convert to JSON
	payloadBytes, err := json.Marshal(updateData)
	require.NoError(t, err)

	// Create PATCH request
	req, err := http.NewRequest(
		http.MethodPatch,
		fmt.Sprintf("%s/reviews/%d", baseURL, review.ID),
		bytes.NewBuffer(payloadBytes),
	)
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	// Send request
	client := &http.Client{}
	resp, err = client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	// Check status code
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Get reviews again to check the update
	resp, err = http.Get(fmt.Sprintf("%s/products/%d/reviews", baseURL, product.ID))
	require.NoError(t, err)
	defer resp.Body.Close()

	// Parse response
	body, err = io.ReadAll(resp.Body)
	require.NoError(t, err)

	err = json.Unmarshal(body, &response)
	require.NoError(t, err)

	reviews = response.Reviews
	// Find our updated review
	found = false
	for _, r := range reviews {
		if r.ID == review.ID {
			found = true
			assert.Equal(t, "Updated review content", r.Content)
			break
		}
	}
	assert.True(t, found, "Updated review not found in results")
}
