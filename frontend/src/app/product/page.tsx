"use client";

import { Suspense, useEffect, useState } from "react";
import { useSearchParams } from "next/navigation";
import Image from "next/image";
import { getProductDetails } from "@/services/api";
import { ProductDetails } from "@/types/product";

function ProductDetailsContent() {
  const searchParams = useSearchParams();
  const url = searchParams.get("url");
  const [product, setProduct] = useState<ProductDetails | null>(null);

  useEffect(() => {
    if (!url) return;

    const loadProduct = async () => {
      try {
        const data = await getProductDetails(url);
        setProduct(data);
      } catch (error) {
        console.error(error);
      }
    };

    loadProduct();
  }, [url]);

  if (!product) {
    return <div className="p-10 text-center">Chargement...</div>;
  }

  return (
    <main className="min-h-screen bg-gray-100">
      <div className="bg-blue-600 p-5 text-white">
        <h1 className="text-3xl font-bold">Détails du produit</h1>
      </div>

      <div className="mx-auto max-w-7xl p-8">
        <button
          onClick={() => window.history.back()}
          className="mb-6 rounded-lg bg-gray-800 px-4 py-2 text-white"
        >
          ← Retour
        </button>

        <div className="rounded-2xl bg-white p-8 shadow-lg">
          <div className="grid gap-10 md:grid-cols-2">
            <div className="relative min-h-96 w-full overflow-hidden rounded-xl border">
              {product.image && (
                <Image
                  src={product.image}
                  alt={product.name}
                  fill
                  sizes="(min-width: 768px) 50vw, 100vw"
                  className="object-contain"
                />
              )}
            </div>

            <div>
              <h1 className="mb-4 text-5xl font-bold">{product.name}</h1>

              <p className="mb-4 text-4xl font-bold text-blue-600">
                {product.price}
              </p>

              {product.availability && (
                <span className="rounded-full bg-green-100 px-4 py-2 text-green-700">
                  {product.availability}
                </span>
              )}
            </div>
          </div>

          <div className="mt-10 border-t pt-8">
            <h2 className="mb-4 text-3xl font-bold">Description du produit</h2>

            <p className="leading-8 text-gray-700">
              {product.description || "Aucune description disponible."}
            </p>
          </div>
        </div>
      </div>
    </main>
  );
}

export default function ProductPage() {
  return (
    <Suspense fallback={<div className="p-10 text-center">Chargement...</div>}>
      <ProductDetailsContent />
    </Suspense>
  );
}
