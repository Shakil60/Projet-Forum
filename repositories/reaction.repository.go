package repositories

import (
	"database/sql"
	"fmt"
)

type ReactionRepository struct {
	db *sql.DB
}

func InitReactionRepository(db *sql.DB) *ReactionRepository {
	return &ReactionRepository{db}
}

func (r *ReactionRepository) Find(messageId int, userId int) (string, error) {
	var reactionType string
	err := r.db.QueryRow("SELECT `type` FROM `reactions` WHERE `message_id` = ? AND `utilisateur_id` = ?;", messageId, userId).Scan(&reactionType)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", nil
		}
		return "", fmt.Errorf("Erreur recuperation reaction - %s", err.Error())
	}
	return reactionType, nil
}

func (r *ReactionRepository) Set(messageId int, userId int, reactionType string) error {
	query := "INSERT INTO `reactions`(`message_id`, `utilisateur_id`, `type`) VALUES (?,?,?) " +
		"ON DUPLICATE KEY UPDATE `type` = VALUES(`type`);"
	if _, err := r.db.Exec(query, messageId, userId, reactionType); err != nil {
		return fmt.Errorf("Erreur enregistrement reaction - %s", err.Error())
	}
	return nil
}

func (r *ReactionRepository) Remove(messageId int, userId int) error {
	if _, err := r.db.Exec("DELETE FROM `reactions` WHERE `message_id` = ? AND `utilisateur_id` = ?;", messageId, userId); err != nil {
		return fmt.Errorf("Erreur suppression reaction - %s", err.Error())
	}
	return nil
}

func (r *ReactionRepository) Counts(messageId int) (int, int, error) {
	var likes, dislikes int
	query := "SELECT " +
		"COALESCE(SUM(CASE WHEN `type` = 'like' THEN 1 ELSE 0 END), 0), " +
		"COALESCE(SUM(CASE WHEN `type` = 'dislike' THEN 1 ELSE 0 END), 0) " +
		"FROM `reactions` WHERE `message_id` = ?;"
	if err := r.db.QueryRow(query, messageId).Scan(&likes, &dislikes); err != nil {
		return 0, 0, fmt.Errorf("Erreur comptage reactions - %s", err.Error())
	}
	return likes, dislikes, nil
}
