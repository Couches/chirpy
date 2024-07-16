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


func endpointGetAllChirps(w http.ResponseWriter, _ *http.Request, config apiConfig) {
  result := config.Database.GetAllChirps()

	if result.Error != nil {
    respondWithError(w, result)
		return
	}

  respondWithJSON(w, result)
}

