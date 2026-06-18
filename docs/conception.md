# Conception — CinéTalk (ProjJS.2 Forum)

Livrables de conception demandés par le sujet : **diagrammes de cas d'utilisation par rôle**
et **modèle logique de données (MLD)**.

## 1. Rôles et héritage

Trois rôles, du moins au plus de droits. Chaque rôle hérite des actions du précédent.

- **Visiteur** : non connecté.
- **Membre** : visiteur connecté (hérite des actions du visiteur).
- **Administrateur** : membre avec droits de modération (hérite des actions du membre).

## 2. Cas d'utilisation par rôle

### Visiteur

```mermaid
flowchart LR
    visiteur([Visiteur])
    subgraph S1[Système CinéTalk]
        v1(Consulter les fils et messages)
        v2(Consulter le catalogue cinéma)
        v3(S'inscrire)
        v4(Se connecter)
    end
    visiteur --> v1
    visiteur --> v2
    visiteur --> v3
    visiteur --> v4
```

### Membre (hérite du Visiteur)

```mermaid
flowchart LR
    membre([Membre])
    subgraph S2[Système CinéTalk]
        m0(Actions du visiteur)
        m1(Créer / modifier / supprimer ses fils)
        m2(Poster / modifier / supprimer ses messages)
        m3(Réagir like / dislike)
        m4(Rechercher des fils)
        m5(Filtrer par genre)
        m6(Se déconnecter)
        auth{{S'authentifier}}
    end
    membre --> m0
    membre --> m1
    membre --> m2
    membre --> m3
    membre --> m4
    membre --> m5
    membre --> m6
    m1 -.->|include| auth
    m2 -.->|include| auth
    m3 -.->|include| auth
```

### Administrateur (hérite du Membre)

```mermaid
flowchart LR
    admin([Administrateur])
    subgraph S3[Système CinéTalk]
        a0(Actions du membre)
        a1(Consulter le tableau de bord)
        a2(Modérer un fil : changer l'état)
        a3(Supprimer n'importe quel fil ou message)
        a4(Bannir / réactiver un compte)
    end
    admin --> a0
    admin --> a1
    admin --> a2
    admin --> a3
    admin --> a4
```

## 3. Modèle logique de données (MLD)

Schéma relationnel issu de `migration/script.sql`.

```mermaid
erDiagram
    UTILISATEURS ||--o{ FILS : "rédige"
    UTILISATEURS ||--o{ MESSAGES : "poste"
    UTILISATEURS ||--o{ REACTIONS : "exprime"
    FILS ||--o{ MESSAGES : "contient"
    MESSAGES ||--o{ REACTIONS : "reçoit"
    FILS ||--o{ FIL_TAGS : ""
    TAGS ||--o{ FIL_TAGS : ""

    UTILISATEURS {
        int id PK
        varchar nom_utilisateur "unique"
        varchar email "unique"
        varchar mot_de_passe
        varchar sel
        enum role "utilisateur / administrateur"
        boolean banni
        timestamp date_creation
    }
    TAGS {
        int id PK
        varchar nom "unique"
    }
    FILS {
        int id PK
        varchar titre
        text contenu
        int utilisateur_id FK
        enum etat "ouvert / ferme / archive"
        timestamp date_creation
    }
    FIL_TAGS {
        int fil_id PK,FK
        int tag_id PK,FK
    }
    MESSAGES {
        int id PK
        int fil_id FK
        int utilisateur_id FK
        text contenu
        timestamp date_envoi
    }
    REACTIONS {
        int id PK
        int message_id FK
        int utilisateur_id FK
        enum type "like / dislike"
    }
```

### Forme textuelle du MLD

```
utilisateurs(#id, nom_utilisateur, email, mot_de_passe, sel, role, banni, date_creation)
tags(#id, nom)
fils(#id, titre, contenu, #utilisateur_id, etat, date_creation)
fil_tags(#fil_id, #tag_id)
messages(#id, #fil_id, #utilisateur_id, contenu, date_envoi)
reactions(#id, #message_id, #utilisateur_id, type)   -- unique(message_id, utilisateur_id)
```

> `#` = clé primaire, soulignement des clés étrangères dans la forme classique.
> La relation **fils ⟷ tags** (plusieurs-à-plusieurs) est résolue par la table d'association `fil_tags`.
> Toutes les clés étrangères sont en `ON DELETE CASCADE`.
