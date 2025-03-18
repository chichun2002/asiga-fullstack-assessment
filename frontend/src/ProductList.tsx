import { useQuery } from "@tanstack/react-query";
import axios from "axios";

interface Product {
    id: number;
    name: string;
    price: number;
}

export default function ProductList() {
    const { data, isLoading, error } = useQuery({
        queryKey: ["products"],
        queryFn: () => 
            axios.get("http://localhost:8080/products")
                .then(res => res.data)
                .catch(err=> {
                    throw new Error("Failed to fetch products: " + err.message);
                }),
    });

    if (isLoading) return <div>Loading...</div>;
    if (error) return <div>Error: {error.message}</div>

    return (
        <div>
            <h1>Product List</h1>
            {data?.map((product: Product) => (
                <div key={product.id}>
                    <h3>{product.name}</h3>
                    <p>${product.price}</p>
                </div>
            ))}
        </div>
    );
}