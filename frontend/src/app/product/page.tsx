"use client";

import Image from "next/image";
import { useSearchParams } from "next/navigation";
import { Suspense, useEffect, useState } from "react";
import { getProductDetails } from "@/services/api";
import { ProductDetails } from "@/types/product";

function ProductDetailsContent() {
  const searchParams = useSearchParams();
  const url = searchParams.get("url");
  const [product, setProduct] = useState<ProductDetails | null>(null);
  const [error, setError] = useState("");
  const [imageFailed, setImageFailed] = useState(false);

  useEffect(() => {
    if (!url) {
      return;
    }

    const loadProduct = async () => {
      try {
        const data = await getProductDetails(url);
        setImageFailed(false);
        setProduct(data);
      } catch {
        setError("Impossible de charger les détails du produit.");
      }
    };

    loadProduct();
  }, [url]);

  const pageError = !url ? "Lien du produit manquant." : error;

  if (pageError) {
    return (
      <main className="flex min-h-screen items-center justify-center bg-[#f5f7fb] p-6">
        <div className="max-w-md rounded-lg border border-zinc-200 bg-white p-8 text-center shadow-sm">
          <p className="mb-5 text-red-600">{pageError}</p>
          <button
            onClick={() => window.history.back()}
            className="rounded-lg bg-zinc-900 px-5 py-2 text-sm font-semibold text-white hover:bg-zinc-700"
          >
            Retour
          </button>
        </div>
      </main>
    );
  }

  if (!product) {
    return (
      <main className="flex min-h-screen items-center justify-center bg-[#f5f7fb]">
        <p className="text-zinc-600">Chargement du produit...</p>
      </main>
    );
  }

  return (
    <main className="min-h-screen bg-[#f5f7fb] text-zinc-900">
      <header className="border-b border-zinc-200 bg-white">
        <div className="mx-auto flex max-w-6xl items-center justify-between gap-4 px-6 py-5">
          <div>
            <h1 className="mt-1 text-2xl font-bold">Détails du produit</h1>
          </div>

          <button
            onClick={() => window.history.back()}
            className="rounded-lg border border-zinc-300 bg-white px-4 py-2 text-sm font-semibold text-zinc-700 hover:border-blue-300 hover:text-blue-700"
          >
            Retour
          </button>
        </div>
      </header>

      <div className="mx-auto max-w-6xl px-6 py-8">
        <section className="grid gap-8 rounded-lg border border-zinc-200 bg-white p-6 shadow-sm lg:grid-cols-[minmax(0,1fr)_400px] lg:p-8">
          <div className="flex min-h-[420px] items-center justify-center rounded-lg border border-zinc-200 bg-zinc-50 p-6">
            <div className="relative h-80 w-full max-w-xl sm:h-96">
              {product.image && !imageFailed ? (
                <Image
                  src={product.image}
                  alt={product.name}
                  fill
                  priority
                  unoptimized
                  sizes="(min-width: 1024px) 55vw, 100vw"
                  className="object-contain"
                  onError={() => setImageFailed(true)}
                />
              ) : (
                <div className="flex h-full items-center justify-center text-zinc-400">
                  Image indisponible
                </div>
              )}
            </div>
          </div>

          <aside className="flex flex-col justify-center">
            <span className="mb-4 w-fit rounded-full bg-emerald-50 px-3 py-1 text-sm font-semibold text-emerald-700">
              Produit collecté
            </span>

            <h2 className="text-3xl font-bold leading-tight">
              {product.name}
            </h2>

            <div className="my-6 rounded-lg border border-blue-100 bg-blue-50 p-5">
              <p className="text-sm font-semibold uppercase text-blue-800">
                Prix
              </p>
              <p className="mt-2 text-4xl font-bold text-blue-700">
                {product.price || "Prix indisponible"}
              </p>
            </div>

            <div className="rounded-lg border border-zinc-200 bg-white p-4">
              <p className="text-sm font-semibold text-zinc-500">
                Disponibilité
              </p>
              <p className="mt-1 font-semibold text-emerald-700">
                {product.availability || "Non précisée"}
              </p>
            </div>
          </aside>
        </section>

        <section className="mt-6 rounded-lg border border-zinc-200 bg-white p-6 shadow-sm lg:p-8">
          <h3 className="mb-4 text-xl font-bold">Description</h3>
          <p className="whitespace-pre-line leading-8 text-zinc-700">
            {product.description || "Aucune description disponible."}
          </p>
        </section>
      </div>
    </main>
  );
}

export default function ProductPage() {
  return (
    <Suspense
      fallback={
        <main className="flex min-h-screen items-center justify-center bg-[#f5f7fb]">
          <p className="text-zinc-600">Chargement du produit...</p>
        </main>
      }
    >
      <ProductDetailsContent />
    </Suspense>
  );
}
