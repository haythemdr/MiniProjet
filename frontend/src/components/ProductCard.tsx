"use client";

import Link from "next/link";
import Image from "next/image";
import { Product } from "@/types/product";

interface Props {
  product: Product;
}

export default function ProductCard({ product }: Props) {
  return (
    <div className="overflow-hidden rounded-xl bg-white shadow-md transition hover:shadow-xl">
      <div className="relative h-52 w-full">
        {product.image && (
          <Image
            src={product.image}
            alt={product.name}
            fill
            sizes="(min-width: 1280px) 25vw, (min-width: 1024px) 33vw, (min-width: 768px) 50vw, 100vw"
            className="object-cover"
          />
        )}
      </div>

      <div className="p-4">
        <h2 className="line-clamp-2 text-lg font-semibold">{product.name}</h2>

        <p className="mt-2 font-bold text-blue-600">{product.price}</p>

        <Link
          href={`/product?url=${encodeURIComponent(product.url)}`}
          className="mt-4 block rounded-lg bg-blue-600 py-2 text-center text-white"
        >
          Voir détails
        </Link>
      </div>
    </div>
  );
}
