package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	ChirpyDatabase "github.com/Couches/chirpy-database"
)

func usersCreateEndpoint(w http.ResponseWriter, request *http.Request, config apiConfig) {
	decoder := json.NewDecoder(request.Body)
	req := ChirpyDatabase.UserRequest{}
	err := decoder.Decode(&req)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Something went wrong while decoding request")
		return
	}

  chirp, error := config.UserDatabase.CreateUser(req)
  if error.Err != nil {
    fmt.Print(error.Msg, error.Code)
    respondWithError(w, error.Code, error.Msg)
    return
  }

	respondWithJSON(w, http.StatusCreated, chirp)
}

func usersGetEndpoint(w http.ResponseWriter, request *http.Request, config apiConfig) {
  userID, err := strconv.Atoi(request.PathValue("userID"))
  if err != nil {
    respondWithError(w, http.StatusBadRequest, "Invalid input")
    return
  }

  chirp, error := config.UserDatabase.GetUser(userID)

  if error.Err != nil {
    fmt.Println(error.Msg, error.Code)
    respondWithError(w, error.Code, error.Msg)
    return
  }

  respondWithJSON(w, http.StatusOK, chirp)
}

func usersGetAllEndpoint(w http.ResponseWriter, request *http.Request, config apiConfig) {
  chirps, error := config.UserDatabase.GetUsers()

  if error.Err != nil {
    fmt.Println(error.Msg, error.Code)
    respondWithError(w, error.Code, error.Msg)
    return
  }

  respondWithJSON(w, http.StatusOK, chirps)
}
