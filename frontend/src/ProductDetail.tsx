import { useQuery } from "@tanstack/react-query";
import axios from "axios";
import { useParams, useNavigate } from "react-router-dom";
import "./ProductDetail.css";

interface Product {
    ID: number;
    name: string;
    price: number;
    Reviews: Review[];
    CreatedAt: string;
    UpdatedAt: string;
}

interface Review {
    ID: number;
    Content: string;
    ProductID: number;
    CreatedAt: string;
}

export default function ProductDetail() {
    const { id } = useParams<{ id: string}>();
    const navigate = useNavigate();

    const { data: product, isLoading, error } = useQuery({
        queryKey: ["product", id],
        queryFn: async () => {
            const response = await axios.get<Product>(`http://localhost:8080/products/${id}`);
            return response.data;
        }
    });

    if (isLoading) return <div>Loading product details...</div>;
    if (error) return <div>Error: {(error as Error).message}</div>;
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
                <h1>{product.name}</h1>
                <div className="product-meta">
                    <span className="product-id">Product ID: {product.ID}</span>
                    <span className="product-price">${product.price.toFixed(2)}</span>
                </div>
                
                <div className="product-dates">
                    <p>Added on: {new Date(product.CreatedAt).toLocaleDateString()}</p>
                    <p>Last updated: {new Date(product.UpdatedAt).toLocaleDateString()}</p>
                </div>
                
                <div className="reviews-section">
                    <h2>Customer Reviews</h2>
                    {product.Reviews && product.Reviews.length > 0 ? (
                        <div className="reviews-list">
                            {product.Reviews.map(review => (
                                <div className="review-card" key={review.ID}>
                                    <p className="review-content">{review.Content}</p>
                                    <p className="review-date">
                                        Posted on {new Date(review.CreatedAt).toLocaleDateString()}
                                    </p>
                                </div>
                            ))}
                        </div>
                    ) : (
                        <p className="no-reviews">No reviews yet for this product.</p>
                    )}
                </div>
            </div>
        </div>
    );
}