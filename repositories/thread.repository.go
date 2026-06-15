package repositories

import (
	"database/sql"
	"forum/models"
	"fmt"
)

type ThreadRepository struct {
	db *sql.DB
}

func InitThreadRepository(db *sql.DB) *ThreadRepository {
	return &ThreadRepository{db}
}

func (r *ThreadRepository) Create(thread models.Thread) (int, error) {
	query := "INSERT INTO `fils`(`titre`, `contenu`, `utilisateur_id`, `etat`) VALUES (?,?,?,?);"

	result, err := r.db.Exec(query, thread.Titre, thread.Contenu, thread.AuthorId, thread.Etat)
	if err != nil {
		return -1, fmt.Errorf("Erreur creation fil - %s", err.Error())
	}

	id, idErr := result.LastInsertId()
	if idErr != nil {
		return -1, fmt.Errorf("Erreur creation fil - recuperation identifiant : %s", idErr.Error())
	}

	return int(id), nil
}

func (r *ThreadRepository) SetTags(threadId int, tagIds []int) error {
	if _, err := r.db.Exec("DELETE FROM `fil_tags` WHERE `fil_id` = ?;", threadId); err != nil {
		return fmt.Errorf("Erreur association tags - %s", err.Error())
	}

	for _, tagId := range tagIds {
		if _, err := r.db.Exec("INSERT INTO `fil_tags`(`fil_id`, `tag_id`) VALUES (?,?);", threadId, tagId); err != nil {
			return fmt.Errorf("Erreur association tags - %s", err.Error())
		}
	}

	return nil
}

func (r *ThreadRepository) loadTags(threadId int) ([]models.Tag, error) {
	var tags []models.Tag
	query := "SELECT t.`id`, t.`nom` FROM `tags` t JOIN `fil_tags` ft ON ft.`tag_id` = t.`id` WHERE ft.`fil_id` = ? ORDER BY t.`nom`;"

	rows, err := r.db.Query(query, threadId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var tag models.Tag
		if scanErr := rows.Scan(&tag.Id, &tag.Nom); scanErr != nil {
			return nil, scanErr
		}
		tags = append(tags, tag)
	}

	return tags, nil
}

func (r *ThreadRepository) ReadById(id int) (models.Thread, error) {
	var thread models.Thread
	query := "SELECT f.`id`, f.`titre`, f.`contenu`, f.`utilisateur_id`, u.`nom_utilisateur`, f.`etat`, f.`date_creation`, " +
		"(SELECT COUNT(*) FROM `messages` m WHERE m.`fil_id` = f.`id`) " +
		"FROM `fils` f JOIN `utilisateurs` u ON u.`id` = f.`utilisateur_id` WHERE f.`id` = ?;"

	err := r.db.QueryRow(query, id).Scan(
		&thread.Id, &thread.Titre, &thread.Contenu, &thread.AuthorId,
		&thread.AuthorName, &thread.Etat, &thread.CreatedAt, &thread.MessageCount,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.Thread{}, nil
		}
		return models.Thread{}, fmt.Errorf("Erreur recuperation fil - %s", err.Error())
	}

	tags, tagsErr := r.loadTags(thread.Id)
	if tagsErr != nil {
		return models.Thread{}, tagsErr
	}
	thread.Tags = tags

	return thread, nil
}

func (r *ThreadRepository) buildFilters(tagId int, search string) (string, []any) {
	clause := " WHERE f.`etat` <> 'archive'"
	args := []any{}

	if tagId > 0 {
		clause += " AND f.`id` IN (SELECT `fil_id` FROM `fil_tags` WHERE `tag_id` = ?)"
		args = append(args, tagId)
	}

	if search != "" {
		like := "%" + search + "%"
		clause += " AND (f.`titre` LIKE ? OR f.`id` IN (SELECT ft.`fil_id` FROM `fil_tags` ft JOIN `tags` t ON t.`id` = ft.`tag_id` WHERE t.`nom` LIKE ?))"
		args = append(args, like, like)
	}

	return clause, args
}

func (r *ThreadRepository) CountVisible(tagId int, search string) (int, error) {
	clause, args := r.buildFilters(tagId, search)
	query := "SELECT COUNT(*) FROM `fils` f" + clause + ";"

	var count int
	if err := r.db.QueryRow(query, args...).Scan(&count); err != nil {
		return 0, fmt.Errorf("Erreur comptage fils - %s", err.Error())
	}
	return count, nil
}

func (r *ThreadRepository) ReadVisible(tagId int, search string, limit int, offset int) ([]models.Thread, error) {
	clause, args := r.buildFilters(tagId, search)

	query := "SELECT f.`id`, f.`titre`, f.`contenu`, f.`utilisateur_id`, u.`nom_utilisateur`, f.`etat`, f.`date_creation`, " +
		"(SELECT COUNT(*) FROM `messages` m WHERE m.`fil_id` = f.`id`) " +
		"FROM `fils` f JOIN `utilisateurs` u ON u.`id` = f.`utilisateur_id`" + clause +
		" ORDER BY f.`date_creation` DESC"

	if limit > 0 {
		query += " LIMIT ? OFFSET ?"
		args = append(args, limit, offset)
	}
	query += ";"

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("Erreur recuperation fils - %s", err.Error())
	}
	defer rows.Close()

	var threads []models.Thread
	for rows.Next() {
		var thread models.Thread
		if scanErr := rows.Scan(
			&thread.Id, &thread.Titre, &thread.Contenu, &thread.AuthorId,
			&thread.AuthorName, &thread.Etat, &thread.CreatedAt, &thread.MessageCount,
		); scanErr != nil {
			return nil, scanErr
		}
		threads = append(threads, thread)
	}

	for i := range threads {
		tags, tagsErr := r.loadTags(threads[i].Id)
		if tagsErr != nil {
			return nil, tagsErr
		}
		threads[i].Tags = tags
	}

	return threads, nil
}

func (r *ThreadRepository) ReadAllForAdmin() ([]models.Thread, error) {
	query := "SELECT f.`id`, f.`titre`, f.`contenu`, f.`utilisateur_id`, u.`nom_utilisateur`, f.`etat`, f.`date_creation`, " +
		"(SELECT COUNT(*) FROM `messages` m WHERE m.`fil_id` = f.`id`) " +
		"FROM `fils` f JOIN `utilisateurs` u ON u.`id` = f.`utilisateur_id` ORDER BY f.`date_creation` DESC;"

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("Erreur recuperation fils - %s", err.Error())
	}
	defer rows.Close()

	var threads []models.Thread
	for rows.Next() {
		var thread models.Thread
		if scanErr := rows.Scan(
			&thread.Id, &thread.Titre, &thread.Contenu, &thread.AuthorId,
			&thread.AuthorName, &thread.Etat, &thread.CreatedAt, &thread.MessageCount,
		); scanErr != nil {
			return nil, scanErr
		}
		threads = append(threads, thread)
	}

	return threads, nil
}

func (r *ThreadRepository) Update(id int, titre string, contenu string) error {
	result, err := r.db.Exec("UPDATE `fils` SET `titre` = ?, `contenu` = ? WHERE `id` = ?;", titre, contenu, id)
	if err != nil {
		return fmt.Errorf("Erreur modification fil - %s", err.Error())
	}
	if affected, _ := result.RowsAffected(); affected <= 0 {
		return fmt.Errorf("Erreur modification fil - aucune ligne modifiee")
	}
	return nil
}

func (r *ThreadRepository) UpdateState(id int, etat string) error {
	result, err := r.db.Exec("UPDATE `fils` SET `etat` = ? WHERE `id` = ?;", etat, id)
	if err != nil {
		return fmt.Errorf("Erreur changement etat fil - %s", err.Error())
	}
	if affected, _ := result.RowsAffected(); affected <= 0 {
		return fmt.Errorf("Erreur changement etat fil - aucune ligne modifiee")
	}
	return nil
}

func (r *ThreadRepository) Delete(id int) error {
	result, err := r.db.Exec("DELETE FROM `fils` WHERE `id` = ?;", id)
	if err != nil {
		return fmt.Errorf("Erreur suppression fil - %s", err.Error())
	}
	if affected, _ := result.RowsAffected(); affected <= 0 {
		return fmt.Errorf("Erreur suppression fil - aucun fil supprime")
	}
	return nil
}
