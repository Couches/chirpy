package main

import (
	"net/http"
	"strconv"

	ChirpyDatabase "github.com/Couches/chirpy-database"
)


func endpointCreateChirp(w http.ResponseWriter, r *http.Request, config apiConfig) {
  type parameters struct {
    Body string `json:"body"`
  }

  result := decodeRequestBody(r, &parameters{})
  if result.Error != nil {
    respondWithError(w, result)
    return
  }

  req := (*result.Body).(*parameters)

  result = config.Database.CreateChirp(req.Body)

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
