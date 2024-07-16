package main

import (
	"net/http"
	"strings"

	Auth "github.com/Couches/auth"
	ChirpyDatabase "github.com/Couches/chirpy-database"
)

func loginEndpoint(w http.ResponseWriter, r *http.Request, config apiConfig) {
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

	result = config.Database.GetUserByEmail(req.Email)

	if result.Error != nil {
		respondWithError(w, result)
		return
	}

	user := (*result.Body).(ChirpyDatabase.User)

	result = Auth.ComparePassword(req.Password, user.HashedPassword)
	if result.Error != nil {
		respondWithError(w, result)
		return
	}

	result = Auth.CreateJWT(user.Id, config.jwtSecret, config.jwtExpireTime)
	if result.Error != nil {
		respondWithError(w, result)
		return
	}

	token := (*result.Body).(string)

	result = config.Database.CreateRefreshToken(user.Id)
	if result.Error != nil {
		respondWithError(w, result)
		return
	}

	refreshToken := (*result.Body).(ChirpyDatabase.RefreshToken)

	loginResponse := struct {
		Id           int    `json:"id"`
		Email        string `json:"email"`
		Token        string `json:"token"`
		RefreshToken string `json:"refresh_token"`
		IsChirpyRed  bool   `json:"is_chirpy_red"`
	}{
		Id:           user.Id,
		Email:        user.Email,
		Token:        token,
		RefreshToken: refreshToken.Value,
		IsChirpyRed:  user.IsChirpyRed,
	}

	result = ChirpyDatabase.GetOKResult(http.StatusOK, loginResponse)
	respondWithJSON(w, result)
}

func refreshEndpoint(w http.ResponseWriter, r *http.Request, config apiConfig) {
	requestToken := r.Header.Get("Authorization")
	splitToken := strings.Fields(requestToken)
	requestToken = splitToken[1]

	result := config.Database.GetRefreshToken(requestToken)
	if result.Error != nil {
		respondWithError(w, result)
		return
	}

	refreshToken := (*result.Body).(ChirpyDatabase.RefreshToken)

	result = Auth.CreateJWT(refreshToken.UserID, config.jwtSecret, config.jwtExpireTime)
	if result.Error != nil {
		respondWithError(w, result)
		return
	}

	token := (*result.Body).(string)

	refreshResponse := struct {
		Token string `json:"token"`
	}{
		Token: token,
	}

	result = ChirpyDatabase.GetOKResult(http.StatusOK, refreshResponse)
	respondWithJSON(w, result)
}

func revokeEndpoint(w http.ResponseWriter, r *http.Request, config apiConfig) {
	requestToken := r.Header.Get("Authorization")
	splitToken := strings.Fields(requestToken)
	requestToken = splitToken[1]

	result := config.Database.DeleteRefreshToken(requestToken)
	if result.Error != nil {
		respondWithError(w, result)
		return
	}

	respondWithJSON(w, result)
}
