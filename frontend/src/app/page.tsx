"use client";

import { useState, useEffect, useRef } from "react";
import ProductCard from "@/components/ProductCard";
import { searchProducts, streamProducts, } from "@/services/api";
import { Product } from "@/types/product";
import {
  Search,
  Monitor,
  Smartphone,
  Heart,
  Home,
  ChevronRight,
  Sparkles,
  TrendingUp,
  ShoppingBag,
  Zap,
  Star,
  ArrowRight,
  Grid3X3
} from "lucide-react";

type CategoryItem = {
  label: string;
  query: string;
};

const SESSION_CACHE_VERSION = "2";

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

const categoryIcons: Record<string, React.ReactNode> = {
  Informatique: <Monitor className="w-5 h-5" />,
  Téléphonie: <Smartphone className="w-5 h-5" />,
  "Santé & Beauté": <Heart className="w-5 h-5" />,
  Électroménager: <Home className="w-5 h-5" />,
};

export default function Homes() {
  const [search, setSearch] = useState("");
  const [products, setProducts] = useState<Product[]>([]);
  const [currentPage, setCurrentPage] = useState(1);
  const [activeCategory, setActiveCategory] = useState<string | null>(null);
  const eventSourceRef = useRef<EventSource | null>(null);

  const [loading, setLoading] = useState(false);
  const [message, setMessage] = useState(
    "Recherchez un produit ou choisissez une catégorie."
  );

  const PRODUCTS_PER_PAGE = 24;

  useEffect(() => {
    const savedVersion = sessionStorage.getItem("productsCacheVersion");
    if (savedVersion !== SESSION_CACHE_VERSION) {
      sessionStorage.removeItem("search");
      sessionStorage.removeItem("products");
      sessionStorage.removeItem("page");
      sessionStorage.setItem("productsCacheVersion", SESSION_CACHE_VERSION);
      return;
    }

    const savedSearch = sessionStorage.getItem("search");
    const savedProducts = sessionStorage.getItem("products");
    const savedPage = sessionStorage.getItem("page");

    if (savedSearch) {
      setSearch(savedSearch);
    }

    if (savedProducts) {
      setProducts(JSON.parse(savedProducts));
    }

    if (savedPage) {
      setCurrentPage(Number(savedPage));
    }
  }, []);
  useEffect(() => {
  return () => {
    if (eventSourceRef.current) {
      eventSourceRef.current.close();
      eventSourceRef.current = null;
    }
  };
}, []);

  const loadProducts = async (query: string, displayValue = query) => {
    const value = query.trim();

    if (!value) {
      setMessage("Saisissez un mot clé pour lancer la recherche.");
      return;
    }

    setLoading(true);
    setSearch(displayValue);
    setMessage("");
    setActiveCategory(displayValue);

    // Reset previous results
    setProducts([]);
    setCurrentPage(1);

    let firstBatch = true;
    if (eventSourceRef.current) {
      eventSourceRef.current.close();
    }

    eventSourceRef.current = streamProducts(
      value,

      // Called every time a scraper sends a page
      (newProducts) => {
        if (firstBatch) {
          setLoading(false);
          firstBatch = false;
        }

        setProducts((old) => {
          const updated = [...old, ...newProducts];

          sessionStorage.setItem(
            "productsCacheVersion",
            SESSION_CACHE_VERSION
          );
          sessionStorage.setItem("search", displayValue);
          sessionStorage.setItem("products", JSON.stringify(updated));
          sessionStorage.setItem("page", "1");

          return updated;
        });
      },

      // Finished
      () => {
        eventSourceRef.current = null;
        setLoading(false);
      },

      // Error
      () => {
        if (eventSourceRef.current) {
          eventSourceRef.current.close();
          eventSourceRef.current = null;
        }

        setLoading(false);
        setProducts([]);
        setMessage("Impossible de charger les produits.");
      }
    );
  };

  const indexOfLastProduct = currentPage * PRODUCTS_PER_PAGE;
  const indexOfFirstProduct = indexOfLastProduct - PRODUCTS_PER_PAGE;

  const currentProducts = products.slice(indexOfFirstProduct, indexOfLastProduct);

  const totalPages = Math.ceil(products.length / PRODUCTS_PER_PAGE);

  const changePage = (page: number) => {
    setCurrentPage(page);
    sessionStorage.setItem("page", page.toString());
    window.scrollTo({
      top: 0,
      behavior: "smooth",
    });
  };

  const getVisiblePages = () => {
    const pages: (number | string)[] = [];

    if (totalPages <= 7) {
      return Array.from({ length: totalPages }, (_, i) => i + 1);
    }

    const start = Math.max(2, currentPage - 1);
    const end = Math.min(totalPages - 1, currentPage + 1);

    pages.push(1);

    if (start > 2) {
      pages.push("...");
    }

    for (let i = start; i <= end; i++) {
      pages.push(i);
    }

    if (end < totalPages - 1) {
      pages.push("...");
    }

    pages.push(totalPages);

    return pages;
  };

  const Pagination = () => {
    if (totalPages <= 1) return null;

    return (
      <div className="flex flex-wrap items-center justify-center gap-2 py-6">
        <button
          disabled={currentPage === 1}
          onClick={() => changePage(currentPage - 1)}
          className="group flex items-center gap-2 rounded-xl border border-zinc-200 bg-white px-4 py-2.5 text-sm font-medium text-zinc-700 transition-all hover:border-zinc-300 hover:bg-zinc-50 hover:shadow-md disabled:opacity-40 disabled:hover:bg-white disabled:hover:shadow-none"
        >
          <ArrowRight className="w-4 h-4 rotate-180 transition-transform group-hover:-translate-x-0.5" />
          Previous
        </button>

        {getVisiblePages().map((page, index) =>
          page === "..." ? (
            <span key={index} className="px-2 text-zinc-400">
              ...
            </span>
          ) : (
            <button
              key={`${page}-${index}`}
              onClick={() => changePage(page as number)}
              className={`min-w-[2.5rem] h-10 rounded-xl text-sm font-semibold transition-all duration-200 ${currentPage === page
                ? "bg-gradient-to-br from-violet-600 to-indigo-600 text-white shadow-lg shadow-indigo-500/25 scale-105"
                : "border border-zinc-200 bg-white text-zinc-700 hover:border-zinc-300 hover:bg-zinc-50 hover:shadow-md"
                }`}
            >
              {page}
            </button>
          )
        )}

        <button
          disabled={currentPage === totalPages}
          onClick={() => changePage(currentPage + 1)}
          className="group flex items-center gap-2 rounded-xl border border-zinc-200 bg-white px-4 py-2.5 text-sm font-medium text-zinc-700 transition-all hover:border-zinc-300 hover:bg-zinc-50 hover:shadow-md disabled:opacity-40 disabled:hover:bg-white disabled:hover:shadow-none"
        >
          Next
          <ArrowRight className="w-4 h-4 transition-transform group-hover:translate-x-0.5" />
        </button>
      </div>
    );
  };

  return (
    <main className="min-h-screen bg-gradient-to-br from-slate-50 via-white to-violet-50/30">
      {/* Animated background pattern */}
      <div className="fixed inset-0 -z-10">
        <div className="absolute inset-0 bg-[linear-gradient(to_right,#80808012_1px,transparent_1px),linear-gradient(to_bottom,#80808012_1px,transparent_1px)] bg-[size:24px_24px]"></div>
        <div className="absolute top-0 left-0 right-0 h-96 bg-gradient-to-br from-violet-500/5 via-transparent to-rose-500/5 blur-3xl"></div>
      </div>

      {/* Header */}
      <header className="relative border-b border-zinc-200/60 bg-white/80 backdrop-blur-xl">
        <div className="mx-auto max-w-7xl px-6 py-8">
          <div className="flex flex-col gap-4 md:flex-row md:items-center md:justify-between">
            <div className="space-y-2">
              <div className="flex items-center gap-3">
                <div className="flex items-center gap-2 rounded-2xl bg-gradient-to-br from-violet-600 to-indigo-600 px-4 py-2 shadow-lg shadow-indigo-500/25">
                  <ShoppingBag className="w-6 h-6 text-white" />
                  <span className="text-lg font-bold text-white">TechScout</span>
                </div>
                <span className="hidden sm:inline text-sm font-medium text-zinc-500">
                  Explorer
                </span>
              </div>
              <h1 className="text-4xl font-bold tracking-tight md:text-5xl lg:text-6xl">
                <span className="bg-gradient-to-r from-violet-600 via-indigo-600 to-blue-600 bg-clip-text text-transparent">
                  Découvrez les meilleures offres
                </span>
              </h1>
              <p className="text-zinc-600 max-w-2xl">
                Explorez des milliers de produits aux prix imbattables. Votre prochain achat commence ici.
              </p>
            </div>

            <div className="hidden lg:flex items-center gap-3">
              <div className="flex items-center gap-2 rounded-xl border border-green-200 bg-green-50 px-4 py-2.5">
                <TrendingUp className="w-4 h-4 text-green-600" />
                <span className="text-sm font-semibold text-green-700">Prix compétitifs</span>
              </div>
              <div className="flex items-center gap-2 rounded-xl border border-amber-200 bg-amber-50 px-4 py-2.5">
                <Star className="w-4 h-4 text-amber-600" />
                <span className="text-sm font-semibold text-amber-700">Top marques</span>
              </div>
            </div>
          </div>
        </div>
      </header>

      <div className="mx-auto max-w-7xl px-6 py-8">
        {/* Search Bar */}
        <div className="relative mb-8">
          <div className="absolute -inset-1 rounded-2xl bg-gradient-to-r from-violet-400 via-indigo-400 to-blue-400 opacity-50 blur"></div>
          <div className="relative rounded-2xl border border-zinc-200 bg-white p-6 shadow-xl">
            <div className="flex flex-col gap-3 sm:flex-row">
              <div className="relative flex-1">
                <Search className="absolute left-4 top-1/2 -translate-y-1/2 w-5 h-5 text-zinc-400" />
                <input
                  type="text"
                  placeholder="Recherchez un produit... (ex: smartphone, aspirateur, ordinateur)"
                  value={search}
                  onChange={(event) => setSearch(event.target.value)}
                  onKeyDown={(event) => {
                    if (event.key === "Enter") {
                      loadProducts(search);
                    }
                  }}
                  className="w-full h-14 rounded-xl border-2 border-zinc-200 bg-zinc-50 pl-12 pr-4 text-sm outline-none transition-all placeholder:text-zinc-400 hover:border-zinc-300 focus:border-violet-500 focus:bg-white focus:ring-4 focus:ring-violet-500/10"
                />
              </div>
              <button
                onClick={() => loadProducts(search)}
                className="group flex items-center gap-2 h-14 rounded-xl bg-gradient-to-r from-violet-600 to-indigo-600 px-8 text-sm font-semibold text-white shadow-lg shadow-indigo-500/25 transition-all hover:shadow-xl hover:shadow-indigo-500/30 hover:scale-105 active:scale-100"
              >
                <Search className="w-4 h-4" />
                Rechercher
                <ArrowRight className="w-4 h-4 transition-transform group-hover:translate-x-0.5" />
              </button>
            </div>

            {/* Quick stats when products are loaded */}
            {products.length > 0 && (
              <div className="mt-4 flex flex-wrap gap-3 text-xs font-medium text-zinc-600">
                <div className="flex items-center gap-1.5">
                  <div className="w-1.5 h-1.5 rounded-full bg-green-500"></div>
                  {products.length} résultats trouvés
                </div>
                <div className="flex items-center gap-1.5">
                  <Grid3X3 className="w-3.5 h-3.5" />
                  {totalPages} page{totalPages > 1 ? 's' : ''}
                </div>
              </div>
            )}
          </div>
        </div>

        <div className="grid gap-8 lg:grid-cols-[300px_minmax(0,1fr)]">
          {/* Sidebar Categories */}
          <aside className="lg:sticky lg:top-8 lg:self-start">
            <div className="rounded-2xl border border-zinc-200 bg-white/80 backdrop-blur p-6 shadow-xl">
              <div className="flex items-center gap-2 mb-6">
                <div className="p-2 rounded-lg bg-violet-100">
                  <Grid3X3 className="w-5 h-5 text-violet-600" />
                </div>
                <h2 className="text-lg font-bold text-zinc-900">Catégories</h2>
              </div>

              <div className="space-y-6">
                {Object.entries(categories).map(([mainCategory, subCategories]) => (
                  <div key={mainCategory}>
                    <div className="flex items-center gap-2 mb-3">
                      <div className="p-1.5 rounded-md bg-gradient-to-br from-zinc-100 to-zinc-50 border border-zinc-200">
                        {categoryIcons[mainCategory]}
                      </div>
                      <h3 className="text-sm font-bold uppercase tracking-wider text-zinc-500">
                        {mainCategory}
                      </h3>
                    </div>

                    <div className="flex flex-col gap-0.5">
                      {subCategories.map((subCategory) => (
                        <button
                          key={subCategory.label}
                          onClick={() =>
                            loadProducts(subCategory.query, subCategory.label)
                          }
                          className={`group flex items-center justify-between rounded-xl px-3 py-2.5 text-sm font-medium transition-all ${activeCategory === subCategory.label
                            ? "bg-gradient-to-r from-violet-50 to-indigo-50 text-violet-700 border border-violet-200"
                            : "text-zinc-600 hover:bg-zinc-50 hover:text-zinc-900"
                            }`}
                        >
                          <span>{subCategory.label}</span>
                          <ChevronRight className={`w-4 h-4 transition-all ${activeCategory === subCategory.label
                            ? "opacity-100 translate-x-0"
                            : "opacity-0 -translate-x-2 group-hover:opacity-100 group-hover:translate-x-0"
                            }`} />
                        </button>
                      ))}
                    </div>
                  </div>
                ))}
              </div>
            </div>
          </aside>

          {/* Products Section */}
          <section>
            {/* Results Header */}
            <div className="mb-6 flex items-center justify-between">
              <div className="flex items-center gap-3">
                <h2 className="text-2xl font-bold text-zinc-900">
                  {activeCategory || "Tous les produits"}
                </h2>
                {activeCategory && (
                  <span className="rounded-full bg-gradient-to-r from-violet-100 to-indigo-100 px-3 py-1 text-xs font-semibold text-violet-700 border border-violet-200">
                    {products.length} résultat{products.length !== 1 ? 's' : ''}
                  </span>
                )}
              </div>
              {products.length > 0 && (
                <div className="hidden sm:flex items-center gap-2 text-sm text-zinc-500">
                  <span>Page</span>
                  <span className="font-semibold text-zinc-900">{currentPage}</span>
                  <span>sur</span>
                  <span className="font-semibold text-zinc-900">{totalPages}</span>
                </div>
              )}
            </div>

            {/* Loading State (only before the first batch of products has arrived) */}
            {loading && products.length === 0 && (
              <div className="space-y-4">
                <div className="rounded-2xl border border-zinc-200 bg-white p-12 text-center shadow-sm">
                  <div className="flex flex-col items-center gap-4">
                    <div className="relative">
                      <div className="w-16 h-16 rounded-full border-4 border-zinc-100 border-t-violet-600 animate-spin"></div>
                      <Search className="absolute inset-0 m-auto w-6 h-6 text-violet-600" />
                    </div>
                    <div>
                      <p className="text-lg font-semibold text-zinc-700">Recherche en cours</p>
                      <p className="text-sm text-zinc-500">Nous parcourons les meilleures offres pour vous...</p>
                    </div>
                  </div>
                </div>
                {/* Skeleton cards */}
                <div className="grid gap-5 sm:grid-cols-2 xl:grid-cols-3">
                  {[...Array(6)].map((_, i) => (
                    <div key={i} className="rounded-2xl border border-zinc-200 bg-white p-6 shadow-sm animate-pulse">
                      <div className="w-full h-48 bg-zinc-100 rounded-xl mb-4"></div>
                      <div className="h-4 bg-zinc-100 rounded w-3/4 mb-2"></div>
                      <div className="h-4 bg-zinc-100 rounded w-1/2"></div>
                    </div>
                  ))}
                </div>
              </div>
            )}

            {/* Streaming indicator: more products are still arriving in the background */}
            {loading && products.length > 0 && (
              <div className="mb-4 flex items-center gap-2 rounded-xl border border-violet-100 bg-violet-50 px-4 py-2.5 text-sm font-medium text-violet-700">
                <div className="w-4 h-4 rounded-full border-2 border-violet-200 border-t-violet-600 animate-spin"></div>
                Chargement de produits supplémentaires...
              </div>
            )}

            {/* Empty State */}
            {!loading && products.length === 0 && (
              <div className="rounded-2xl border-2 border-dashed border-zinc-300 bg-white p-16 text-center shadow-sm">
                <div className="flex flex-col items-center gap-4">
                  <div className="rounded-2xl bg-gradient-to-br from-violet-50 to-indigo-50 p-6 border border-violet-100">
                    <Sparkles className="w-12 h-12 text-violet-500" />
                  </div>
                  <div className="space-y-2">
                    <p className="text-xl font-semibold text-zinc-700">Commencez votre recherche</p>
                    <p className="text-sm text-zinc-500 max-w-md mx-auto">
                      {message || "Utilisez la barre de recherche ou sélectionnez une catégorie pour découvrir nos produits"}
                    </p>
                  </div>
                  <div className="flex gap-2 mt-2">
                    <div className="flex items-center gap-2 rounded-xl border border-zinc-200 bg-zinc-50 px-4 py-2 text-sm text-zinc-600">
                      <TrendingUp className="w-4 h-4" />
                      Meilleurs prix
                    </div>
                    <div className="flex items-center gap-2 rounded-xl border border-zinc-200 bg-zinc-50 px-4 py-2 text-sm text-zinc-600">
                      <Star className="w-4 h-4" />
                      Qualité garantie
                    </div>
                  </div>
                </div>
              </div>
            )}

            {/* Products Grid */}
            {products.length > 0 && (
              <>
                <Pagination />
                <div className="grid gap-5 sm:grid-cols-2 xl:grid-cols-3">
                  {currentProducts.map((product, index) => (
                    <div
                      key={product.url || product.name}
                      className="animate-fadeInUp"
                      style={{
                        animationDelay: `${index * 50}ms`,
                        animationFillMode: 'both'
                      }}
                    >
                      <ProductCard product={product} />
                    </div>
                  ))}
                </div>
                <Pagination />
              </>
            )}
          </section>
        </div>
      </div>

      <style jsx>{`
        @keyframes fadeInUp {
          from {
            opacity: 0;
            transform: translateY(20px);
          }
          to {
            opacity: 1;
            transform: translateY(0);
          }
        }
        .animate-fadeInUp {
          animation: fadeInUp 0.5s ease-out;
        }
      `}</style>
    </main>
  );
}