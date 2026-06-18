package models

// Represente un message poste dans un fil de discussion.

type Message struct {
	Id           int    `json:"id"`
	ThreadId     int    `json:"fil_id"`
	AuthorId     int    `json:"utilisateur_id"`
	AuthorName   string `json:"auteur"`
	Contenu      string `json:"contenu"`
	CreatedAt    string `json:"date_envoi"`
	Likes        int    `json:"likes"`
	Dislikes     int    `json:"dislikes"`
	Score        int    `json:"score"`
	UserReaction string `json:"reaction_utilisateur"`
}
