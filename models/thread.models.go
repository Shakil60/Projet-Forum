package models

type Thread struct {
	Id           int    `json:"id"`
	Titre        string `json:"titre"`
	Contenu      string `json:"contenu"`
	AuthorId     int    `json:"utilisateur_id"`
	AuthorName   string `json:"auteur"`
	Etat         string `json:"etat"`
	CreatedAt    string `json:"date_creation"`
	Tags         []Tag  `json:"tags"`
	MessageCount int    `json:"nb_messages"`
}
