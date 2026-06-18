# CinéTalk 🎬

CinéTalk est une plateforme communautaire de discussion dédiée au cinéma et aux séries,
inspirée de services comme Stack Overflow, Reddit et IMDb. Les membres peuvent ouvrir des
fils de discussion, échanger des messages, réagir aux contributions de la communauté et
organiser les contenus par genres.

## Membres de l'équipe

-Shakil Khaldi
-Hugo Berton
-Nohan Reis

## Fonctionnalités

- Inscription et connexion (par nom d'utilisateur ou e-mail) avec mots de passe hachés en SHA-512 + sel
- Authentification par jeton **JWT** stocké dans un cookie sécurisé (HttpOnly)
- Création, modification et suppression de fils de discussion
- Gestion des états d'un fil : `ouvert`, `fermé`, `archivé`
- Publication, modification et suppression de messages
- Système de réactions like / dislike avec calcul d'un score de popularité (AJAX)
- Tri des messages : plus récents, chronologique ou par popularité
- Pagination configurable (10, 20, 30 ou tous les éléments)
- Association des fils à un ou plusieurs genres / catégories et filtrage par genre
- Recherche de fils par titre ou par genre (réservée aux membres connectés)
- **Catalogue TMDB** : films, séries, acteurs et réalisateurs via l'API The Movie Database
- Tableau de bord d'administration : modération des fils, suppression de contenus,
  changement d'état et bannissement de comptes

## Technologies utilisées

- **Go** (langage principal : routes, contrôleurs, logique métier, accès aux données, sécurité)
- **gorilla/mux** pour le routage HTML
- **MySQL** (base de données relationnelle)
- **golang-jwt** pour la génération et la validation des jetons JWT
- **html/template** pour le rendu des vues côté serveur
- **JavaScript** (uniquement pour les réactions asynchrones et certaines interactions)

## Architecture

Le projet est monolithique et organisé en couches inspirées du modèle MVC, avec une
séparation claire des responsabilités :

```
forum/
├── main.go              Point d'entrée de l'application
├── app/                 Assemblage des dépendances et démarrage
├── config/              Chargement de la configuration et connexion à la base
├── auth/                Hachage des mots de passe, génération et validation des JWT
├── middleware/          Contrôle d'authentification, de rôle et de bannissement
├── routers/             Définition des routes (Router)
├── controllers/         Réception des requêtes et validation des entrées (Controller)
├── services/            Logique métier (Service)
├── repositories/        Accès aux données et requêtes SQL (Repository)
├── models/              Représentation des données (Model / Entity)
├── dto/                 Objets de transfert et pagination
├── helper/              Rendu des vues et réponses JSON
├── views/               Pages HTML (View)
├── static/              Feuilles de style et scripts JavaScript
└── migration/           Scripts de création et de remplissage de la base
```

## Prérequis

- [Go](https://go.dev/dl/) version 1.22 ou supérieure
- [MySQL](https://dev.mysql.com/downloads/) version 8 ou supérieure (ou MariaDB compatible)

## Installation

1. **Récupérer le projet**

   ```bash
   git clone <url-du-depot>
   cd "Projet Forum"
   ```

2. **Créer et remplir la base de données**

   ```bash
   mysql -u root -p < migration/script.sql
   mysql -u root -p < migration/seed.sql
   ```

   Le premier script crée la base `cinetalk` et ses tables, le second insère les données de test.

3. **Configurer les variables d'environnement**

   Copier le fichier `.env.example` en `.env` puis adapter les valeurs à votre installation :

   ```bash
   cp .env.example .env
   ```

   ```
   DB_NAME=cinetalk
   DB_USER=root
   DB_PWD=votre_mot_de_passe
   DB_HOST=localhost
   DB_PORT=3306
   JWT_SECRET=une_cle_secrete_longue
   SERVER_PORT=8080
   TMDB_API_KEY=votre_cle_api_tmdb
   ```

   Pour obtenir une clé TMDB : créez un compte sur [themoviedb.org](https://www.themoviedb.org/),
   puis demandez une clé API (type Developer) dans Paramètres → API.

4. **Installer les dépendances**

   ```bash
   go mod download
   ```

5. **Lancer l'application**

   ```bash
   go run .
   ```

6. **Ouvrir le forum**

   Rendez-vous sur [http://localhost:8080](http://localhost:8080).
   Le catalogue cinéma est accessible sur [http://localhost:8080/catalog](http://localhost:8080/catalog).

## Comptes de test

| Identifiant | Mot de passe      | Rôle           |
|-------------|-------------------|----------------|
| `admin`     | `Admin@CineTalk1` | Administrateur |
| `alice`     | `Cinephile@2026`  | Utilisateur    |
| `bob`       | `Cinephile@2026`  | Utilisateur    |
| `claire`    | `Cinephile@2026`  | Utilisateur    |
| `david`     | `Cinephile@2026`  | Utilisateur    |

## Sécurité

- Les mots de passe ne sont jamais stockés en clair : ils sont hachés en **SHA-512** avec un
  sel unique par utilisateur.
- Les règles de création de mot de passe imposent au minimum 12 caractères, une majuscule et
  un caractère spécial.
- L'identité des membres est vérifiée à chaque requête grâce au jeton JWT, et les comptes
  bannis perdent immédiatement l'accès aux fonctionnalités authentifiées.
