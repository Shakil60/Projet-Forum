package services

// Service d'appel a l'API externe TMDB pour les films, series et personnes.

import (
	"encoding/json"
	"errors"
	"fmt"
	"forum/dto"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const tmdbBaseURL = "https://api.themoviedb.org/3"

type TMDBService struct {
	apiKey     string
	httpClient *http.Client
}

func InitTMDBService(apiKey string) *TMDBService {
	return &TMDBService{
		apiKey: strings.TrimSpace(apiKey),
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (s *TMDBService) IsConfigured() bool {
	return s.apiKey != ""
}

// Renvoie une erreur explicite si la cle API n'est pas configuree.
func (s *TMDBService) ConfigError() error {
	if s.IsConfigured() {
		return nil
	}
	return errors.New("cle API TMDB manquante : ajoutez TMDB_API_KEY dans le fichier .env")
}

// Construit l'URL, ajoute la cle et la langue, puis effectue l'appel HTTP vers TMDB.
func (s *TMDBService) request(path string, params url.Values) ([]byte, error) {
	if err := s.ConfigError(); err != nil {
		return nil, err
	}

	if params == nil {
		params = url.Values{}
	}
	params.Set("api_key", s.apiKey)
	params.Set("language", "fr-FR")

	requestURL := fmt.Sprintf("%s%s?%s", tmdbBaseURL, path, params.Encode())

	response, err := s.httpClient.Get(requestURL)
	if err != nil {
		return nil, fmt.Errorf("erreur reseau TMDB : %w", err)
	}
	defer response.Body.Close()

	body, readErr := io.ReadAll(response.Body)
	if readErr != nil {
		return nil, fmt.Errorf("erreur lecture reponse TMDB : %w", readErr)
	}

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("erreur TMDB (%d) : %s", response.StatusCode, string(body))
	}

	return body, nil
}

// Decode le JSON renvoye par TMDB dans la structure cible (fonction generique).
func decode[T any](body []byte, target *T) error {
	if err := json.Unmarshal(body, target); err != nil {
		return fmt.Errorf("erreur decodage TMDB : %w", err)
	}
	return nil
}

func (s *TMDBService) GetPopularMovies(page int) (dto.TMDBPagedResponse[dto.TMDBMediaItem], error) {
	var result dto.TMDBPagedResponse[dto.TMDBMediaItem]
	body, err := s.request("/movie/popular", url.Values{"page": {strconv.Itoa(page)}})
	if err != nil {
		return result, err
	}
	return result, decode(body, &result)
}

func (s *TMDBService) GetPopularSeries(page int) (dto.TMDBPagedResponse[dto.TMDBMediaItem], error) {
	var result dto.TMDBPagedResponse[dto.TMDBMediaItem]
	body, err := s.request("/tv/popular", url.Values{"page": {strconv.Itoa(page)}})
	if err != nil {
		return result, err
	}
	return result, decode(body, &result)
}

func (s *TMDBService) GetPopularPeople(page int) (dto.TMDBPagedResponse[dto.TMDBMediaItem], error) {
	var result dto.TMDBPagedResponse[dto.TMDBMediaItem]
	body, err := s.request("/person/popular", url.Values{"page": {strconv.Itoa(page)}})
	if err != nil {
		return result, err
	}
	return result, decode(body, &result)
}

func (s *TMDBService) SearchMovies(query string, page int) (dto.TMDBPagedResponse[dto.TMDBMediaItem], error) {
	var result dto.TMDBPagedResponse[dto.TMDBMediaItem]
	body, err := s.request("/search/movie", url.Values{
		"query": {query},
		"page":  {strconv.Itoa(page)},
	})
	if err != nil {
		return result, err
	}
	return result, decode(body, &result)
}

func (s *TMDBService) SearchSeries(query string, page int) (dto.TMDBPagedResponse[dto.TMDBMediaItem], error) {
	var result dto.TMDBPagedResponse[dto.TMDBMediaItem]
	body, err := s.request("/search/tv", url.Values{
		"query": {query},
		"page":  {strconv.Itoa(page)},
	})
	if err != nil {
		return result, err
	}
	return result, decode(body, &result)
}

func (s *TMDBService) SearchPeople(query string, page int) (dto.TMDBPagedResponse[dto.TMDBMediaItem], error) {
	var result dto.TMDBPagedResponse[dto.TMDBMediaItem]
	body, err := s.request("/search/person", url.Values{
		"query": {query},
		"page":  {strconv.Itoa(page)},
	})
	if err != nil {
		return result, err
	}
	return result, decode(body, &result)
}

// Lance une recherche combinee sur les films, les series et les personnes.
func (s *TMDBService) SearchAll(query string, page int) (dto.TMDBSearchResults, error) {
	results := dto.TMDBSearchResults{
		Query: query,
		Page:  page,
	}

	movies, moviesErr := s.SearchMovies(query, page)
	if moviesErr != nil {
		return results, moviesErr
	}
	results.Movies = movies.Results

	series, seriesErr := s.SearchSeries(query, page)
	if seriesErr != nil {
		return results, seriesErr
	}
	results.Series = series.Results

	people, peopleErr := s.SearchPeople(query, page)
	if peopleErr != nil {
		return results, peopleErr
	}
	results.People = people.Results

	results.HasMore = page < movies.TotalPages || page < series.TotalPages || page < people.TotalPages
	return results, nil
}

func (s *TMDBService) GetMovie(id int) (dto.TMDBMovieDetail, error) {
	var result dto.TMDBMovieDetail
	body, err := s.request(fmt.Sprintf("/movie/%d", id), nil)
	if err != nil {
		return result, err
	}
	return result, decode(body, &result)
}

func (s *TMDBService) GetSeries(id int) (dto.TMDBTVDetail, error) {
	var result dto.TMDBTVDetail
	body, err := s.request(fmt.Sprintf("/tv/%d", id), nil)
	if err != nil {
		return result, err
	}
	return result, decode(body, &result)
}

func (s *TMDBService) GetPerson(id int) (dto.TMDBPersonDetail, error) {
	var result dto.TMDBPersonDetail
	body, err := s.request(fmt.Sprintf("/person/%d", id), nil)
	if err != nil {
		return result, err
	}
	return result, decode(body, &result)
}

func (s *TMDBService) GetMovieCredits(id int) (dto.TMDBCredits, error) {
	var result dto.TMDBCredits
	body, err := s.request(fmt.Sprintf("/movie/%d/credits", id), nil)
	if err != nil {
		return result, err
	}
	return result, decode(body, &result)
}

func (s *TMDBService) GetSeriesCredits(id int) (dto.TMDBCredits, error) {
	var result dto.TMDBCredits
	body, err := s.request(fmt.Sprintf("/tv/%d/credits", id), nil)
	if err != nil {
		return result, err
	}
	return result, decode(body, &result)
}

// Ne garde que les realisateurs dans la liste de l'equipe technique.
func FilterDirectors(crew []dto.TMDBCrewMember) []dto.TMDBCrewMember {
	var directors []dto.TMDBCrewMember
	for _, member := range crew {
		if strings.EqualFold(member.Job, "Director") {
			directors = append(directors, member)
		}
	}
	return directors
}

// Limite la liste des acteurs aux premiers du casting.
func TopCast(cast []dto.TMDBCastMember, limit int) []dto.TMDBCastMember {
	if len(cast) <= limit {
		return cast
	}
	return cast[:limit]
}
