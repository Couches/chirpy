package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type LoginRequest struct {
	Email      string `json:"email"`
	Password   string `json:"password"`
	ExpireTime int    `json:"expires_in_seconds,omitempty"`
}

func loginEndpoint(w http.ResponseWriter, request *http.Request, config apiConfig) {
	decoder := json.NewDecoder(request.Body)
	req := LoginRequest{}
	err := decoder.Decode(&req)

	fmt.Println(req)
  expiresIn := time.Duration(req.ExpireTime) * time.Second
  if expiresIn == 0 || expiresIn > 24 * time.Hour {
    expiresIn = 24 * time.Hour
  }

  expiresAt := time.Now().Add(expiresIn)

  claims := jwt.RegisteredClaims{
    Issuer: "chirpy",
    IssuedAt: jwt.NewNumericDate(time.Now()),
    ExpiresAt: jwt.NewNumericDate(expiresAt),
  }
  fmt.Println(claims)
  token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{})
  fmt.Println(token)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Something went wrong while decoding request")
		return
	}

	user, error := config.UserDatabase.GetUserByEmail(req.Email)

	if error.Err != nil {
		respondWithError(w, error.Code, error.Msg)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect user password")
		return
	}

	respondWithJSON(w, http.StatusOK, user.ToUserResponse())
}
