"use client";

import Image from "next/image";
import Link from "next/link";
import { Product } from "@/types/product";
import { useState } from "react";


interface Props {
  product: Product;
}

export default function ProductCard({ product }: Props) {
  const [imageFailed, setImageFailed] = useState(false);
  const showImage = product.image && !imageFailed;

  return (
    <article className="group overflow-hidden rounded-lg border border-zinc-200 bg-white shadow-sm transition hover:-translate-y-0.5 hover:shadow-md">
      <div className="relative h-52 w-full border-b border-zinc-100 bg-zinc-50">
        {showImage ? (
          <Image
            src={product.image}
            alt={product.name}
            fill
            unoptimized
            sizes="(min-width: 1280px) 30vw, (min-width: 768px) 50vw, 100vw"
            className="object-contain p-4"
            onError={() => setImageFailed(true)}
          />
        ) : (
          <div className="flex h-full items-center justify-center text-sm text-zinc-400">
            Image indisponible
          </div>
        )}
      </div>

      <div className="flex min-h-48 flex-col p-4">
        <h3 className="line-clamp-2 text-base font-semibold leading-6 text-zinc-900">
          {product.name}
        </h3>
        <div className="mt-2">
          <span
            className={`rounded-full px-3 py-1 text-xs font-semibold ${product.store === "TunisiaNet"
                ? "bg-red-100 text-red-700"
                : "bg-blue-100 text-blue-700"
              }`}
          >
            {product.store}
          </span>
        </div>

        <p className="mt-3 text-xl font-bold text-blue-700">
          {product.price || "Prix indisponible"}
        </p>

        <Link
          href={`/product?url=${encodeURIComponent(product.url)}`}
          className="mt-auto rounded-lg border border-blue-200 bg-blue-50 px-4 py-2 text-center text-sm font-semibold text-blue-700 hover:border-blue-700 hover:bg-blue-700 hover:text-white"
        >
          Voir détails
        </Link>
      </div>
    </article>
  );
}
