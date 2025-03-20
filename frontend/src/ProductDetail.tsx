import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { useCallback, useState } from "react";
import axios from "axios";
import { useParams, useNavigate } from "react-router-dom";
import "./ProductDetail.css";
import ProductReviews from "./ProductReviews";

interface Product {
    ID: number;
    name: string;
    description: string;
    price: number;
    CreatedAt: string;
    UpdatedAt: string;
}

interface UpdateProductRequest {
    name?: string;
    description?: string;
    price?: number;
}

export default function ProductDetail() {
    const { id } = useParams<{ id: string}>();
    const navigate = useNavigate();
    const queryClient = useQueryClient();
    
    // Edit product state
    const [showEditForm, setShowEditForm] = useState(false);
    const [editName, setEditName] = useState("");
    const [editDescription, setEditDescription] = useState("");
    const [editPrice, setEditPrice] = useState("");
    const [editError, setEditError] = useState("");

    // Fetch Product
    const { data: product, isLoading: productLoading, error : productError } = useQuery({
        queryKey: ["product", id],
        queryFn: async () => {
            const response = await axios.get<Product>(`http://localhost:8080/products/${id}`);
            return response.data;
        }
    });

    const updateProductMutation = useMutation({
        mutationFn: async (productData: UpdateProductRequest) => {
            const response = await axios.patch(`http://localhost:8080/products/${id}`, productData);
            return response.data;
        },
        onSuccess: () => {
            setShowEditForm(false);
            setEditError("");
            // Invalidate product data to refresh UI
            queryClient.invalidateQueries({ queryKey: ["product", id] });
        },
        onError: (err: Error) => {
            setEditError(`Failed to update product: ${err.message}`);
        }
    });

    const handleEditProduct = (e: React.FormEvent) => {
        e.preventDefault();
        
        if (!editName.trim() && !editPrice.trim() && !editDescription.trim()) {
            setEditError("At least one field must be modified");
            return;
        }

        const updateData: UpdateProductRequest = {};
        
        if (editName.trim()) {
            updateData.name = editName.trim();
        }

        if (editDescription.trim()) {
          updateData.description = editDescription.trim();
        }
        
        if (editPrice.trim()) {
            const priceValue = parseFloat(editPrice);
            if (isNaN(priceValue) || priceValue <= 0) {
                setEditError("Price must be a positive number");
                return;
            }
            updateData.price = priceValue;
        }
        
        setEditError("");
        updateProductMutation.mutate(updateData);
    };

    // Initialize edit form values when product data is loaded
    const handleShowEditForm = useCallback(() => {
        if (product) {
            setEditName(product.name);
            setEditDescription(product.description);
            setEditPrice(product.price.toString());
            setShowEditForm(true);
        }
    }, [product]);

    if (productLoading) return <div>Loading product details...</div>;
    if (productError ) return <div>Error: {(productError  as Error).message}</div>;
    if (!product) return <div>Product not found</div>

    return (
        <div className="product-detail-container">
          <button
            className="back-button"
            onClick={() => navigate('/products')}
          >
            ← Back to Products
          </button>
          
          <div className="product-detail-card">
            {!showEditForm ? (
              <>
                {/* Details */}
                <div className="product-header">
                  <h1>{product.name}</h1>
                  <h2>{product.description}</h2>
                  <button 
                    className="edit-product-button"
                    onClick={handleShowEditForm}
                  >
                    Edit Product
                  </button>
                </div>
                
                <div className="product-meta">
                  <span className="product-id">Product ID: {product.ID}</span>
                  <span className="product-price">${product.price.toFixed(2)}</span>
                </div>
              </>
            ) : (
              <div className="edit-form-container">
                {/* Edit Details */}
                <h2>Edit Product</h2>
                
                {editError && <div className="error-message">{editError}</div>}
                
                <form onSubmit={handleEditProduct}>
                  <div className="form-group">
                    <label htmlFor="product-name">Name</label>
                    <input
                      id="product-name"
                      type="text"
                      value={editName}
                      onChange={(e) => setEditName(e.target.value)}
                      placeholder="Product name"
                      disabled={updateProductMutation.isPending}
                    />
                  </div>

                  <div className="form-group">
                    <label htmlFor="product-description">Name</label>
                    <input
                      id="product-description"
                      type="text"
                      value={editDescription}
                      onChange={(e) => setEditDescription(e.target.value)}
                      placeholder="Product description"
                      disabled={updateProductMutation.isPending}
                    />
                  </div>
                  
                  <div className="form-group">
                    <label htmlFor="product-price">Price</label>
                    <input
                      id="product-price"
                      type="number"
                      step="0.01"
                      min="0.01"
                      value={editPrice}
                      onChange={(e) => setEditPrice(e.target.value)}
                      placeholder="Product price"
                      disabled={updateProductMutation.isPending}
                    />
                  </div>
                  
                  <div className="form-actions">
                    <button
                      type="button"
                      className="cancel-button"
                      onClick={() => {
                        setShowEditForm(false);
                        setEditError("");
                      }}
                      disabled={updateProductMutation.isPending}
                    >
                      Cancel
                    </button>
                    
                    <button
                      type="submit"
                      className="submit-button"
                      disabled={updateProductMutation.isPending}
                    >
                      {updateProductMutation.isPending ? "Updating..." : "Update Product"}
                    </button>
                  </div>
                </form>
              </div>
            )}
            
            <div className="product-dates">
              <p>Added on: {new Date(product.CreatedAt).toLocaleDateString()}</p>
              <p>Last updated: {new Date(product.UpdatedAt).toLocaleDateString()}</p>
            </div>
            
            {/* Render the Reviews component */}
            <ProductReviews productId={id || ""} />
          </div>
        </div>
    );
}