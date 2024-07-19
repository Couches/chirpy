package main

import (
	"net/http"
	"strconv"
	"strings"

	Auth "github.com/Couches/auth"
	ChirpyDatabase "github.com/Couches/chirpy-database"
)

func endpointCreateChirp(w http.ResponseWriter, r *http.Request, config apiConfig) {
	requestToken := r.Header.Get("Authorization")
	splitToken := strings.Fields(requestToken)
	requestToken = splitToken[1]

	result := Auth.ValidateJWT(requestToken, config.jwtSecret)
	if result.Error != nil {
		respondWithError(w, result)
		return
	}

	subject := (*result.Body).(string)
	userID, _ := strconv.Atoi(subject)

	type parameters struct {
		Body string `json:"body"`
	}

	result = decodeRequestBody(r, &parameters{})
	if result.Error != nil {
		respondWithError(w, result)
		return
	}

	req := (*result.Body).(*parameters)

	result = config.Database.CreateChirp(req.Body, userID)

	if result.Error != nil {
		respondWithError(w, result)
		return
	}

	respondWithJSON(w, result)
}

func endpointDeleteChirp(w http.ResponseWriter, r *http.Request, config apiConfig) {
	requestToken := r.Header.Get("Authorization")
	splitToken := strings.Fields(requestToken)
	requestToken = splitToken[1]

	result := Auth.ValidateJWT(requestToken, config.jwtSecret)
	if result.Error != nil {
		respondWithError(w, result)
		return
	}

	chirpID, err := strconv.Atoi(r.PathValue("chirpID"))
	if err != nil {
		error := ChirpyDatabase.GetErrorResult(http.StatusInternalServerError, err)
		respondWithError(w, error)
		return
	}

	subject := (*result.Body).(string)
	userID, _ := strconv.Atoi(subject)

	result = config.Database.DeleteChirp(chirpID, userID)
	if result.Error != nil {
		respondWithError(w, result)
		return
	}

	respondWithJSON(w, result)
}

func endpointGetChirp(w http.ResponseWriter, r *http.Request, config apiConfig) {
	chirpID, err := strconv.Atoi(r.PathValue("chirpID"))
	if err != nil {
		error := ChirpyDatabase.GetErrorResult(http.StatusInternalServerError, err)
		respondWithError(w, error)
		return
	}

	result := config.Database.GetChirp(chirpID)

	if result.Error != nil {
		respondWithError(w, result)
		return
	}

	respondWithJSON(w, result)
}

func endpointGetChirps(w http.ResponseWriter, r *http.Request, config apiConfig) {
	authorIDString := r.URL.Query().Get("author_id")
	authorID := 0
	if len(authorIDString) > 0 {
    parsed, err := strconv.Atoi(authorIDString)

		if err != nil {
			error := ChirpyDatabase.GetErrorResult(http.StatusInternalServerError, err)
			respondWithError(w, error)
			return
		}

    authorID = parsed
	}

	sortBy := r.URL.Query().Get("sort")

	result := config.Database.GetChirpsFromAuthor(authorID, sortBy)

	if result.Error != nil {
		respondWithError(w, result)
		return
	}

	respondWithJSON(w, result)
}
