package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type LoginRequest struct {
	Email      string `json:"email"`
	Password   string `json:"password"`
	ExpireTime int    `json:"expires_in_seconds,omitempty"`
}

type LoginResponse struct {
	Id    int    `json:"id"`
	Email string `json:"email"`
	Token string `json:"token"`
}

func loginEndpoint(w http.ResponseWriter, request *http.Request, config apiConfig) {
	decoder := json.NewDecoder(request.Body)
	req := LoginRequest{}
	err := decoder.Decode(&req)

	expiresIn := time.Duration(req.ExpireTime) * time.Second
	if expiresIn == 0 || expiresIn > 24*time.Hour {
		expiresIn = 24 * time.Hour
	}

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Something went wrong while decoding request")
		return
	}

	user, error := config.UserDatabase.GetUserByEmail(req.Email)

	if error.Err != nil {
		respondWithError(w, error.Code, error.Msg)
		return
	}

	expiresAt := time.Now().Add(expiresIn)

	claims := jwt.RegisteredClaims{
		Issuer:    "chirpy",
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(expiresAt),
		Subject:   strconv.Itoa(user.Id),
	}

	new_token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
  token, err := new_token.SignedString([]byte(config.jwtSecret))
  if err != nil {
    respondWithError(w, http.StatusInternalServerError, err.Error())
    return
  }

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	loginResponse := LoginResponse{
		Id: user.Id,
    Email: user.Email,
    Token: token,
	}

	respondWithJSON(w, http.StatusOK, loginResponse)
}

