package services

// Logique metier des reactions (like / dislike) sur les messages.

import (
	"errors"
	"forum/dto"
	"forum/repositories"
)

type ReactionService struct {
	reactionRepository *repositories.ReactionRepository
	messageRepository  *repositories.MessageRepository
}

func InitReactionService(reactionRepository *repositories.ReactionRepository, messageRepository *repositories.MessageRepository) *ReactionService {
	return &ReactionService{
		reactionRepository: reactionRepository,
		messageRepository:  messageRepository,
	}
}

// Ajoute, retire ou change la reaction d'un utilisateur puis renvoie les compteurs a jour.
func (s *ReactionService) Toggle(userId int, messageId int, reactionType string) (dto.ReactionResponseDto, error) {
	if reactionType != "like" && reactionType != "dislike" {
		return dto.ReactionResponseDto{}, errors.New("type de reaction invalide")
	}

	message, err := s.messageRepository.ReadById(messageId)
	if err != nil {
		return dto.ReactionResponseDto{}, err
	}
	if message.Id == 0 {
		return dto.ReactionResponseDto{}, errors.New("message introuvable")
	}

	current, findErr := s.reactionRepository.Find(messageId, userId)
	if findErr != nil {
		return dto.ReactionResponseDto{}, findErr
	}

	userReaction := reactionType
	if current == reactionType {
		if removeErr := s.reactionRepository.Remove(messageId, userId); removeErr != nil {
			return dto.ReactionResponseDto{}, removeErr
		}
		userReaction = ""
	} else {
		if setErr := s.reactionRepository.Set(messageId, userId, reactionType); setErr != nil {
			return dto.ReactionResponseDto{}, setErr
		}
	}

	likes, dislikes, countErr := s.reactionRepository.Counts(messageId)
	if countErr != nil {
		return dto.ReactionResponseDto{}, countErr
	}

	return dto.ReactionResponseDto{
		Likes:        likes,
		Dislikes:     dislikes,
		Score:        likes - dislikes,
		UserReaction: userReaction,
	}, nil
}
