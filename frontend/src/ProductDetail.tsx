import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { useCallback, useState } from "react";
import axios from "axios";
import { useParams, useNavigate } from "react-router-dom";
import "./ProductDetail.css";

interface Product {
    ID: number;
    name: string;
    price: number;
    CreatedAt: string;
    UpdatedAt: string;
}

interface Review {
    ID: number;
    content: string;
    product_id: number;
    CreatedAt: string;
}

interface ReviewResponse {
    reviews: Review[];
    pagination: {
        total: number;
        page: number;
        limit: number;
        pages: number;
    };
}

interface AddReviewRequest {
    content: string;
    product_id: number;
}

export default function ProductDetail() {
    const { id } = useParams<{ id: string}>();
    const navigate = useNavigate();
    const queryClient = useQueryClient();
    const [reviewPage, setReviewPage] = useState(1);
    const [reviewLimit] = useState(5); // Number of reviews per page

    const [reviewContent, setReviewContent] = useState("");
    const [showReviewForm, setShowReviewForm] = useState(false);
    const [reviewError, setReviewError] = useState("");

    // Fetch Product
    const { data: product, isLoading: productLoading, error : productError } = useQuery({
        queryKey: ["product", id],
        queryFn: async () => {
            const response = await axios.get<Product>(`http://localhost:8080/products/${id}`);
            return response.data;
        }
    });

    // Fetch Reviews
    const { data: reviewsData, isLoading: reviewsLoading, error: reviewsError} = useQuery({
        queryKey: ["productReviews", id, reviewPage, reviewLimit],
        queryFn: async () => {
            const response = await axios.get<ReviewResponse>(`http://localhost:8080/products/${id}/reviews`,
                {
                    params: {
                        page: reviewPage,
                        limit: reviewLimit,
                        sort: "created_at",
                        order: "desc"
                    }
                }
            );
            return response.data;
        },
        enabled: !!id
    });

    const addReviewMutation = useMutation({
        mutationFn: async (reviewData: AddReviewRequest) => {
            const response = await axios.post("http://localhost:8080/reviews", reviewData);
            return response.data;
        },
        onSuccess: () => {
            setReviewContent("");
            setShowReviewForm(false);

            queryClient.invalidateQueries({ queryKey: ["productReviews", id] });
        },
        onError: (err: Error) => {
            setReviewError(`Failed to add review': ${err.message}`)
        }
    });

    const handleSubmitReview = (e: React.FormEvent) => {
        e.preventDefault();

        if (!reviewContent.trim()) {
            setReviewError("Review content is required");
            return;
        }
        setReviewError("");

        addReviewMutation.mutate({
            content: reviewContent.trim(),
            product_id: parseInt(id || "0", 10)
        });
    };

    if (productLoading) return <div>Loading product details...</div>;
    if (productError ) return <div>Error: {(productError  as Error).message}</div>;
    if (!product) return <div>Product not found</div>

    const reviews = reviewsData?.reviews || [];
    const pagination = reviewsData?.pagination;

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
              <div className="reviews-header">
                <h2>Customer Reviews</h2>
                {!showReviewForm && (
                  <button 
                    className="add-review-button"
                    onClick={() => setShowReviewForm(true)}
                  >
                    + Add Review
                  </button>
                )}
              </div>
              
              {/* Inline Review Form */}
              {showReviewForm && (
                <div className="review-form-container">
                  <h3>Write a Review</h3>
                  
                  {reviewError && <div className="error-message">{reviewError}</div>}
                  
                  <form onSubmit={handleSubmitReview}>
                    <div className="form-group">
                      <textarea
                        value={reviewContent}
                        onChange={(e) => setReviewContent(e.target.value)}
                        placeholder="Write your review here..."
                        rows={4}
                        disabled={addReviewMutation.isPending}
                      />
                    </div>
                    
                    <div className="form-actions">
                      <button
                        type="button"
                        className="cancel-button"
                        onClick={() => {
                          setShowReviewForm(false);
                          setReviewError("");
                        }}
                        disabled={addReviewMutation.isPending}
                      >
                        Cancel
                      </button>
                      
                      <button
                        type="submit"
                        className="submit-button"
                        disabled={addReviewMutation.isPending}
                      >
                        {addReviewMutation.isPending ? "Submitting..." : "Submit Review"}
                      </button>
                    </div>
                  </form>
                </div>
              )}
              
              {/* Reviews List */}
              {reviewsLoading && !reviewsData ? (
                <div className="loading-reviews">Loading reviews...</div>
              ) : reviewsError ? (
                <div className="error">Error loading reviews: {(reviewsError as Error).message}</div>
              ) : reviews.length > 0 ? (
                <>
                  <div className="reviews-list">
                    {reviews.map(review => (
                      <div className="review-card" key={review.ID}>
                        <p className="review-content">{review.content}</p>
                        <p className="review-date">
                          Posted on {new Date(review.CreatedAt).toLocaleDateString()}
                        </p>
                      </div>
                    ))}
                  </div>
                  
                  {/* Pagination controls for reviews */}
                  {pagination && pagination.pages > 1 && (
                    <div className="reviews-pagination">
                      <button
                        onClick={() => setReviewPage(p => Math.max(p - 1, 1))}
                        disabled={reviewPage === 1 || reviewsLoading}
                        className="page-button small"
                      >
                        Previous
                      </button>
                      
                      <span className="page-info">
                        Page {pagination.page} of {pagination.pages}
                      </span>
                      
                      <button
                        onClick={() => setReviewPage(p => Math.min(p + 1, pagination.pages))}
                        disabled={reviewPage >= pagination.pages || reviewsLoading}
                        className="page-button small"
                      >
                        Next
                      </button>
                    </div>
                  )}
                </>
              ) : (
                <p className="no-reviews">No reviews yet for this product.</p>
              )}
            </div>
          </div>
        </div>
    );
}