# Tunisianet Explorer

Application web qui permet de rechercher et consulter des produits collectés depuis Tunisianet.

## Stack

- Backend: Go, Echo, GoQuery
- Frontend: Next.js, React, TypeScript, Tailwind CSS

## Fonctionnalités

- Page d'accueil avec champ de recherche.
- Liste des catégories demandées dans le sujet.
- Affichage des produits trouvés avec image, nom et prix.
- Page détail pour consulter les informations d'un produit.
- API backend séparée du frontend.

## Lancer le backend

```bash
cd backend
go run .
```

Le backend démarre sur `http://localhost:8080`.

## Lancer le frontend

```bash
cd frontend
npm install
npm run dev
```

Le frontend démarre sur `http://localhost:3000`.

## API

- `GET /products?search=ordinateur`
- `GET /product/details?url=<product-url>`

## Vérification

```bash
cd backend
go test ./...
```

```bash
cd frontend
npm run lint
```
