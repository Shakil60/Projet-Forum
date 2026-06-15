package repositories

import (
	"database/sql"
	"forum/models"
	"fmt"
)

type MessageRepository struct {
	db *sql.DB
}

func InitMessageRepository(db *sql.DB) *MessageRepository {
	return &MessageRepository{db}
}

func (r *MessageRepository) Create(message models.Message) (int, error) {
	query := "INSERT INTO `messages`(`fil_id`, `utilisateur_id`, `contenu`) VALUES (?,?,?);"

	result, err := r.db.Exec(query, message.ThreadId, message.AuthorId, message.Contenu)
	if err != nil {
		return -1, fmt.Errorf("Erreur publication message - %s", err.Error())
	}

	id, idErr := result.LastInsertId()
	if idErr != nil {
		return -1, fmt.Errorf("Erreur publication message - recuperation identifiant : %s", idErr.Error())
	}

	return int(id), nil
}

func orderClause(sort string) string {
	switch sort {
	case "ancien":
		return " ORDER BY m.`date_envoi` ASC, m.`id` ASC"
	case "populaire":
		return " ORDER BY score DESC, m.`date_envoi` DESC"
	default:
		return " ORDER BY m.`date_envoi` DESC, m.`id` DESC"
	}
}

func (r *MessageRepository) CountByThread(threadId int) (int, error) {
	var count int
	err := r.db.QueryRow("SELECT COUNT(*) FROM `messages` WHERE `fil_id` = ?;", threadId).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("Erreur comptage messages - %s", err.Error())
	}
	return count, nil
}

func (r *MessageRepository) ReadByThread(threadId int, sort string, currentUserId int, limit int, offset int) ([]models.Message, error) {
	query := "SELECT m.`id`, m.`fil_id`, m.`utilisateur_id`, u.`nom_utilisateur`, m.`contenu`, m.`date_envoi`, " +
		"COALESCE(SUM(CASE WHEN r.`type` = 'like' THEN 1 ELSE 0 END), 0) AS likes, " +
		"COALESCE(SUM(CASE WHEN r.`type` = 'dislike' THEN 1 ELSE 0 END), 0) AS dislikes, " +
		"COALESCE(SUM(CASE WHEN r.`type` = 'like' THEN 1 WHEN r.`type` = 'dislike' THEN -1 ELSE 0 END), 0) AS score, " +
		"COALESCE(MAX(CASE WHEN r.`utilisateur_id` = ? THEN r.`type` END), '') AS reaction_utilisateur " +
		"FROM `messages` m " +
		"JOIN `utilisateurs` u ON u.`id` = m.`utilisateur_id` " +
		"LEFT JOIN `reactions` r ON r.`message_id` = m.`id` " +
		"WHERE m.`fil_id` = ? " +
		"GROUP BY m.`id`, m.`fil_id`, m.`utilisateur_id`, u.`nom_utilisateur`, m.`contenu`, m.`date_envoi`" +
		orderClause(sort)

	args := []any{currentUserId, threadId}
	if limit > 0 {
		query += " LIMIT ? OFFSET ?"
		args = append(args, limit, offset)
	}
	query += ";"

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("Erreur recuperation messages - %s", err.Error())
	}
	defer rows.Close()

	var messages []models.Message
	for rows.Next() {
		var message models.Message
		if scanErr := rows.Scan(
			&message.Id, &message.ThreadId, &message.AuthorId, &message.AuthorName,
			&message.Contenu, &message.CreatedAt, &message.Likes, &message.Dislikes,
			&message.Score, &message.UserReaction,
		); scanErr != nil {
			return nil, scanErr
		}
		messages = append(messages, message)
	}

	return messages, nil
}

func (r *MessageRepository) ReadById(id int) (models.Message, error) {
	var message models.Message
	query := "SELECT m.`id`, m.`fil_id`, m.`utilisateur_id`, u.`nom_utilisateur`, m.`contenu`, m.`date_envoi` " +
		"FROM `messages` m JOIN `utilisateurs` u ON u.`id` = m.`utilisateur_id` WHERE m.`id` = ?;"

	err := r.db.QueryRow(query, id).Scan(
		&message.Id, &message.ThreadId, &message.AuthorId,
		&message.AuthorName, &message.Contenu, &message.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.Message{}, nil
		}
		return models.Message{}, fmt.Errorf("Erreur recuperation message - %s", err.Error())
	}

	return message, nil
}

func (r *MessageRepository) Update(id int, contenu string) error {
	result, err := r.db.Exec("UPDATE `messages` SET `contenu` = ? WHERE `id` = ?;", contenu, id)
	if err != nil {
		return fmt.Errorf("Erreur modification message - %s", err.Error())
	}
	if affected, _ := result.RowsAffected(); affected <= 0 {
		return fmt.Errorf("Erreur modification message - aucune ligne modifiee")
	}
	return nil
}

func (r *MessageRepository) Delete(id int) error {
	result, err := r.db.Exec("DELETE FROM `messages` WHERE `id` = ?;", id)
	if err != nil {
		return fmt.Errorf("Erreur suppression message - %s", err.Error())
	}
	if affected, _ := result.RowsAffected(); affected <= 0 {
		return fmt.Errorf("Erreur suppression message - aucun message supprime")
	}
	return nil
}
