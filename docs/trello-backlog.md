# Backlog Trello — CinéTalk (ProjJS.2 Forum)

Aligné sur le sujet Ynov. Le tableau doit contenir **au minimum** : une colonne « À faire »,
une colonne « En cours », une colonne « Terminé », des tâches de **conception**, des tâches de
**développement des fonctionnalités obligatoires**, et une **répartition** entre les membres.

Équipe : **Shakil Khaldi (SK)** · **Hugo Berton (HB)** · **Nohan Reis (NR)**

## Colonnes (listes)
**À faire · En cours · Terminé**
(optionnel : « En revue » pour la relecture avant fusion)

## Étiquettes
`Conception` · `Auth` · `Fils` · `Messages` · `Réactions` · `Recherche` · `Admin` · `Catalogue` · `BDD` · `UI` · `DevOps` · `Doc`

---

## Conception (livrables de cadrage)
- [Conception] Identité de la plateforme (thème cinéma, ambiance graphique) — HB
- [Conception] Rôles et actions (visiteur, membre, administrateur) — SK
- [Conception] Diagramme de cas d'utilisation par rôle — SK
- [Conception] Modèle logique de données (MLD) — NR
- [Conception] Mise en place du dépôt Git + README — HB
- [Conception] Tableau Trello et répartition des tâches — HB
- [Conception] Diagramme C4 (contexte / conteneurs) pour la présentation — HB

## Base de données
- [BDD] Schéma relationnel (utilisateurs, fils, messages, réactions, genres) — NR
- [BDD] Script de création des tables — NR
- [BDD] Jeu de données de test (seed) — NR

## Authentification (obligatoire)
- [Auth] Inscription (nom d'utilisateur ou e-mail, règles de mot de passe) — SK
- [Auth] Hachage SHA-512 + sel — SK
- [Auth] Connexion et jeton JWT en cookie HttpOnly — SK
- [Auth] Middleware d'authentification, de rôle et de bannissement — SK
- [Auth] Déconnexion — SK

## Fils de discussion — CRUD (obligatoire)
- [Fils] Création, modification, suppression d'un fil — HB
- [Fils] États du fil : ouvert / fermé / archivé — HB
- [Fils] Association d'un fil à un ou plusieurs genres — HB
- [Fils] Pagination configurable (10 / 20 / 30 / tous) — HB

## Messages — CRUD (obligatoire)
- [Messages] Publication, modification, suppression d'un message — NR
- [Messages] Tri : récents / chronologique / popularité — NR

## Réactions
- [Réactions] Like / dislike avec calcul du score — NR
- [Réactions] Mise à jour asynchrone (AJAX) sans rechargement — NR

## Moteur de recherche interne (obligatoire)
- [Recherche] Recherche de fils par titre ou par genre — HB
- [Recherche] Filtrage par genre — HB

## Tableau de bord administrateur (obligatoire)
- [Admin] Statistiques (nombre d'utilisateurs, de fils) — SK
- [Admin] Modération des fils et changement d'état — SK
- [Admin] Bannissement / réactivation des comptes — SK

## Catalogue cinéma (bonus)
- [Catalogue] Service d'appel à l'API TMDB — HB
- [Catalogue] Pages films, séries, acteurs et réalisateurs — HB
- [Catalogue] Recherche dans le catalogue — HB

## Interface (UI)
- [UI] Charte graphique et CSS (variables, responsive) — HB
- [UI] Accessibilité (focus clavier, aria, skip-link) — HB
- [UI] Formulaires et états vides — HB

## DevOps & Documentation
- [DevOps] Dockerfile + docker-compose (app + MySQL + migrations) — HB
- [Doc] README et guide d'installation — HB
- [Doc] Support de présentation (soutenance 10 min + 5 min questions) — SK/HB/NR

---

## État actuel (juin 2026)
Le code est **terminé** : auth, CRUD fils/messages, réactions, recherche, admin, catalogue,
Docker et UI polie. Sur le tableau, la majorité des cartes « développement » vont donc en
**Terminé** ; restent en cours/à faire les livrables de présentation (C4, support, MLD/cas
d'utilisation s'ils ne sont pas encore formalisés).

## Remplissage automatique (API Trello)
Possible dès que tu fournis un **token** (la clé seule ne suffit pas). Voir la procédure dans
la réponse de l'assistant.
