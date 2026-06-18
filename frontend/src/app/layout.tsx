import type { Metadata } from "next";
import "./globals.css";

export const metadata: Metadata = {
  title: "Tunisianet Explorer",
  description: "Rechercher et explorer les produits de Tunisianet.",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="fr" className="h-full antialiased">
      <body className="flex min-h-full flex-col">{children}</body>
    </html>
  );
}
