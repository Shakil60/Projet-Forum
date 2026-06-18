package auth

// Gestion des mots de passe : hachage, sel et regles de validation.

import (
	"crypto/rand"
	"crypto/sha512"
	"encoding/hex"
	"errors"
	"unicode"
)

// Genere un sel aleatoire de 16 octets en hexadecimal.
func GenerateSalt() (string, error) {
	buffer := make([]byte, 16)
	if _, err := rand.Read(buffer); err != nil {
		return "", err
	}
	return hex.EncodeToString(buffer), nil
}

// Hache le mot de passe avec son sel en SHA-512.
func HashPassword(password string, salt string) string {
	sum := sha512.Sum512([]byte(salt + password))
	return hex.EncodeToString(sum[:])
}

// Verifie qu'un mot de passe correspond bien au hash stocke.
func CheckPassword(password string, salt string, hash string) bool {
	return HashPassword(password, salt) == hash
}

// Verifie que le mot de passe respecte les regles (longueur, majuscule, caractere special).
func ValidatePasswordRules(password string) error {
	if len(password) < 12 {
		return errors.New("le mot de passe doit contenir au moins 12 caracteres")
	}

	hasUpper := false
	hasSpecial := false
	for _, c := range password {
		switch {
		case unicode.IsUpper(c):
			hasUpper = true
		case !unicode.IsLetter(c) && !unicode.IsDigit(c):
			hasSpecial = true
		}
	}

	if !hasUpper {
		return errors.New("le mot de passe doit contenir au moins une majuscule")
	}
	if !hasSpecial {
		return errors.New("le mot de passe doit contenir au moins un caractere special")
	}

	return nil
}
