"use client";

import { useState } from "react";
import ProductCard from "@/components/ProductCard";
import { searchProducts } from "@/services/api";
import { Product } from "@/types/product";

type CategoryItem = {
  label: string;
  query: string;
};

const categories: Record<string, CategoryItem[]> = {
  Informatique: [
    { label: "Composants", query: "composant pc" },
    { label: "Ordinateurs", query: "ordinateur" },
    { label: "Réseaux et connectivité", query: "routeur" },
    { label: "Périphériques", query: "clavier souris" },
    { label: "Stockages", query: "disque dur ssd" },
  ],
  Téléphonie: [
    { label: "Téléphones portables", query: "téléphone portable" },
    { label: "Smartphones", query: "smartphone" },
    { label: "Accessoires", query: "accessoire téléphone" },
    { label: "Téléphones fixes", query: "téléphone fixe" },
    { label: "Smartwatch", query: "smartwatch" },
  ],
  "Santé & Beauté": [
    { label: "Toiletries", query: "gel douche" },
    { label: "Moniteurs de santé", query: "tensiomètre" },
    { label: "Bébé & enfants", query: "bébé" },
    { label: "Pharmaceutiques & médicaments", query: "thermomètre" },
    { label: "Produits pour soins personnels", query: "rasoir" },
  ],
  Électroménager: [
    { label: "Aspirateurs", query: "aspirateur" },
    { label: "Machine à laver", query: "machine à laver" },
    { label: "Sèche-linge", query: "sèche linge" },
    { label: "Lave-vaisselle", query: "lave vaisselle" },
    { label: "Fours", query: "four" },
  ],
};

export default function Home() {
  const [search, setSearch] = useState("");
  const [products, setProducts] = useState<Product[]>([]);
  const [currentPage, setCurrentPage] = useState(1);

  const PRODUCTS_PER_PAGE = 24;
  const [loading, setLoading] = useState(false);
  const [message, setMessage] = useState(
    "Recherchez un produit ou choisissez une catégorie."
  );

  const loadProducts = async (query: string, displayValue = query) => {
    const value = query.trim();

    if (!value) {
      setMessage("Saisissez un mot clé pour lancer la recherche.");
      return;
    }

    setLoading(true);
    setSearch(displayValue);
    setMessage("");

    try {
      const data = await searchProducts(value);
      setProducts(data);
      setCurrentPage(1);

      if (data.length === 0) {
        setMessage("Aucun produit trouvé pour cette recherche.");
      }
    } catch {
      setProducts([]);
      setMessage("Impossible de charger les produits pour le moment.");
    } finally {
      setLoading(false);
    }
  };
  const indexOfLastProduct = currentPage * PRODUCTS_PER_PAGE;
  const indexOfFirstProduct = indexOfLastProduct - PRODUCTS_PER_PAGE;

  const currentProducts = products.slice(
    indexOfFirstProduct,
    indexOfLastProduct
  );

  const totalPages = Math.ceil(products.length / PRODUCTS_PER_PAGE);
  return (
    <main className="min-h-screen bg-[#f5f7fb] text-zinc-900">
      <header className="border-b border-zinc-200 bg-white">
        <div className="mx-auto max-w-7xl px-6 py-6">
          <div className="mt-2 flex flex-col justify-between gap-3 md:flex-row md:items-end">
            <div>
              <h1 className="text-3xl font-bold tracking-normal md:text-4xl">
                Tunisianet Explorer
              </h1>
            </div>
          </div>
        </div>
      </header>

      <div className="mx-auto max-w-7xl px-6 py-8">
        <div className="mb-8 rounded-lg border border-zinc-200 bg-white p-4 shadow-sm">
          <div className="flex flex-col gap-3 sm:flex-row">
            <input
              type="text"
              placeholder="Exemple : ordinateur, smartphone, aspirateur..."
              value={search}
              onChange={(event) => setSearch(event.target.value)}
              onKeyDown={(event) => {
                if (event.key === "Enter") {
                  loadProducts(search);
                }
              }}
              className="min-h-12 flex-1 rounded-lg border border-zinc-300 bg-white px-4 text-sm outline-none focus:border-blue-600 focus:ring-2 focus:ring-blue-100"
            />

            <button
              onClick={() => loadProducts(search)}
              className="min-h-12 rounded-lg bg-blue-700 px-8 text-sm font-semibold text-white hover:bg-blue-800"
            >
              Rechercher
            </button>
          </div>
        </div>

        <div className="grid gap-8 lg:grid-cols-[280px_minmax(0,1fr)]">
          <aside className="h-fit rounded-lg border border-zinc-200 bg-white p-5 shadow-sm">
            <h2 className="mb-5 text-lg font-bold">Catégories</h2>

            {Object.entries(categories).map(([mainCategory, subCategories]) => (
              <div key={mainCategory} className="mb-6 last:mb-0">
                <h3 className="mb-3 text-sm font-bold uppercase text-zinc-500">
                  {mainCategory}
                </h3>

                <div className="flex flex-col gap-1">
                  {subCategories.map((subCategory) => (
                    <button
                      key={subCategory.label}
                      onClick={() =>
                        loadProducts(subCategory.query, subCategory.label)
                      }
                      className="rounded-md px-3 py-2 text-left text-sm text-zinc-700 hover:bg-blue-50 hover:text-blue-700"
                    >
                      {subCategory.label}
                    </button>
                  ))}
                </div>
              </div>
            ))}
          </aside>

          <section>
            <div className="mb-4 flex min-h-8 items-center justify-between">
              <h2 className="text-xl font-bold">Produits</h2>
              {!loading && products.length > 0 && (
                <p className="text-sm font-medium text-zinc-600">
                  {products.length} produits trouvés
                </p>
              )}
            </div>

            {loading && (
              <div className="rounded-lg border border-zinc-200 bg-white p-10 text-center text-zinc-600 shadow-sm">
                Chargement des produits...
              </div>
            )}

            {!loading && products.length === 0 && (
              <div className="rounded-lg border border-dashed border-zinc-300 bg-white p-10 text-center text-zinc-500">
                {message}
              </div>
            )}

            {!loading && products.length > 0 && (
              <div className="grid gap-5 sm:grid-cols-2 xl:grid-cols-3">
                {currentProducts.map((product) => (
                  <ProductCard
                    key={product.url || product.name}
                    product={product}
                  />
                ))}
              </div>
            )}
            
            {!loading && products.length > 0 && totalPages > 1 && (
              <div className="mt-8 flex justify-center items-center gap-2">
                <button
                  disabled={currentPage === 1}
                  onClick={() => setCurrentPage(currentPage - 1)}
                  className="px-4 py-2 rounded bg-gray-200 disabled:opacity-50"
                >
                  Previous
                </button>

                {Array.from({ length: totalPages }, (_, i) => (
                  <button
                    key={i}
                    onClick={() => setCurrentPage(i + 1)}
                    className={`px-4 py-2 rounded ${currentPage === i + 1
                        ? "bg-blue-700 text-white"
                        : "bg-gray-200 hover:bg-gray-300"
                      }`}
                  >
                    {i + 1}
                  </button>
                ))}

                <button
                  disabled={currentPage === totalPages}
                  onClick={() => setCurrentPage(currentPage + 1)}
                  className="px-4 py-2 rounded bg-gray-200 disabled:opacity-50"
                >
                  Next
                </button>
              </div>
            )}
          </section>

        </div>
      </div>
    </main>
  );
}
