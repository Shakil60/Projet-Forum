package repositories

// Acces a la base de donnees pour les utilisateurs.

import (
	"database/sql"
	"forum/models"
	"fmt"
)

type UserRepository struct {
	db *sql.DB
}

func InitUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db}
}

// Insere un nouvel utilisateur et renvoie son identifiant.
func (r *UserRepository) Create(user models.User) (int, error) {
	query := "INSERT INTO `utilisateurs`(`nom_utilisateur`, `email`, `mot_de_passe`, `sel`, `role`) VALUES (?,?,?,?,?);"

	result, err := r.db.Exec(query, user.Username, user.Email, user.Password, user.Salt, user.Role)
	if err != nil {
		return -1, fmt.Errorf("Erreur creation utilisateur - %s", err.Error())
	}

	id, idErr := result.LastInsertId()
	if idErr != nil {
		return -1, fmt.Errorf("Erreur creation utilisateur - recuperation identifiant : %s", idErr.Error())
	}

	return int(id), nil
}

// Indique si un nom d'utilisateur est deja pris.
func (r *UserRepository) ExistsByUsername(username string) (bool, error) {
	var count int
	err := r.db.QueryRow("SELECT COUNT(*) FROM `utilisateurs` WHERE `nom_utilisateur` = ?;", username).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// Indique si un email est deja utilise.
func (r *UserRepository) ExistsByEmail(email string) (bool, error) {
	var count int
	err := r.db.QueryRow("SELECT COUNT(*) FROM `utilisateurs` WHERE `email` = ?;", email).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// Recherche un utilisateur par nom d'utilisateur ou email (pour la connexion).
func (r *UserRepository) FindByLogin(identifiant string) (models.User, error) {
	var user models.User
	query := "SELECT `id`, `nom_utilisateur`, `email`, `mot_de_passe`, `sel`, `role`, `banni`, `date_creation` FROM `utilisateurs` WHERE `nom_utilisateur` = ? OR `email` = ?;"

	err := r.db.QueryRow(query, identifiant, identifiant).
		Scan(&user.Id, &user.Username, &user.Email, &user.Password, &user.Salt, &user.Role, &user.Banned, &user.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.User{}, nil
		}
		return models.User{}, fmt.Errorf("Erreur recuperation utilisateur - %s", err.Error())
	}

	return user, nil
}

// Recupere un utilisateur par son identifiant.
func (r *UserRepository) FindById(id int) (models.User, error) {
	var user models.User
	query := "SELECT `id`, `nom_utilisateur`, `email`, `mot_de_passe`, `sel`, `role`, `banni`, `date_creation` FROM `utilisateurs` WHERE `id` = ?;"

	err := r.db.QueryRow(query, id).
		Scan(&user.Id, &user.Username, &user.Email, &user.Password, &user.Salt, &user.Role, &user.Banned, &user.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.User{}, nil
		}
		return models.User{}, fmt.Errorf("Erreur recuperation utilisateur - %s", err.Error())
	}

	return user, nil
}

// Recupere tous les utilisateurs (pour l'administration).
func (r *UserRepository) ReadAll() ([]models.User, error) {
	var users []models.User
	query := "SELECT `id`, `nom_utilisateur`, `email`, `role`, `banni`, `date_creation` FROM `utilisateurs` ORDER BY `date_creation` DESC;"

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("Erreur recuperation utilisateurs - %s", err.Error())
	}
	defer rows.Close()

	for rows.Next() {
		var user models.User
		if scanErr := rows.Scan(&user.Id, &user.Username, &user.Email, &user.Role, &user.Banned, &user.CreatedAt); scanErr != nil {
			return nil, scanErr
		}
		users = append(users, user)
	}

	return users, nil
}

// Bannit ou debannit un utilisateur.
func (r *UserRepository) SetBanned(id int, banned bool) error {
	result, err := r.db.Exec("UPDATE `utilisateurs` SET `banni` = ? WHERE `id` = ?;", banned, id)
	if err != nil {
		return fmt.Errorf("Erreur mise a jour utilisateur - %s", err.Error())
	}
	if affected, _ := result.RowsAffected(); affected <= 0 {
		return fmt.Errorf("Erreur mise a jour utilisateur - aucun utilisateur modifie")
	}
	return nil
}
