package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	ChirpyDatabase "github.com/Couches/chirpy-database"
)

func chirpsCreateEndpoint(w http.ResponseWriter, request *http.Request, config apiConfig) {
	decoder := json.NewDecoder(request.Body)
	req := ChirpyDatabase.ChirpRequest{}
	err := decoder.Decode(&req)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Something went wrong while decoding request")
		return
	}

  chirp, error := config.ChirpDatabase.CreateChirp(req)
  if error.Err != nil {
    fmt.Print(error.Msg, error.Code)
    respondWithError(w, error.Code, error.Msg)
    return
  }

	respondWithJSON(w, http.StatusCreated, chirp)
}

func chirpsGetEndpoint(w http.ResponseWriter, request *http.Request, config apiConfig) {
  chirpID, err := strconv.Atoi(request.PathValue("chirpID"))
  if err != nil {
    respondWithError(w, http.StatusBadRequest, "Invalid input")
    return
  }

  chirp, error := config.ChirpDatabase.GetChirp(chirpID)

  if error.Err != nil {
    fmt.Println(error.Msg, error.Code)
    respondWithError(w, error.Code, error.Msg)
    return
  }

  respondWithJSON(w, http.StatusOK, chirp)
}

func chirpsGetAllEndpoint(w http.ResponseWriter, request *http.Request, config apiConfig) {
  chirps, error := config.ChirpDatabase.GetChirps()

  if error.Err != nil {
    fmt.Println(error.Msg, error.Code)
    respondWithError(w, error.Code, error.Msg)
    return
  }

  respondWithJSON(w, http.StatusOK, chirps)
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	http.Error(w, msg, code)
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	res, err := json.Marshal(payload)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}
  

	w.Header().Add("Content-Type", "text/json; charset=utf-8")
	w.WriteHeader(code)
	w.Write(res)
}

