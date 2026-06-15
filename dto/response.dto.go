package dto

type ResponseDto struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type ReactionResponseDto struct {
	Likes        int    `json:"likes"`
	Dislikes     int    `json:"dislikes"`
	Score        int    `json:"score"`
	UserReaction string `json:"reaction_utilisateur"`
}
