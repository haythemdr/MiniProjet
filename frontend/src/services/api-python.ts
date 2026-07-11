import { Product } from "@/types/product";

const API_URL = "http://127.0.0.1:8000";

export interface SearchResponse {
  store: string;
  source: string;
  lastUpdated: string;
  products: Product[];
}

export function streamProducts(
  search: string,
  onProducts: (response: SearchResponse) => void,
  onDone?: () => void,
  onError?: () => void
) {
  const eventSource = new EventSource(
    `${API_URL}/products/?query=${encodeURIComponent(search)}`
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
}