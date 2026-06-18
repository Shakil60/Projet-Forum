package models

// Represente une reaction (like ou dislike) d'un utilisateur sur un message.

type Reaction struct {
	Id        int    `json:"id"`
	MessageId int    `json:"message_id"`
	UserId    int    `json:"utilisateur_id"`
	Type      string `json:"type"`
}
