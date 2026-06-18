package services

// Logique metier reservee a l'administration des utilisateurs.

import (
	"errors"
	"forum/models"
	"forum/repositories"
)

type AdminService struct {
	userRepository *repositories.UserRepository
}

func InitAdminService(userRepository *repositories.UserRepository) *AdminService {
	return &AdminService{userRepository: userRepository}
}

func (s *AdminService) ListUsers() ([]models.User, error) {
	return s.userRepository.ReadAll()
}

// Bannit ou debannit un utilisateur en verifiant les regles (pas soi-meme, pas un admin).
func (s *AdminService) SetBanned(actingUserId int, targetUserId int, banned bool) error {
	if actingUserId == targetUserId {
		return errors.New("vous ne pouvez pas bannir votre propre compte")
	}

	target, err := s.userRepository.FindById(targetUserId)
	if err != nil {
		return err
	}
	if target.Id == 0 {
		return errors.New("utilisateur introuvable")
	}
	if target.IsAdmin() {
		return errors.New("impossible de bannir un administrateur")
	}

	return s.userRepository.SetBanned(targetUserId, banned)
}
