package services

// Logique metier des favoris et commentaires de films.

import (
	"errors"
	"forum/models"
	"forum/repositories"
	"strings"
)

type FilmService struct {
	filmRepository *repositories.FilmRepository
}

func InitFilmService(filmRepository *repositories.FilmRepository) *FilmService {
	return &FilmService{filmRepository: filmRepository}
}

// Verifie que le type de media est valide (movie ou tv).
func validMediaType(mediaType string) bool {
	return mediaType == "movie" || mediaType == "tv"
}

func (s *FilmService) IsFavori(userId, tmdbId int, mediaType string) bool {
	favori, _ := s.filmRepository.IsFavori(userId, tmdbId, mediaType)
	return favori
}

func (s *FilmService) CountFavoris(tmdbId int, mediaType string) int {
	n, _ := s.filmRepository.CountFavoris(tmdbId, mediaType)
	return n
}

// Ajoute ou retire un favori apres validation.
func (s *FilmService) ToggleFavori(userId, tmdbId int, mediaType, titre, affiche string) error {
	if !validMediaType(mediaType) {
		return errors.New("type de media invalide")
	}
	titre = strings.TrimSpace(titre)
	if titre == "" {
		titre = "Sans titre"
	}
	_, err := s.filmRepository.ToggleFavori(userId, tmdbId, mediaType, titre, affiche)
	return err
}

func (s *FilmService) ListFavoris(userId int) ([]models.FilmFavori, error) {
	return s.filmRepository.ListFavoris(userId)
}

func (s *FilmService) Comments(tmdbId int, mediaType string) []models.FilmComment {
	comments, _ := s.filmRepository.ListComments(tmdbId, mediaType)
	return comments
}

// Ajoute un commentaire apres validation du contenu.
func (s *FilmService) AddComment(tmdbId int, mediaType string, userId int, contenu string) error {
	if !validMediaType(mediaType) {
		return errors.New("type de media invalide")
	}
	contenu = strings.TrimSpace(contenu)
	if contenu == "" {
		return errors.New("le commentaire ne peut pas etre vide")
	}
	return s.filmRepository.AddComment(tmdbId, mediaType, userId, contenu)
}

// Supprime un commentaire si l'utilisateur en est l'auteur ou est admin.
func (s *FilmService) DeleteComment(commentId, userId int, isAdmin bool) error {
	authorId, err := s.filmRepository.CommentAuthor(commentId)
	if err != nil {
		return err
	}
	if authorId == 0 {
		return errors.New("commentaire introuvable")
	}
	if authorId != userId && !isAdmin {
		return errors.New("action non autorisee")
	}
	return s.filmRepository.DeleteComment(commentId)
}
