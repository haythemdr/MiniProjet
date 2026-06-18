import { Product, ProductDetails } from "@/types/product";

const API_URL = "http://localhost:8080";

export async function searchProducts(
  search: string
): Promise<Product[]> {
  const response = await fetch(
    `${API_URL}/products?search=${encodeURIComponent(search)}`
  );

  return response.json();
}

export async function getProductDetails(
  url: string
): Promise<ProductDetails> {
  const response = await fetch(
    `${API_URL}/product/details?url=${encodeURIComponent(url)}`
  );

  return response.json();
}