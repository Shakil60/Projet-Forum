package services

// Logique metier liee aux tags (categories) des fils de discussion.

import (
	"forum/models"
	"forum/repositories"
)

type TagService struct {
	tagRepository *repositories.TagRepository
}

func InitTagService(tagRepository *repositories.TagRepository) *TagService {
	return &TagService{tagRepository: tagRepository}
}

func (s *TagService) ReadAll() ([]models.Tag, error) {
	return s.tagRepository.ReadAll()
}

func (s *TagService) FindById(id int) (models.Tag, error) {
	return s.tagRepository.FindById(id)
}
