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

func RadarrLookupName(term string) ([]types.MovieData, error) {
	values := make(map[string]string)
	values["apikey"] = config.Config.RadarrToken
	values["term"] = term
	u := helpers.NewUrl(config.Config.RadarrHost, "api/movie/lookup", values)
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

func RadarrLookupID(ID string) (types.MovieData, error) {
	values := make(map[string]string)
	values["apikey"] = config.Config.RadarrToken
	values["tmdbId"] = ID
	u := helpers.NewUrl(config.Config.RadarrHost, "api/movie/lookup/tmdb", values)
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

func RadarrAdd(movie types.MovieData) (int, error) {
	postValues := make(map[string]string)
	postValues["apikey"] = config.Config.RadarrToken
	postUrl := helpers.NewUrl(config.Config.RadarrHost, "api/movie", postValues)
	body, _ := json.Marshal(movie)
	resp, err := http.Post(postUrl.String(), "application/json", bytes.NewBuffer((body)))
	// Check if the request went through
	if err != nil {
		return -1, err
	}
	// Check if the movie was successfully added
	code := resp.StatusCode
	return code, nil
}
