import { Product, ProductDetails } from "@/types/product";

const API_URL = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080";

export interface SearchResponse {
  store: string;
  source: string;
  lastUpdated: string;
  products: Product[];
}

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

export async function getProductDetails(
  url: string
): Promise<ProductDetails> {
  return fetchJson<ProductDetails>(
    `${API_URL}/product/details?url=${encodeURIComponent(url)}`
  );
}

export function streamProducts(
  search: string,
  onProducts: (response: SearchResponse) => void,
  onDone?: () => void,
  onError?: () => void
) {
  const eventSource = new EventSource(
    `${API_URL}/products/stream?search=${encodeURIComponent(search)}`
  );

  eventSource.onmessage = (event) => {
    const response: SearchResponse = JSON.parse(event.data);
    onProducts(response);
  };

  eventSource.addEventListener("done", () => {
    eventSource.close();
    onDone?.();
  });

  eventSource.onerror = () => {
    eventSource.close();
    onError?.();
  };

  return eventSource;
}/*
export async function getSuggestions(query: string): Promise<string[]> {
  const response = await fetch(
    `${API_URL}/elastic/suggest?q=${encodeURIComponent(query)}`
  );

  if (!response.ok) {
    throw new Error("Suggestion error");
  }

  return response.json();
}*/
export async function getSuggestions(query: string): Promise<string[]> {
  return [];
}