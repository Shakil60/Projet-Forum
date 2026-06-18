package repositories

// Acces a la base pour les favoris et commentaires de films.

import (
	"database/sql"
	"fmt"
	"forum/models"
)

type FilmRepository struct {
	db *sql.DB
}

func InitFilmRepository(db *sql.DB) *FilmRepository {
	return &FilmRepository{db}
}

// Indique si un membre a ce film/serie en favori.
func (r *FilmRepository) IsFavori(userId, tmdbId int, mediaType string) (bool, error) {
	var x int
	err := r.db.QueryRow("SELECT 1 FROM `film_favoris` WHERE `utilisateur_id` = ? AND `tmdb_id` = ? AND `media_type` = ?;", userId, tmdbId, mediaType).Scan(&x)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("Erreur lecture favori - %s", err.Error())
	}
	return true, nil
}

// Compte combien de membres ont mis ce film/serie en favori.
func (r *FilmRepository) CountFavoris(tmdbId int, mediaType string) (int, error) {
	var n int
	err := r.db.QueryRow("SELECT COUNT(*) FROM `film_favoris` WHERE `tmdb_id` = ? AND `media_type` = ?;", tmdbId, mediaType).Scan(&n)
	if err != nil {
		return 0, fmt.Errorf("Erreur comptage favoris - %s", err.Error())
	}
	return n, nil
}

// Ajoute ou retire le favori (bascule) et renvoie le nouvel etat.
func (r *FilmRepository) ToggleFavori(userId, tmdbId int, mediaType, titre, affiche string) (bool, error) {
	favori, err := r.IsFavori(userId, tmdbId, mediaType)
	if err != nil {
		return false, err
	}
	if favori {
		_, err = r.db.Exec("DELETE FROM `film_favoris` WHERE `utilisateur_id` = ? AND `tmdb_id` = ? AND `media_type` = ?;", userId, tmdbId, mediaType)
		return false, err
	}
	_, err = r.db.Exec("INSERT INTO `film_favoris`(`utilisateur_id`, `tmdb_id`, `media_type`, `titre`, `affiche`) VALUES (?,?,?,?,?);", userId, tmdbId, mediaType, titre, affiche)
	return true, err
}

// Liste les favoris d'un membre, du plus recent au plus ancien.
func (r *FilmRepository) ListFavoris(userId int) ([]models.FilmFavori, error) {
	rows, err := r.db.Query("SELECT `tmdb_id`, `media_type`, `titre`, `affiche`, `date_ajout` FROM `film_favoris` WHERE `utilisateur_id` = ? ORDER BY `date_ajout` DESC;", userId)
	if err != nil {
		return nil, fmt.Errorf("Erreur liste favoris - %s", err.Error())
	}
	defer rows.Close()

	favoris := []models.FilmFavori{}
	for rows.Next() {
		var f models.FilmFavori
		var affiche sql.NullString
		if err := rows.Scan(&f.TmdbId, &f.MediaType, &f.Titre, &affiche, &f.CreatedAt); err != nil {
			return nil, err
		}
		f.Affiche = affiche.String
		favoris = append(favoris, f)
	}
	return favoris, nil
}

// Ajoute un commentaire sous un film/serie.
func (r *FilmRepository) AddComment(tmdbId int, mediaType string, userId int, contenu string) error {
	_, err := r.db.Exec("INSERT INTO `film_commentaires`(`tmdb_id`, `media_type`, `utilisateur_id`, `contenu`) VALUES (?,?,?,?);", tmdbId, mediaType, userId, contenu)
	if err != nil {
		return fmt.Errorf("Erreur ajout commentaire - %s", err.Error())
	}
	return nil
}

// Liste les commentaires d'un film/serie, du plus recent au plus ancien.
func (r *FilmRepository) ListComments(tmdbId int, mediaType string) ([]models.FilmComment, error) {
	query := "SELECT c.`id`, c.`utilisateur_id`, u.`nom_utilisateur`, c.`contenu`, c.`date_creation` " +
		"FROM `film_commentaires` c JOIN `utilisateurs` u ON u.`id` = c.`utilisateur_id` " +
		"WHERE c.`tmdb_id` = ? AND c.`media_type` = ? ORDER BY c.`date_creation` DESC;"
	rows, err := r.db.Query(query, tmdbId, mediaType)
	if err != nil {
		return nil, fmt.Errorf("Erreur liste commentaires - %s", err.Error())
	}
	defer rows.Close()

	comments := []models.FilmComment{}
	for rows.Next() {
		var c models.FilmComment
		if err := rows.Scan(&c.Id, &c.AuthorId, &c.AuthorName, &c.Contenu, &c.CreatedAt); err != nil {
			return nil, err
		}
		comments = append(comments, c)
	}
	return comments, nil
}

// Renvoie l'auteur d'un commentaire (0 si introuvable).
func (r *FilmRepository) CommentAuthor(commentId int) (int, error) {
	var authorId int
	err := r.db.QueryRow("SELECT `utilisateur_id` FROM `film_commentaires` WHERE `id` = ?;", commentId).Scan(&authorId)
	if err == sql.ErrNoRows {
		return 0, nil
	}
	if err != nil {
		return 0, fmt.Errorf("Erreur lecture commentaire - %s", err.Error())
	}
	return authorId, nil
}

// Supprime un commentaire.
func (r *FilmRepository) DeleteComment(commentId int) error {
	_, err := r.db.Exec("DELETE FROM `film_commentaires` WHERE `id` = ?;", commentId)
	if err != nil {
		return fmt.Errorf("Erreur suppression commentaire - %s", err.Error())
	}
	return nil
}
