package api

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/AdrienCos/pidarr_bot/internal/config"
	"github.com/AdrienCos/pidarr_bot/internal/helpers"
	"github.com/AdrienCos/pidarr_bot/internal/types"
)

// RadarrLookupName searches online for all movies matching the given term
func RadarrLookupName(term string) ([]types.MovieData, error) {
	values := make(map[string]string)
	values["apikey"] = config.Config.RadarrToken
	values["term"] = term
	u := helpers.NewURL(config.Config.RadarrHost, "api/movie/lookup", values)
	response, err := http.Get(u.String())
	if err != nil {
		return nil, err
	}
	data, _ := ioutil.ReadAll(response.Body)
	// For all movies returned, extract the ID, name and year
	movies := []types.MovieData{}
	json.Unmarshal(data, &movies)
	return movies, nil
}

// RadarrLookupID searches online for the movie matching the given TMDB ID
func RadarrLookupID(ID string) (types.MovieData, error) {
	values := make(map[string]string)
	values["apikey"] = config.Config.RadarrToken
	values["tmdbId"] = ID
	u := helpers.NewURL(config.Config.RadarrHost, "api/movie/lookup/tmdb", values)
	response, err := http.Get(u.String())
	if err != nil {
		return types.MovieData{}, err
	}
	// Extract the movie data
	data, _ := ioutil.ReadAll(response.Body)
	movie := types.MovieData{}
	json.Unmarshal(data, &movie)
	return movie, nil
}

// RadarrAdd adds the given movie to the Radarr library
func RadarrAdd(movie types.MovieData) (int, error) {
	postValues := make(map[string]string)
	postValues["apikey"] = config.Config.RadarrToken
	postURL := helpers.NewURL(config.Config.RadarrHost, "api/movie", postValues)
	body, _ := json.Marshal(movie)
	resp, err := http.Post(postURL.String(), "application/json", bytes.NewBuffer((body)))
	// Check if the request went through
	if err != nil {
		return -1, err
	}
	// Check if the movie was successfully added
	code := resp.StatusCode
	return code, nil
}

// RadarrGetAll returns a list of all movies in the Radarr collection
func RadarrGetAll() ([]types.MovieData, error) {
	URL := helpers.NewURL(config.Config.RadarrHost, "api/movie", nil)
	resp, err := http.Get(URL.String())
	if err != nil {
		return nil, err
	}
	data, _ := ioutil.ReadAll(resp.Body)
	movies := []types.MovieData{}
	json.Unmarshal(data, &movies)
	return movies, nil
}
