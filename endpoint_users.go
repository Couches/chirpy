package main

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	Auth "github.com/Couches/auth"
	ChirpyDatabase "github.com/Couches/chirpy-database"
)

func endpointCreateUser(w http.ResponseWriter, r *http.Request, config apiConfig) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	result := decodeRequestBody(r, &parameters{})
	if result.Error != nil {
		respondWithError(w, result)
		return
	}

	req := (*result.Body).(*parameters)

	result = Auth.HashPassword(req.Password)
	if result.Error != nil {
		respondWithError(w, result)
		return
	}

	hashedPassword := (*result.Body).(string)

	result = config.Database.CreateUser(req.Email, hashedPassword)

	if result.Error != nil {
		respondWithError(w, result)
		return
	}

	user := (*result.Body).(ChirpyDatabase.User)

	createResponse := struct {
		Id          int    `json:"id"`
		Email       string `json:"email"`
		IsChirpyRed bool   `json:"is_chirpy_red"`
	}{
		Id:          user.Id,
		Email:       user.Email,
		IsChirpyRed: user.IsChirpyRed,
	}

	result = ChirpyDatabase.GetOKResult(http.StatusCreated, createResponse)
	respondWithJSON(w, result)
}

func endpointUpdateUserLogin(w http.ResponseWriter, r *http.Request, config apiConfig) {
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
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	result = decodeRequestBody(r, &parameters{})
	if result.Error != nil {
		respondWithError(w, result)
		return
	}

	req := (*result.Body).(*parameters)

	result = Auth.HashPassword(req.Password)
	if result.Error != nil {
		respondWithError(w, result)
		return
	}

	hashedPassword := (*result.Body).(string)

	result = config.Database.UpdateUserLogin(userID, req.Email, hashedPassword)
	if result.Error != nil {
		respondWithError(w, result)
		return
	}

	user := (*result.Body).(ChirpyDatabase.User)

	updateResponse := struct {
		Id          int    `json:"id"`
		Email       string `json:"email"`
		IsChirpyRed bool   `json:"is_chirpy_red"`
	}{
		Id:          user.Id,
		Email:       user.Email,
		IsChirpyRed: user.IsChirpyRed,
	}

	result = ChirpyDatabase.GetOKResult(http.StatusOK, updateResponse)
	respondWithJSON(w, result)
}

func endpointGetUser(w http.ResponseWriter, r *http.Request, config apiConfig) {
	userID, err := strconv.Atoi(r.PathValue("userID"))
	if err != nil {
		error := ChirpyDatabase.GetErrorResult(http.StatusInternalServerError, err)
		respondWithError(w, error)
		return
	}

	result := config.Database.GetUser(userID)

	if result.Error != nil {
		respondWithError(w, result)
		return
	}

	respondWithJSON(w, result)
}

func endpointGetAllUsers(w http.ResponseWriter, r *http.Request, config apiConfig) {
	result := config.Database.GetAllUsers()

	if result.Error != nil {
		respondWithError(w, result)
		return
	}

	respondWithJSON(w, result)
}

func upgradeUserEndpoint(w http.ResponseWriter, r *http.Request, config apiConfig) {
	apiKey := r.Header.Get("Authorization")
	splitToken := strings.Fields(apiKey)

  if len(splitToken) < 2 {
    respondWithError(w, ChirpyDatabase.GetErrorResult(http.StatusUnauthorized, errors.New("Unauthorized")))
    return
  }

	apiKey = splitToken[1]

  if apiKey != config.apiKey {
    respondWithError(w, ChirpyDatabase.GetErrorResult(http.StatusUnauthorized, errors.New("Unauthorized")))
    return
  }

	type parameters struct {
		Event string `json:"event"`
		Data  struct {
			UserID int `json:"user_id"`
		} `json:"data"`
	}

  result := decodeRequestBody(r, &parameters{})
	if result.Error != nil {
		respondWithError(w, result)
		return
	}

	req := (*result.Body).(*parameters)

  if req.Event != "user.upgraded" {
    respondWithJSON(w, ChirpyDatabase.GetOKResult(http.StatusNoContent, nil))
    return
  }

  result = config.Database.UpgradeUser(req.Data.UserID)
  if result.Error != nil {
    respondWithError(w, result)
    return
  }

  respondWithJSON(w, result)
}
