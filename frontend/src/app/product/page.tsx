"use client";

import { useSearchParams } from "next/navigation";
import { Suspense, useEffect, useState, useRef, useCallback } from "react";
import { getProductDetails } from "@/services/api";
import { ProductDetails } from "@/types/product";
import { 
  ArrowLeft, 
  ShoppingBag, 
  Tag, 
  Package, 
  CheckCircle2, 
  AlertCircle,
  Shield,
  Truck,
  RefreshCw,
  Star,
  Clock,
  Award,
  ChevronRight,
  ImageIcon,
  Sparkles
} from "lucide-react";

function ProductDetailsContent() {
  const searchParams = useSearchParams();
  const url = searchParams.get("url");
  const [product, setProduct] = useState<ProductDetails | null>(null);
  const [error, setError] = useState("");
  const [imageFailed, setImageFailed] = useState(false);
  const [imageLoaded, setImageLoaded] = useState(false);
  const imgRef = useRef<HTMLImageElement>(null);

  // Check if image is already loaded (e.g., from cache)
  useEffect(() => {
    if (!product?.image || imageFailed) return;

    const img = imgRef.current;
    if (img) {
      // Check if image is already cached and loaded
      if (img.complete && img.naturalWidth > 0) {
        setImageLoaded(true);
      }
    }
  }, [product?.image, imageFailed]);

  const handleImageLoad = useCallback(() => {
    setImageLoaded(true);
  }, []);

  const handleImageError = useCallback(() => {
    setImageFailed(true);
    setImageLoaded(false);
  }, []);

  useEffect(() => {
    if (!url) {
      return;
    }

    const loadProduct = async () => {
      try {
        setImageLoaded(false);
        setImageFailed(false);
        const data = await getProductDetails(url);
        setProduct(data);
      } catch {
        setError("Impossible de charger les détails du produit.");
      }
    };

    loadProduct();
  }, [url]);

  const pageError = !url ? "Lien du produit manquant." : error;

  // ... rest of the component remains the same
  if (pageError) {
    return (
      <main className="flex min-h-screen items-center justify-center bg-gradient-to-br from-slate-50 via-white to-violet-50/30 p-6">
        <div className="relative">
          <div className="absolute -inset-1 rounded-2xl bg-gradient-to-r from-red-400 to-rose-400 opacity-25 blur"></div>
          <div className="relative max-w-md rounded-2xl border border-zinc-200 bg-white p-8 text-center shadow-xl">
            <div className="mb-6 flex justify-center">
              <div className="rounded-full bg-red-50 p-4">
                <AlertCircle className="w-10 h-10 text-red-500" />
              </div>
            </div>
            <p className="mb-2 text-lg font-semibold text-zinc-900">Une erreur est survenue</p>
            <p className="mb-6 text-sm text-red-600">{pageError}</p>
            <button
              onClick={() => window.history.back()}
              className="group inline-flex items-center gap-2 rounded-xl bg-gradient-to-r from-violet-600 to-indigo-600 px-6 py-3 text-sm font-semibold text-white shadow-lg shadow-indigo-500/25 transition-all hover:shadow-xl hover:shadow-indigo-500/30 hover:scale-105 active:scale-100"
            >
              <ArrowLeft className="w-4 h-4 transition-transform group-hover:-translate-x-0.5" />
              Retour
            </button>
          </div>
        </div>
      </main>
    );
  }

  if (!product) {
    return (
      <main className="flex min-h-screen items-center justify-center bg-gradient-to-br from-slate-50 via-white to-violet-50/30">
        <div className="flex flex-col items-center gap-6">
          <div className="relative">
            <div className="w-20 h-20 rounded-2xl border-4 border-zinc-100 border-t-violet-600 animate-spin"></div>
            <Package className="absolute inset-0 m-auto w-8 h-8 text-violet-600" />
          </div>
          <div className="text-center space-y-2">
            <p className="text-xl font-semibold text-zinc-700">Chargement du produit</p>
            <p className="text-sm text-zinc-500">Récupération des détails en cours...</p>
          </div>
        </div>
      </main>
    );
  }

  return (
    <main className="min-h-screen bg-gradient-to-br from-slate-50 via-white to-violet-50/30">
      {/* Background pattern */}
      <div className="fixed inset-0 -z-10">
        <div className="absolute inset-0 bg-[linear-gradient(to_right,#80808012_1px,transparent_1px),linear-gradient(to_bottom,#80808012_1px,transparent_1px)] bg-[size:24px_24px]"></div>
        <div className="absolute top-0 left-0 right-0 h-96 bg-gradient-to-br from-violet-500/5 via-transparent to-rose-500/5 blur-3xl"></div>
      </div>

      {/* Header */}
      <header className="relative border-b border-zinc-200/60 bg-white/80 backdrop-blur-xl">
        <div className="mx-auto flex max-w-7xl items-center justify-between gap-4 px-6 py-6">
          <div className="flex items-center gap-4">
            <button
              onClick={() => window.history.back()}
              className="group flex items-center gap-2 rounded-xl border border-zinc-200 bg-white px-4 py-2.5 text-sm font-medium text-zinc-700 transition-all hover:border-zinc-300 hover:bg-zinc-50 hover:shadow-md"
            >
              <ArrowLeft className="w-4 h-4 transition-transform group-hover:-translate-x-0.5" />
              Retour
            </button>
            <div className="hidden sm:block">
              <div className="flex items-center gap-2 text-sm text-zinc-500">
                <ShoppingBag className="w-4 h-4" />
                <span>Produits</span>
                <ChevronRight className="w-4 h-4" />
                <span className="font-medium text-zinc-700">Détails</span>
              </div>
              <h1 className="text-2xl font-bold text-zinc-900">Détails du produit</h1>
            </div>
          </div>

          <div className="flex items-center gap-3">
            <div className="hidden md:flex items-center gap-2 rounded-xl border border-green-200 bg-green-50 px-3 py-2">
              <Shield className="w-4 h-4 text-green-600" />
              <span className="text-xs font-semibold text-green-700">Achat sécurisé</span>
            </div>
            <div className="hidden md:flex items-center gap-2 rounded-xl border border-blue-200 bg-blue-50 px-3 py-2">
              <Truck className="w-4 h-4 text-blue-600" />
              <span className="text-xs font-semibold text-blue-700">Livraison disponible</span>
            </div>
          </div>
        </div>
      </header>

      <div className="mx-auto max-w-7xl px-6 py-8">
        {/* Main Product Section */}
        <section className="relative overflow-hidden rounded-2xl border border-zinc-200 bg-white shadow-xl">
          <div className="absolute top-0 right-0 w-96 h-96 bg-gradient-to-bl from-violet-500/5 to-transparent rounded-full blur-3xl"></div>
          
          <div className="relative grid gap-8 p-6 lg:grid-cols-[minmax(0,1fr)_420px] lg:p-10">
            {/* Image Section */}
            <div className="flex min-h-[420px] items-center justify-center rounded-2xl bg-gradient-to-br from-zinc-50 to-violet-50/30 p-8">
              <div className="relative h-80 w-full max-w-xl sm:h-96">
                {/* Image loading skeleton */}
                {!imageLoaded && product.image && !imageFailed && (
                  <div className="absolute inset-0 flex items-center justify-center bg-zinc-50 rounded-xl z-10">
                    <div className="flex flex-col items-center gap-3">
                      <div className="w-12 h-12 border-4 border-zinc-200 border-t-violet-600 rounded-full animate-spin"></div>
                      <span className="text-sm text-zinc-500 font-medium">Chargement de l'image...</span>
                    </div>
                  </div>
                )}
                {product.image && !imageFailed ? (
                  <>
                    {/* eslint-disable-next-line @next/next/no-img-element */}
                    <img
                      ref={imgRef}
                      src={product.image}
                      alt={product.name}
                      className={`w-full h-full object-contain transition-opacity duration-500 ${
                        imageLoaded ? 'opacity-100' : 'opacity-0'
                      }`}
                      onLoad={handleImageLoad}
                      onError={handleImageError}
                    />
                  </>
                ) : (
                  <div className="flex h-full flex-col items-center justify-center gap-3 text-zinc-400">
                    <div className="rounded-2xl bg-zinc-100 p-6">
                      <ImageIcon className="w-16 h-16" />
                    </div>
                    <p className="text-sm font-medium">Image indisponible</p>
                  </div>
                )}
              </div>
            </div>

            {/* Product Info Section */}
            <aside className="flex flex-col justify-center">
              <div className="flex items-center gap-2 mb-6">
                <span className="inline-flex items-center gap-1.5 rounded-full bg-gradient-to-r from-emerald-50 to-green-50 px-4 py-1.5 text-sm font-semibold text-emerald-700 border border-emerald-200">
                  <CheckCircle2 className="w-4 h-4" />
                  Produit collecté
                </span>
                {product.availability && (
                  <span className="inline-flex items-center gap-1.5 rounded-full bg-blue-50 px-4 py-1.5 text-sm font-semibold text-blue-700 border border-blue-200">
                    <Clock className="w-4 h-4" />
                    {product.availability}
                  </span>
                )}
              </div>

              <h2 className="text-3xl font-bold leading-tight text-zinc-900 lg:text-4xl">
                {product.name}
              </h2>

              <div className="relative mt-8">
                <div className="absolute -inset-0.5 rounded-2xl bg-gradient-to-r from-blue-400 to-violet-400 opacity-20 blur"></div>
                <div className="relative rounded-2xl border-2 border-blue-100 bg-gradient-to-br from-blue-50 to-indigo-50 p-6">
                  <div className="flex items-center gap-2 mb-1">
                    <Tag className="w-4 h-4 text-blue-600" />
                    <p className="text-sm font-semibold uppercase tracking-wider text-blue-600">
                      Prix
                    </p>
                  </div>
                  <p className="mt-2 text-5xl font-bold bg-gradient-to-r from-blue-600 to-violet-600 bg-clip-text text-transparent">
                    {product.price || "Prix indisponible"}
                  </p>
                  <div className="mt-3 flex items-center gap-2 text-xs text-blue-600">
                    <RefreshCw className="w-3 h-3" />
                    <span>Prix mis à jour</span>
                  </div>
                </div>
              </div>

              <div className="mt-4 rounded-2xl border border-zinc-200 bg-white p-5">
                <div className="flex items-center gap-2 mb-1">
                  <Package className="w-4 h-4 text-zinc-500" />
                  <p className="text-sm font-semibold text-zinc-500">
                    Disponibilité
                  </p>
                </div>
                <p className="mt-2 text-lg font-bold text-emerald-600">
                  {product.availability || "Non précisée"}
                </p>
              </div>

              <div className="mt-6 grid grid-cols-2 gap-3">
                <div className="flex items-center gap-2 rounded-xl border border-zinc-200 bg-zinc-50 p-3">
                  <Shield className="w-4 h-4 text-green-600 flex-shrink-0" />
                  <span className="text-xs font-medium text-zinc-600">Achat sécurisé</span>
                </div>
                <div className="flex items-center gap-2 rounded-xl border border-zinc-200 bg-zinc-50 p-3">
                  <Truck className="w-4 h-4 text-blue-600 flex-shrink-0" />
                  <span className="text-xs font-medium text-zinc-600">Livraison rapide</span>
                </div>
                <div className="flex items-center gap-2 rounded-xl border border-zinc-200 bg-zinc-50 p-3">
                  <Award className="w-4 h-4 text-amber-600 flex-shrink-0" />
                  <span className="text-xs font-medium text-zinc-600">Qualité garantie</span>
                </div>
                <div className="flex items-center gap-2 rounded-xl border border-zinc-200 bg-zinc-50 p-3">
                  <Star className="w-4 h-4 text-violet-600 flex-shrink-0" />
                  <span className="text-xs font-medium text-zinc-600">Top marques</span>
                </div>
              </div>
            </aside>
          </div>
        </section>

        {/* Description Section */}
        <section className="relative mt-6 overflow-hidden rounded-2xl border border-zinc-200 bg-white shadow-xl">
          <div className="absolute top-0 left-0 w-96 h-96 bg-gradient-to-br from-violet-500/5 to-transparent rounded-full blur-3xl"></div>
          
          <div className="relative p-6 lg:p-10">
            <div className="flex items-center gap-3 mb-6">
              <div className="p-2 rounded-xl bg-violet-100">
                <Sparkles className="w-5 h-5 text-violet-600" />
              </div>
              <h3 className="text-2xl font-bold text-zinc-900">Description du produit</h3>
            </div>
            
            {product.description ? (
              <div className="prose prose-zinc max-w-none">
                <div className="whitespace-pre-line leading-8 text-zinc-700 bg-zinc-50 rounded-xl p-6 border border-zinc-200">
                  {product.description}
                </div>
              </div>
            ) : (
              <div className="rounded-xl border-2 border-dashed border-zinc-200 bg-zinc-50 p-12 text-center">
                <div className="flex flex-col items-center gap-3">
                  <AlertCircle className="w-10 h-10 text-zinc-400" />
                  <p className="text-zinc-500 font-medium">Aucune description disponible pour ce produit.</p>
                </div>
              </div>
            )}
          </div>
        </section>

        {/* Additional Info Footer */}
        <div className="mt-6 grid grid-cols-1 sm:grid-cols-3 gap-4">
          <div className="rounded-xl border border-zinc-200 bg-white/80 backdrop-blur p-5 text-center">
            <div className="flex justify-center mb-3">
              <div className="p-2 rounded-lg bg-green-100">
                <Shield className="w-5 h-5 text-green-600" />
              </div>
            </div>
            <p className="text-sm font-semibold text-zinc-700">Paiement sécurisé</p>
            <p className="text-xs text-zinc-500 mt-1">Vos données sont protégées</p>
          </div>
          <div className="rounded-xl border border-zinc-200 bg-white/80 backdrop-blur p-5 text-center">
            <div className="flex justify-center mb-3">
              <div className="p-2 rounded-lg bg-blue-100">
                <Truck className="w-5 h-5 text-blue-600" />
              </div>
            </div>
            <p className="text-sm font-semibold text-zinc-700">Livraison en Tunisie</p>
            <p className="text-xs text-zinc-500 mt-1">Rapide et fiable</p>
          </div>
          <div className="rounded-xl border border-zinc-200 bg-white/80 backdrop-blur p-5 text-center">
            <div className="flex justify-center mb-3">
              <div className="p-2 rounded-lg bg-amber-100">
                <RefreshCw className="w-5 h-5 text-amber-600" />
              </div>
            </div>
            <p className="text-sm font-semibold text-zinc-700">Retour facile</p>
            <p className="text-xs text-zinc-500 mt-1">Satisfait ou remboursé</p>
          </div>
        </div>
      </div>
    </main>
  );
}

export default function ProductPage() {
  return (
    <Suspense
      fallback={
        <main className="flex min-h-screen items-center justify-center bg-gradient-to-br from-slate-50 via-white to-violet-50/30">
          <div className="flex flex-col items-center gap-4">
            <div className="relative">
              <div className="w-16 h-16 rounded-2xl border-4 border-zinc-100 border-t-violet-600 animate-spin"></div>
              <Package className="absolute inset-0 m-auto w-6 h-6 text-violet-600" />
            </div>
            <p className="text-lg font-semibold text-zinc-700">Chargement...</p>
          </div>
        </main>
      }
    >
      <ProductDetailsContent />
    </Suspense>
  );
}