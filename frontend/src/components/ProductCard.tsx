"use client";

import Image from "next/image";
import Link from "next/link";
import { Product } from "@/types/product";
import { useState, useCallback } from "react";
import { 
  ShoppingBag, 
  Store, 
  ArrowRight, 
  Eye, 
  Tag,
  ImageIcon
} from "lucide-react";

interface Props {
  product: Product;
}

export default function ProductCard({ product }: Props) {
  const [imageFailed, setImageFailed] = useState(false);
  const [imageLoaded, setImageLoaded] = useState(false);
  const [isHovered, setIsHovered] = useState(false);
  const showImage = product.image && !imageFailed;

  const handleImageLoad = useCallback(() => {
    setImageLoaded(true);
  }, []);

  const handleImageError = useCallback(() => {
    setImageFailed(true);
    setImageLoaded(false);
  }, []);

  return (
    <article
      className="group relative overflow-hidden rounded-2xl border border-zinc-200 bg-white shadow-sm transition-all duration-300 hover:-translate-y-1 hover:shadow-xl hover:shadow-violet-500/5"
      onMouseEnter={() => setIsHovered(true)}
      onMouseLeave={() => setIsHovered(false)}
    >
      {/* Store Badge */}
      <div className="absolute top-3 left-3 z-10">
        <span
          className={`inline-flex items-center gap-1.5 rounded-lg px-3 py-1.5 text-xs font-bold backdrop-blur-sm ${
            product.store === "TunisiaNet"
              ? "bg-red-500/90 text-white shadow-lg shadow-red-500/20"
              : "bg-blue-500/90 text-white shadow-lg shadow-blue-500/20"
          }`}
        >
          <Store className="w-3 h-3" />
          {product.store}
        </span>
      </div>

      {/* Quick View Overlay 
      <div
        className={`absolute inset-0 z-20 flex items-center justify-center bg-black/40 backdrop-blur-sm transition-all duration-300 ${
          isHovered && showImage && imageLoaded ? "opacity-100" : "opacity-0 pointer-events-none"
        }`}
      >
        <div className="flex items-center gap-2 rounded-xl bg-white/90 px-4 py-2 text-sm font-semibold text-zinc-900 shadow-lg backdrop-blur-sm">
          <Eye className="w-4 h-4" />
          Aperçu rapide
        </div>
      </div>*/}

      {/* Image Section */}
      <div className="relative h-56 w-full bg-gradient-to-br from-zinc-50 to-violet-50/30 overflow-hidden">
        {/* Image loading skeleton */}
        {!imageLoaded && showImage && (
          <div className="absolute inset-0 flex items-center justify-center bg-zinc-50">
            <div className="flex flex-col items-center gap-3">
              <div className="w-10 h-10 border-4 border-zinc-200 border-t-violet-600 rounded-full animate-spin"></div>
              <span className="text-xs text-zinc-500 font-medium">Chargement...</span>
            </div>
          </div>
        )}

        {showImage ? (
          <>
            {/* eslint-disable-next-line @next/next/no-img-element */}
            <img
              src={product.image}
              alt={product.name}
              className={`w-full h-full object-contain p-6 transition-all duration-500 ${
                imageLoaded ? 'opacity-100 scale-100' : 'opacity-0 scale-95'
              } group-hover:scale-110`}
              onLoad={handleImageLoad}
              onError={handleImageError}
            />
            {/* Shimmer effect */}
            {imageLoaded && (
              <div className="absolute inset-0 bg-gradient-to-r from-transparent via-white/10 to-transparent -translate-x-full animate-shimmer"></div>
            )}
          </>
        ) : (
          <div className="flex h-full flex-col items-center justify-center gap-3">
            <div className="rounded-xl bg-zinc-100 p-4">
              <ImageIcon className="w-8 h-8 text-zinc-400" />
            </div>
            <span className="text-sm font-medium text-zinc-400">
              Image indisponible
            </span>
          </div>
        )}

        {/* Gradient overlay at bottom */}
        <div className="absolute bottom-0 left-0 right-0 h-20 bg-gradient-to-t from-white/10 to-transparent pointer-events-none"></div>
      </div>

      {/* Content Section */}
      <div className="flex min-h-[200px] flex-col p-5">
        {/* Product Name */}
        <h3 className="line-clamp-2 text-base font-semibold leading-6 text-zinc-900 group-hover:text-violet-700 transition-colors duration-200">
          {product.name}
        </h3>

        {/* Price Section */}
        <div className="mt-3 flex items-center gap-2">
          <div className="flex items-center gap-1.5 rounded-lg bg-gradient-to-r from-blue-50 to-violet-50 px-3 py-1.5 border border-blue-100">
            <Tag className="w-3.5 h-3.5 text-blue-600" />
            <span className="text-lg font-bold bg-gradient-to-r from-blue-600 to-violet-600 bg-clip-text text-transparent">
              {product.price || "Prix indisponible"}
            </span>
          </div>
        </div>

        {/* Divider */}
        <div className="mt-4 mb-4 h-px bg-gradient-to-r from-zinc-200 via-zinc-100 to-transparent"></div>

        {/* Features */}
        <div className="flex flex-wrap gap-2 mb-4">
          <span className="inline-flex items-center gap-1 rounded-md bg-emerald-50 px-2 py-1 text-xs font-medium text-emerald-700 border border-emerald-100">
            <div className="w-1.5 h-1.5 rounded-full bg-emerald-500"></div>
            Disponible
          </span>
          {product.store === "TunisiaNet" && (
            <span className="inline-flex items-center gap-1 rounded-md bg-amber-50 px-2 py-1 text-xs font-medium text-amber-700 border border-amber-100">
              <ShoppingBag className="w-3 h-3" />
              Premium
            </span>
          )}
        </div>

        {/* Action Button */}
        <Link
          href={`/product?url=${encodeURIComponent(product.url)}`}
          className="group/btn relative mt-auto flex items-center justify-center gap-2 overflow-hidden rounded-xl bg-gradient-to-r from-violet-600 to-indigo-600 px-4 py-3 text-sm font-semibold text-white shadow-lg shadow-indigo-500/25 transition-all duration-300 hover:shadow-xl hover:shadow-indigo-500/40 hover:scale-[1.02] active:scale-100"
        >
          <span className="relative z-10 flex items-center gap-2">
            Voir détails
            <ArrowRight className="w-4 h-4 transition-transform duration-300 group-hover/btn:translate-x-0.5" />
          </span>
          {/* Button shine effect */}
          <div className="absolute inset-0 -translate-x-full bg-gradient-to-r from-transparent via-white/20 to-transparent group-hover/btn:translate-x-full transition-transform duration-700"></div>
        </Link>
      </div>

      <style jsx>{`
        @keyframes shimmer {
          0% {
            transform: translateX(-100%);
          }
          100% {
            transform: translateX(100%);
          }
        }
        .animate-shimmer {
          animation: shimmer 2s infinite;
        }
      `}</style>
    </article>
  );
}