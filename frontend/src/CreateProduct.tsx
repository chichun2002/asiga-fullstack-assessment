import { useState } from "react";
import { useMutation } from "@tanstack/react-query";
import axios from "axios";
import { useNavigate } from "react-router-dom";
import "./CreateProduct.css";

interface CreateProductRequest {
  name: string;
  price: number;
}

export default function CreateProduct() {
  const navigate = useNavigate();
  const [name, setName] = useState("");
  const [price, setPrice] = useState("");
  const [error, setError] = useState("");

  const createProductMutation = useMutation({
    mutationFn: async (newProduct: CreateProductRequest) => {
      const response = await axios.post("http://localhost:8080/products", newProduct);
      return response.data;
    },
    onSuccess: (data) => {
      // Navigate to the product detail page of the newly created product
      navigate(`/products/${data.ID}`);
    },
    onError: (err: Error) => {
      setError(`Failed to create product: ${err.message}`);
    }
  });

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    
    // Validate form
    if (!name.trim()) {
      setError("Product name is required");
      return;
    }
    
    const numericPrice = parseFloat(price);
    if (isNaN(numericPrice) || numericPrice <= 0) {
      setError("Price must be a positive number");
      return;
    }
    
    // Clear any previous errors
    setError("");
    
    // Submit the form
    createProductMutation.mutate({
      name: name.trim(),
      price: numericPrice
    });
  };

  return (
    <div className="create-product-container">
      <button 
        className="back-button" 
        onClick={() => navigate('/products')}
      >
        ← Back to Products
      </button>
      
      <div className="form-card">
        <h1>Create New Product</h1>
        
        {error && <div className="error-message">{error}</div>}
        
        <form onSubmit={handleSubmit}>
          <div className="form-group">
            <label htmlFor="name">Product Name</label>
            <input
              type="text"
              id="name"
              value={name}
              onChange={(e) => setName(e.target.value)}
              placeholder="Enter product name"
              disabled={createProductMutation.isPending}
            />
          </div>
          
          <div className="form-group">
            <label htmlFor="price">Price ($)</label>
            <input
              type="number"
              id="price"
              value={price}
              onChange={(e) => setPrice(e.target.value)}
              placeholder="Enter price"
              step="0.01"
              min="0.01"
              disabled={createProductMutation.isPending}
            />
          </div>
          
          <div className="form-actions">
            <button
              type="button"
              className="cancel-button"
              onClick={() => navigate(-1)}
              disabled={createProductMutation.isPending}
            >
              Cancel
            </button>
            
            <button
              type="submit"
              className="submit-button"
              disabled={createProductMutation.isPending}
            >
              {createProductMutation.isPending ? "Creating..." : "Create Product"}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
}