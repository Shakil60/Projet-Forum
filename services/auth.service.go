package services

// Logique metier d'authentification : inscription et connexion des utilisateurs.

import (
	"errors"
	"forum/auth"
	"forum/models"
	"forum/repositories"
	"strings"
)

type AuthService struct {
	userRepository *repositories.UserRepository
}

func InitAuthService(userRepository *repositories.UserRepository) *AuthService {
	return &AuthService{userRepository: userRepository}
}

// Valide les champs, verifie l'unicite puis cree un nouvel utilisateur avec mot de passe hache.
func (s *AuthService) Register(username string, email string, password string) (models.User, error) {
	username = strings.TrimSpace(username)
	email = strings.TrimSpace(email)

	if username == "" || email == "" || password == "" {
		return models.User{}, errors.New("tous les champs sont obligatoires")
	}
	if !strings.Contains(email, "@") || !strings.Contains(email, ".") {
		return models.User{}, errors.New("adresse e-mail invalide")
	}
	if err := auth.ValidatePasswordRules(password); err != nil {
		return models.User{}, err
	}

	usernameTaken, err := s.userRepository.ExistsByUsername(username)
	if err != nil {
		return models.User{}, err
	}
	if usernameTaken {
		return models.User{}, errors.New("ce nom d'utilisateur est deja utilise")
	}

	emailTaken, err := s.userRepository.ExistsByEmail(email)
	if err != nil {
		return models.User{}, err
	}
	if emailTaken {
		return models.User{}, errors.New("cette adresse e-mail est deja utilisee")
	}

	salt, saltErr := auth.GenerateSalt()
	if saltErr != nil {
		return models.User{}, saltErr
	}

	user := models.User{
		Username: username,
		Email:    email,
		Password: auth.HashPassword(password, salt),
		Salt:     salt,
		Role:     "utilisateur",
	}

	id, createErr := s.userRepository.Create(user)
	if createErr != nil {
		return models.User{}, createErr
	}
	user.Id = id

	return user, nil
}

// Verifie les identifiants et renvoie un jeton si la connexion est valide.
func (s *AuthService) Login(identifiant string, password string) (string, models.User, error) {
	identifiant = strings.TrimSpace(identifiant)
	if identifiant == "" || password == "" {
		return "", models.User{}, errors.New("identifiant et mot de passe obligatoires")
	}

	user, err := s.userRepository.FindByLogin(identifiant)
	if err != nil {
		return "", models.User{}, err
	}
	if user.Id == 0 {
		return "", models.User{}, errors.New("identifiants invalides")
	}
	if !auth.CheckPassword(password, user.Salt, user.Password) {
		return "", models.User{}, errors.New("identifiants invalides")
	}
	if user.Banned {
		return "", models.User{}, errors.New("ce compte a ete banni")
	}

	token, tokenErr := auth.GenerateToken(user.Id, user.Username, user.Role)
	if tokenErr != nil {
		return "", models.User{}, tokenErr
	}

	return token, user, nil
}
