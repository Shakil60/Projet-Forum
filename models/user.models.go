package models

type User struct {
	Id        int    `json:"id"`
	Username  string `json:"nom_utilisateur"`
	Email     string `json:"email"`
	Password  string `json:"-"`
	Salt      string `json:"-"`
	Role      string `json:"role"`
	Banned    bool   `json:"banni"`
	CreatedAt string `json:"date_creation"`
}

func (u User) IsAdmin() bool {
	return u.Role == "administrateur"
}
