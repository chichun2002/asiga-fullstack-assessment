import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { BrowserRouter as Router, Route, Routes} from 'react-router-dom';
import './App.css'
import ProductList from './ProductList';
import ProductDetail from './ProductDetail';
import CreateProduct from './CreateProduct';

const queryClient = new QueryClient();

function App() {
  return (
    <>
      <QueryClientProvider client={queryClient}>
        <Router>
          <Routes>
            <Route path="/" element={<ProductList />} />
            <Route path="/products" element={<ProductList />} />
            <Route path="/products/:id" element={<ProductDetail />} />
            <Route path="/products/create" element={<CreateProduct />} />
          </Routes>
        </Router>
      </QueryClientProvider>
    </>
  )
}

export default App
