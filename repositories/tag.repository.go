package repositories

// Acces a la base de donnees pour les tags.

import (
	"database/sql"
	"forum/models"
	"fmt"
	"strings"
)

type TagRepository struct {
	db *sql.DB
}

func InitTagRepository(db *sql.DB) *TagRepository {
	return &TagRepository{db}
}

// Recupere tous les tags tries par nom.
func (r *TagRepository) ReadAll() ([]models.Tag, error) {
	var tags []models.Tag
	rows, err := r.db.Query("SELECT `id`, `nom` FROM `tags` ORDER BY `nom` ASC;")
	if err != nil {
		return nil, fmt.Errorf("Erreur recuperation tags - %s", err.Error())
	}
	defer rows.Close()

	for rows.Next() {
		var tag models.Tag
		if scanErr := rows.Scan(&tag.Id, &tag.Nom); scanErr != nil {
			return nil, scanErr
		}
		tags = append(tags, tag)
	}

	return tags, nil
}

func (r *TagRepository) FindById(id int) (models.Tag, error) {
	var tag models.Tag
	err := r.db.QueryRow("SELECT `id`, `nom` FROM `tags` WHERE `id` = ?;", id).Scan(&tag.Id, &tag.Nom)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.Tag{}, nil
		}
		return models.Tag{}, fmt.Errorf("Erreur recuperation tag - %s", err.Error())
	}
	return tag, nil
}

// Recherche un tag par son nom.
func (r *TagRepository) FindByName(nom string) (models.Tag, error) {
	var tag models.Tag
	err := r.db.QueryRow("SELECT `id`, `nom` FROM `tags` WHERE `nom` = ?;", nom).Scan(&tag.Id, &tag.Nom)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.Tag{}, nil
		}
		return models.Tag{}, fmt.Errorf("Erreur recuperation tag - %s", err.Error())
	}
	return tag, nil
}

// Renvoie l'identifiant d'un tag existant ou le cree s'il n'existe pas.
func (r *TagRepository) FindOrCreate(nom string) (int, error) {
	nom = strings.TrimSpace(nom)
	if nom == "" {
		return 0, fmt.Errorf("Erreur tag - nom vide")
	}

	existing, err := r.FindByName(nom)
	if err != nil {
		return 0, err
	}
	if existing.Id != 0 {
		return existing.Id, nil
	}

	result, insertErr := r.db.Exec("INSERT INTO `tags`(`nom`) VALUES (?);", nom)
	if insertErr != nil {
		return 0, fmt.Errorf("Erreur creation tag - %s", insertErr.Error())
	}

	id, idErr := result.LastInsertId()
	if idErr != nil {
		return 0, idErr
	}
	return int(id), nil
}
