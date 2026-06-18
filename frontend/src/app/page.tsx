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
  const [loading, setLoading] = useState(false);

  const loadProducts = async (query: string, displayValue = query) => {
    if (!query.trim()) return;

    setLoading(true);

    try {
      const data = await searchProducts(query);
      setProducts(Array.isArray(data) ? data : []);
      setSearch(displayValue);
    } catch (error) {
      console.error(error);
      setProducts([]);
    } finally {
      setLoading(false);
    }
  };

  return (
    <main className="min-h-screen bg-gray-100">
      <div className="bg-blue-600 py-6 text-white shadow">
        <h1 className="text-center text-5xl font-extrabold">
          Tunisianet Explorer
        </h1>

        <p className="mt-2 text-center text-blue-100">
          Rechercher et explorer les produits de Tunisianet
        </p>
      </div>

      <div className="mx-auto max-w-7xl p-8">
        <div className="mx-auto mb-10 flex max-w-3xl gap-4">
          <input
            type="text"
            placeholder="Rechercher un produit..."
            value={search}
            onChange={(event) => setSearch(event.target.value)}
            onKeyDown={(event) => {
              if (event.key === "Enter") {
                loadProducts(search);
              }
            }}
            className="flex-1 rounded-xl border bg-white px-4 py-3"
          />

          <button
            onClick={() => loadProducts(search)}
            className="rounded-xl bg-blue-600 px-8 text-white hover:bg-blue-700"
          >
            Rechercher
          </button>
        </div>

        <div className="flex gap-8">
          <aside className="h-fit w-72 rounded-xl bg-white p-5 shadow">
            <h2 className="mb-5 text-xl font-bold">Catégories</h2>

            {Object.entries(categories).map(([mainCategory, subCategories]) => (
              <div key={mainCategory} className="mb-6">
                <h3 className="mb-2 font-bold text-blue-600">
                  {mainCategory}
                </h3>

                <div className="flex flex-col gap-2">
                  {subCategories.map((subCategory) => (
                    <button
                      key={subCategory.label}
                      onClick={() =>
                        loadProducts(subCategory.query, subCategory.label)
                      }
                      className="text-left hover:text-blue-600 hover:underline"
                    >
                      • {subCategory.label}
                    </button>
                  ))}
                </div>
              </div>
            ))}
          </aside>

          <section className="flex-1">
            {loading && <div className="py-10 text-center">Chargement...</div>}

            {!loading && products.length > 0 && (
              <p className="mb-4 font-semibold text-gray-700">
                {products.length} produits trouvés
              </p>
            )}

            {!loading && products.length === 0 && (
              <div className="rounded-xl bg-white p-10 text-center text-gray-500 shadow">
                Recherchez un produit ou choisissez une catégorie.
              </div>
            )}

            <div className="grid gap-6 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4">
              {products.map((product) => (
                <ProductCard
                  key={product.url || product.name}
                  product={product}
                />
              ))}
            </div>
          </section>
        </div>
      </div>
    </main>
  );
}
