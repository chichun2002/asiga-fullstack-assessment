import { keepPreviousData, useQuery } from "@tanstack/react-query";
import axios from "axios";
import { useCallback, useState } from "react";
import { useNavigate } from "react-router-dom";
import "./ProductList.css";

interface Product {
    ID: number;
    name: string;
    price: number;
}

interface ProductResponse {
    products: Product[];
    pagination: {
        total: number;
        page: number;
        limit: number;
        pages: number;
    };
}

type SortField = "name" | "price" | "created_at";
type SortOrder = "asc" | "desc"

export default function ProductList() {
    const navigate = useNavigate();
    const [page, setPage] = useState(1);
    const [limit, setLimit] = useState(12);
    const [search, setSearch] = useState("");
    const [searchInput, setSearchInput] = useState("");
    const [sortBy, setSortBy] = useState<SortField>("created_at");
    const [sortOrder, setSortOrder] = useState<SortOrder>("desc");

    // The explicit generic type parameters for useQuery are important
    const { data, isLoading, error, refetch } = useQuery<ProductResponse, Error>({
        queryKey: ["products", page, limit, sortBy, sortOrder, search],
        queryFn: async () => {
            console.log("Fetching with params:", { page, limit, sort: sortBy, order: sortOrder, search: search });
            const response = await axios.get<ProductResponse>("http://localhost:8080/products", {
                params: {
                    page,
                    limit,
                    sort: sortBy,
                    order: sortOrder,
                    search: search,
                }
            });
            return response.data;
        },
        placeholderData: keepPreviousData,
    });

    const handleSortChange = useCallback((field: SortField) => {
        if (sortBy === field) {
            setSortOrder(prevOrder => prevOrder === "asc" ? "desc" : "asc");
        } else {
            setSortBy(field);
            setSortOrder("desc");
        }
        
        setPage(1);
    }, [sortBy, sortOrder]);

    const getSortIcon = useCallback((field: SortField) => {
        if (sortBy !== field) return null;
        return sortOrder === "asc" ? "↑" : "↓";
    }, [sortBy, sortOrder]);

    const handleViewDetails = (productId: number) => {
        navigate(`/products/${productId}`);
    };
    
    const handleCreateProduct= () => {
        navigate(`/products/create`);
    };

    const handleDeleteProduct = (productId: number) => {
        axios
            .delete(`http://localhost:8080/products/${productId}`)
            .then(() => {
                refetch();
            })
            .catch((error) => {
                console.error(error)
            });
    };

    const handleSearch = (e: React.FormEvent) => {
        e.preventDefault();
        setSearch(searchInput);
        setPage(1); // Reset to first page when searching
    };

    if (isLoading) return <div>Loading...</div>;
    if (error) return <div>Error: {error.message}</div>;

    // Now TypeScript knows that data is of type ProductResponse
    const products = data?.products || [];
    const isEmpty = products.length === 0;

    return (
        <div className="product-container">
            <h1>Product List</h1>
            <button
                className="create-product-button"
                onClick={handleCreateProduct}
            >
                + Add New Product
            </button>

            <form onSubmit={handleSearch} className="search-form">
                <input
                    type="text"
                    placeholder="Search products..."
                    value={searchInput}
                    onChange={(e) => setSearchInput(e.target.value)}
                    className="search-input"
                />
                <button type="submit" className="search-button">
                    Search
                </button>
                {search && (
                    <button 
                        type="button" 
                        className="clear-search-button"
                        onClick={() => {
                            setSearchInput("");
                            setSearch("");
                            setPage(1);
                        }}
                    >
                        Clear
                    </button>
                )}
            </form>

            <div className="sort-controls">
                <span>Sort by:</span>
                <button
                    className={`sort-button ${sortBy === "name" ? "active" : ""}`}
                    onClick={() => handleSortChange("name")}
                >
                    Name {getSortIcon("name")}
                </button>
                <button
                    className={`sort-button ${sortBy === "price" ? "active" : ""}`}
                    onClick={() => handleSortChange("price")}
                >
                    Price {getSortIcon("price")}
                </button>
                <button
                    className={`sort-button ${sortBy === "created_at" ? "active" : ""}`}
                    onClick={() => handleSortChange("created_at")}
                >
                    Date {getSortIcon("created_at")}
                </button>
            </div>

            <div className="product-grid">
                {!isEmpty ? (
                    products.map((product) => (
                        <div key={`product-${product.ID}`} className="product-card">
                            <h3>{product.name}</h3>
                            <button
                                onClick={() => handleDeleteProduct(product.ID)}
                                style={{
                                    position: 'absolute',
                                    top: '10px',
                                    right: '10px',
                                }}
                            >
                                X
                            </button>
                            <p className="product-price">${product.price.toFixed(2)}</p>
                            <button 
                                className="view-button"
                                onClick={() => handleViewDetails(product.ID)}
                            >
                                View Details
                            </button>
                        </div>
                    ))
                ) : (
                    <div className="no-products">No products found</div>
                )}
            </div>
            
            {data?.pagination && (
                <div className="pagination-controls">
                    <button
                        onClick={() => setPage(page => Math.max(page - 1, 1))}
                        disabled={page === 1}
                        className="page-button"
                    >
                        Previous
                    </button>

                    <span className="page-info">
                        Page {data.pagination.page} of {data.pagination.pages}
                    </span>

                    <button
                        onClick={() => setPage(page => page < data.pagination.pages ? page + 1 : page)}
                        disabled={page >= data.pagination.pages}
                        className="page-button"
                    >
                        Next
                    </button>
                </div>
            )}
        </div>
    );
}