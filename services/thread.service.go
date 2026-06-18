package services

import (
	"errors"
	"forum/dto"
	"forum/models"
	"forum/repositories"
	"strings"
)

var validStates = map[string]bool{
	"ouvert":  true,
	"ferme":   true,
	"archive": true,
}

type ThreadService struct {
	threadRepository *repositories.ThreadRepository
	tagRepository    *repositories.TagRepository
}

func InitThreadService(threadRepository *repositories.ThreadRepository, tagRepository *repositories.TagRepository) *ThreadService {
	return &ThreadService{
		threadRepository: threadRepository,
		tagRepository:    tagRepository,
	}
}

func (s *ThreadService) Create(authorId int, titre string, contenu string, tagIds []int, newTags []string) (int, error) {
	titre = strings.TrimSpace(titre)
	contenu = strings.TrimSpace(contenu)

	if titre == "" || contenu == "" {
		return -1, errors.New("le titre et le contenu sont obligatoires")
	}

	// Un fil doit etre associe a au moins un genre (regle de gestion FT-3)
	resolvedTags, tagErr := s.resolveTags(tagIds, newTags)
	if tagErr != nil {
		return -1, tagErr
	}
	if len(resolvedTags) == 0 {
		return -1, errors.New("un fil doit etre associe a au moins un genre")
	}

	thread := models.Thread{
		Titre:    titre,
		Contenu:  contenu,
		AuthorId: authorId,
		Etat:     "ouvert",
	}

	threadId, err := s.threadRepository.Create(thread)
	if err != nil {
		return -1, err
	}

	if setErr := s.threadRepository.SetTags(threadId, resolvedTags); setErr != nil {
		return -1, setErr
	}

	return threadId, nil
}

func (s *ThreadService) resolveTags(tagIds []int, newTags []string) ([]int, error) {
	unique := map[int]bool{}
	resolved := []int{}

	for _, id := range tagIds {
		if id > 0 && !unique[id] {
			unique[id] = true
			resolved = append(resolved, id)
		}
	}

	for _, name := range newTags {
		name = strings.TrimSpace(name)
		if name == "" {
			continue
		}
		id, err := s.tagRepository.FindOrCreate(name)
		if err != nil {
			return nil, err
		}
		if !unique[id] {
			unique[id] = true
			resolved = append(resolved, id)
		}
	}

	return resolved, nil
}

func (s *ThreadService) ReadForList(tagId int, search string, page int, size int) ([]models.Thread, dto.Pagination, error) {
	total, err := s.threadRepository.CountVisible(tagId, search)
	if err != nil {
		return nil, dto.Pagination{}, err
	}

	pagination := dto.NewPagination(page, size, total)

	threads, threadsErr := s.threadRepository.ReadVisible(tagId, search, pagination.Limit(), pagination.Offset())
	if threadsErr != nil {
		return nil, dto.Pagination{}, threadsErr
	}

	return threads, pagination, nil
}

func (s *ThreadService) ReadVisibleById(id int) (models.Thread, error) {
	thread, err := s.threadRepository.ReadById(id)
	if err != nil {
		return models.Thread{}, err
	}
	if thread.Id == 0 {
		return models.Thread{}, errors.New("fil introuvable")
	}
	if thread.Etat == "archive" {
		return models.Thread{}, errors.New("fil introuvable")
	}
	return thread, nil
}

func (s *ThreadService) ReadById(id int) (models.Thread, error) {
	thread, err := s.threadRepository.ReadById(id)
	if err != nil {
		return models.Thread{}, err
	}
	if thread.Id == 0 {
		return models.Thread{}, errors.New("fil introuvable")
	}
	return thread, nil
}

func (s *ThreadService) requireOwnerOrAdmin(threadId int, userId int, isAdmin bool) (models.Thread, error) {
	thread, err := s.threadRepository.ReadById(threadId)
	if err != nil {
		return models.Thread{}, err
	}
	if thread.Id == 0 {
		return models.Thread{}, errors.New("fil introuvable")
	}
	if !isAdmin && thread.AuthorId != userId {
		return models.Thread{}, errors.New("action non autorisee")
	}
	return thread, nil
}

func (s *ThreadService) Update(userId int, isAdmin bool, threadId int, titre string, contenu string, tagIds []int, newTags []string) error {
	if _, err := s.requireOwnerOrAdmin(threadId, userId, isAdmin); err != nil {
		return err
	}

	titre = strings.TrimSpace(titre)
	contenu = strings.TrimSpace(contenu)
	if titre == "" || contenu == "" {
		return errors.New("le titre et le contenu sont obligatoires")
	}

	// Un fil doit rester associe a au moins un genre (regle de gestion FT-3)
	resolvedTags, tagErr := s.resolveTags(tagIds, newTags)
	if tagErr != nil {
		return tagErr
	}
	if len(resolvedTags) == 0 {
		return errors.New("un fil doit etre associe a au moins un genre")
	}

	if err := s.threadRepository.Update(threadId, titre, contenu); err != nil {
		return err
	}

	return s.threadRepository.SetTags(threadId, resolvedTags)
}

func (s *ThreadService) ChangeState(userId int, isAdmin bool, threadId int, etat string) error {
	if !validStates[etat] {
		return errors.New("etat invalide")
	}
	if _, err := s.requireOwnerOrAdmin(threadId, userId, isAdmin); err != nil {
		return err
	}
	return s.threadRepository.UpdateState(threadId, etat)
}

func (s *ThreadService) Delete(userId int, isAdmin bool, threadId int) error {
	if _, err := s.requireOwnerOrAdmin(threadId, userId, isAdmin); err != nil {
		return err
	}
	return s.threadRepository.Delete(threadId)
}

func (s *ThreadService) GetForEdit(userId int, isAdmin bool, threadId int) (models.Thread, error) {
	return s.requireOwnerOrAdmin(threadId, userId, isAdmin)
}

func (s *ThreadService) ReadAllForAdmin() ([]models.Thread, error) {
	return s.threadRepository.ReadAllForAdmin()
}
