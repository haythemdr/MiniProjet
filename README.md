# Tunisianet Explorer

Application web de recherche de produits sur Tunisianet.

Le projet contient un backend Go avec Echo pour exposer une petite API, et un frontend Next.js pour la recherche, les catégories et la page détail produit.

## Fonctionnalités

- Recherche de produits depuis la page d'accueil.
- Liste des catégories demandées dans le sujet.
- Affichage des produits trouvés avec image, nom et prix.
- Page détail avec image, prix, disponibilité et description.
- Backend et frontend séparés.

## Technologies

- Go
- Echo
- GoQuery
- Next.js
- TypeScript
- Tailwind CSS

## Lancer le projet

Backend :

```bash
cd backend
go run .
```

Le backend démarre sur `http://localhost:8080`.

Frontend :

```bash
cd frontend
npm install
npm run dev
```

Le frontend démarre sur `http://localhost:3000`.

## API

- `GET /products?search=ordinateur`
- `GET /product/details?url=https://www.tunisianet.com.tn/...`

## Vérification

```bash
cd backend
go test ./...
```

```bash
cd frontend
npm run lint
npm run build
```
