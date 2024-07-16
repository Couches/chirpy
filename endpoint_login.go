package main

import (
	"net/http"

	Auth "github.com/Couches/auth"
	ChirpyDatabase "github.com/Couches/chirpy-database"
)

func loginEndpoint(w http.ResponseWriter, r *http.Request, config apiConfig) {
	type parameters struct {
		Email      string `json:"email"`
		Password   string `json:"password"`
		ExpireTime int    `json:"expires_in_seconds,omitempty"`
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

	result = Auth.CreateJWT(user.Id, config.jwtSecret, req.ExpireTime)
	if result.Error != nil {
		respondWithError(w, result)
		return
	}

	token := (*result.Body).(string)

	loginResponse := struct {
		Id    int    `json:"id"`
		Email string `json:"email"`
		Token string `json:"token"`
	}{
		Id:    user.Id,
		Email: user.Email,
		Token: token,
	}

	result = ChirpyDatabase.GetOKResult(http.StatusOK, loginResponse)
	respondWithJSON(w, result)
}
