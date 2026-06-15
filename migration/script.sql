DROP DATABASE IF EXISTS cinetalk;
CREATE DATABASE cinetalk
    DEFAULT CHARACTER SET utf8mb4
    DEFAULT COLLATE utf8mb4_unicode_ci;

USE cinetalk;

CREATE TABLE utilisateurs (
    id INT AUTO_INCREMENT PRIMARY KEY,
    nom_utilisateur VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(150) UNIQUE NOT NULL,
    mot_de_passe VARCHAR(255) NOT NULL,
    sel VARCHAR(64) NOT NULL,
    role ENUM('utilisateur', 'administrateur') NOT NULL DEFAULT 'utilisateur',
    banni BOOLEAN NOT NULL DEFAULT FALSE,
    date_creation TIMESTAMP DEFAULT CURRENT_TIMESTAMP
) ENGINE=InnoDB;

CREATE TABLE tags (
    id INT AUTO_INCREMENT PRIMARY KEY,
    nom VARCHAR(50) UNIQUE NOT NULL
) ENGINE=InnoDB;

CREATE TABLE fils (
    id INT AUTO_INCREMENT PRIMARY KEY,
    titre VARCHAR(150) NOT NULL,
    contenu TEXT NOT NULL,
    utilisateur_id INT NOT NULL,
    etat ENUM('ouvert', 'ferme', 'archive') NOT NULL DEFAULT 'ouvert',
    date_creation TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_fils_utilisateurs
        FOREIGN KEY (utilisateur_id)
        REFERENCES utilisateurs(id)
        ON DELETE CASCADE
) ENGINE=InnoDB;

CREATE TABLE fil_tags (
    fil_id INT NOT NULL,
    tag_id INT NOT NULL,
    PRIMARY KEY (fil_id, tag_id),

    CONSTRAINT fk_filtags_fils
        FOREIGN KEY (fil_id)
        REFERENCES fils(id)
        ON DELETE CASCADE,

    CONSTRAINT fk_filtags_tags
        FOREIGN KEY (tag_id)
        REFERENCES tags(id)
        ON DELETE CASCADE
) ENGINE=InnoDB;

CREATE TABLE messages (
    id INT AUTO_INCREMENT PRIMARY KEY,
    fil_id INT NOT NULL,
    utilisateur_id INT NOT NULL,
    contenu TEXT NOT NULL,
    date_envoi TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_messages_fils
        FOREIGN KEY (fil_id)
        REFERENCES fils(id)
        ON DELETE CASCADE,

    CONSTRAINT fk_messages_utilisateurs
        FOREIGN KEY (utilisateur_id)
        REFERENCES utilisateurs(id)
        ON DELETE CASCADE
) ENGINE=InnoDB;

CREATE TABLE reactions (
    id INT AUTO_INCREMENT PRIMARY KEY,
    message_id INT NOT NULL,
    utilisateur_id INT NOT NULL,
    type ENUM('like', 'dislike') NOT NULL,

    UNIQUE KEY uq_reaction_utilisateur (message_id, utilisateur_id),

    CONSTRAINT fk_reactions_messages
        FOREIGN KEY (message_id)
        REFERENCES messages(id)
        ON DELETE CASCADE,

    CONSTRAINT fk_reactions_utilisateurs
        FOREIGN KEY (utilisateur_id)
        REFERENCES utilisateurs(id)
        ON DELETE CASCADE
) ENGINE=InnoDB;
