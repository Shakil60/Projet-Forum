# Diagramme C4 — CinéTalk

Présentation de l'architecture du projet selon le modèle **C4** (Simon Brown).
Le niveau demandé est le **niveau 2 (conteneurs)** ; le niveau 1 (contexte) est ajouté pour situer le système.

## Niveau 1 — Contexte

Vue d'ensemble : qui utilise CinéTalk et à quels systèmes externes il se connecte.

```mermaid
C4Context
    title Diagramme de contexte - CinéTalk

    Person(visiteur, "Visiteur", "Lit les fils publics")
    Person(membre, "Membre", "Publie des fils, des messages et réagit")
    Person(admin, "Administrateur", "Modère et bannit les comptes")

    System(cinetalk, "CinéTalk", "Plateforme de discussion sur le cinéma et les séries")

    System_Ext(tmdb, "TMDB", "API The Movie Database (films, séries, personnes)")

    Rel(visiteur, cinetalk, "Consulte", "HTTPS")
    Rel(membre, cinetalk, "Participe", "HTTPS")
    Rel(admin, cinetalk, "Administre", "HTTPS")
    Rel(cinetalk, tmdb, "Récupère le catalogue", "HTTPS/JSON")
```

## Niveau 2 — Conteneurs

Les grandes briques techniques exécutables et leurs échanges.
Un « conteneur » au sens C4 = une unité exécutable/déployable (≠ conteneur Docker).

```mermaid
C4Container
    title Diagramme de conteneurs - CinéTalk

    Person(membre, "Membre / Visiteur", "Utilise le forum depuis son navigateur")
    Person(admin, "Administrateur", "Modère le contenu et les comptes")

    System_Boundary(cinetalk, "CinéTalk") {
        Container(navigateur, "Navigateur web", "HTML, CSS, JavaScript", "Affiche les pages et envoie les réactions like/dislike en AJAX")
        Container(app, "Application Web", "Go, gorilla/mux, html/template", "Authentification JWT, fils, messages, réactions, catalogue, administration ; rendu des pages côté serveur")
        ContainerDb(db, "Base de données", "MySQL 8", "Utilisateurs, fils, messages, réactions, genres")
    }

    System_Ext(tmdb, "API TMDB", "Service externe REST/JSON")

    Rel(membre, navigateur, "Navigue", "")
    Rel(admin, navigateur, "Navigue", "")
    Rel(navigateur, app, "Requêtes de pages et API réactions", "HTTPS / JSON")
    Rel(app, navigateur, "Pages HTML", "HTTPS")
    Rel(app, db, "Lit et écrit", "SQL / port 3306")
    Rel(app, tmdb, "Interroge films, séries, personnes", "HTTPS / JSON")
```

## Niveau 3 — Composants (bonus, intérieur de l'Application Web)

Découpage interne en couches, qui correspond aux dossiers du dépôt.

```mermaid
C4Component
    title Diagramme de composants - Application Web CinéTalk

    Container(navigateur, "Navigateur web", "HTML/CSS/JS", "Client")
    ContainerDb(db, "Base de données", "MySQL", "Stockage")
    System_Ext(tmdb, "API TMDB", "Service externe")

    Container_Boundary(app, "Application Web") {
        Component(routers, "Routers", "gorilla/mux", "Déclare les routes")
        Component(middleware, "Middleware", "Go", "Auth, rôle, bannissement")
        Component(controllers, "Controllers", "Go", "Valide les requêtes, prépare les réponses")
        Component(services, "Services", "Go", "Logique métier (fils, messages, réactions, TMDB)")
        Component(repositories, "Repositories", "Go", "Accès aux données, requêtes SQL")
        Component(auth, "Auth", "Go, JWT, SHA-512", "Hachage des mots de passe et jetons")
        Component(helper, "Helper", "html/template", "Rendu des vues et réponses JSON")
    }

    Rel(navigateur, routers, "Appelle", "HTTPS")
    Rel(routers, middleware, "Passe par")
    Rel(middleware, controllers, "Transmet")
    Rel(controllers, services, "Appelle")
    Rel(controllers, helper, "Rend la vue")
    Rel(services, repositories, "Utilise")
    Rel(services, auth, "Vérifie / signe")
    Rel(services, tmdb, "Appelle", "HTTPS/JSON")
    Rel(repositories, db, "Lit/écrit", "SQL")
```

## Légende

| Élément | Signification |
|---|---|
| `Person` | Acteur humain |
| `Container` | Brique exécutable (app, navigateur) |
| `ContainerDb` | Brique de stockage |
| `Component` | Regroupement de code dans un conteneur |
| `System_Ext` | Système externe (non développé par l'équipe) |
| `Rel` | Relation orientée avec son protocole |

> Les diagrammes Mermaid se rendent directement sur GitHub.
> Pour exporter une image : coller le bloc dans https://mermaid.live (export PNG/SVG),
> ou exécuter `npx -y @mermaid-js/mermaid-cli -i fichier.mmd -o fichier.png`.
