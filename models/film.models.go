package models

// Donnees liees aux films/series du catalogue (favoris et commentaires).

// Un commentaire poste par un membre sous un film ou une serie.
type FilmComment struct {
	Id         int
	TmdbId     int
	MediaType  string
	AuthorId   int
	AuthorName string
	Contenu    string
	CreatedAt  string
}

// Un film/serie mis en favori par un membre.
type FilmFavori struct {
	TmdbId    int
	MediaType string
	Titre     string
	Affiche   string
	CreatedAt string
}
