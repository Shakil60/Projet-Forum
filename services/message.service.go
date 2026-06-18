package services

// Logique metier liee aux messages postes dans les fils de discussion.

import (
	"errors"
	"forum/dto"
	"forum/models"
	"forum/repositories"
	"strings"
)

type MessageService struct {
	messageRepository *repositories.MessageRepository
	threadRepository  *repositories.ThreadRepository
}

func InitMessageService(messageRepository *repositories.MessageRepository, threadRepository *repositories.ThreadRepository) *MessageService {
	return &MessageService{
		messageRepository: messageRepository,
		threadRepository:  threadRepository,
	}
}

// Cree un message apres avoir verifie que le fil existe et qu'il est encore ouvert.
func (s *MessageService) Create(authorId int, threadId int, contenu string) (int, error) {
	contenu = strings.TrimSpace(contenu)
	if contenu == "" {
		return -1, errors.New("le message ne peut pas etre vide")
	}

	thread, err := s.threadRepository.ReadById(threadId)
	if err != nil {
		return -1, err
	}
	if thread.Id == 0 {
		return -1, errors.New("fil introuvable")
	}
	if thread.Etat != "ouvert" {
		return -1, errors.New("ce fil n'accepte plus de nouveaux messages")
	}

	message := models.Message{
		ThreadId: threadId,
		AuthorId: authorId,
		Contenu:  contenu,
	}

	return s.messageRepository.Create(message)
}

// Renvoie les messages d'un fil avec leur pagination.
func (s *MessageService) ReadForThread(threadId int, sort string, currentUserId int, page int, size int) ([]models.Message, dto.Pagination, error) {
	total, err := s.messageRepository.CountByThread(threadId)
	if err != nil {
		return nil, dto.Pagination{}, err
	}

	pagination := dto.NewPagination(page, size, total)

	messages, messagesErr := s.messageRepository.ReadByThread(threadId, sort, currentUserId, pagination.Limit(), pagination.Offset())
	if messagesErr != nil {
		return nil, dto.Pagination{}, messagesErr
	}

	return messages, pagination, nil
}

// Recupere un message et verifie que l'utilisateur en est bien l'auteur.
func (s *MessageService) requireOwner(messageId int, userId int) (models.Message, error) {
	message, err := s.messageRepository.ReadById(messageId)
	if err != nil {
		return models.Message{}, err
	}
	if message.Id == 0 {
		return models.Message{}, errors.New("message introuvable")
	}
	if message.AuthorId != userId {
		return models.Message{}, errors.New("action non autorisee")
	}
	return message, nil
}

func (s *MessageService) GetForEdit(userId int, messageId int) (models.Message, error) {
	return s.requireOwner(messageId, userId)
}

// Modifie le contenu d'un message appartenant a l'utilisateur et renvoie l'id du fil.
func (s *MessageService) Update(userId int, messageId int, contenu string) (int, error) {
	message, err := s.requireOwner(messageId, userId)
	if err != nil {
		return 0, err
	}

	contenu = strings.TrimSpace(contenu)
	if contenu == "" {
		return 0, errors.New("le message ne peut pas etre vide")
	}

	if updateErr := s.messageRepository.Update(messageId, contenu); updateErr != nil {
		return 0, updateErr
	}
	return message.ThreadId, nil
}

// Supprime un message si l'utilisateur en est l'auteur ou s'il est administrateur.
func (s *MessageService) Delete(userId int, isAdmin bool, messageId int) (int, error) {
	message, err := s.messageRepository.ReadById(messageId)
	if err != nil {
		return 0, err
	}
	if message.Id == 0 {
		return 0, errors.New("message introuvable")
	}
	if !isAdmin && message.AuthorId != userId {
		return 0, errors.New("action non autorisee")
	}

	if deleteErr := s.messageRepository.Delete(messageId); deleteErr != nil {
		return 0, deleteErr
	}
	return message.ThreadId, nil
}
