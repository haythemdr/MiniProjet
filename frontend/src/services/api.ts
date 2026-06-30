import { Product, ProductDetails } from "@/types/product";

const API_URL = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080";

async function fetchJson<T>(url: string): Promise<T> {
  const response = await fetch(url);

  if (!response.ok) {
    throw new Error("Erreur API");
  }

  return response.json();
}

export async function searchProducts(search: string): Promise<Product[]> {
  const response = await fetch(
    `${API_URL}/products?search=${encodeURIComponent(search)}`
  );

  if (!response.ok) {
    throw new Error("Erreur lors de la recherche");
  }

  return response.json();
}

export async function getProductDetails(url: string): Promise<ProductDetails> {
  return fetchJson<ProductDetails>(
    `${API_URL}/product/details?url=${encodeURIComponent(url)}`
  );
}
