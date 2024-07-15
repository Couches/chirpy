package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	ChirpyDatabase "github.com/Couches/chirpy-database"
	"github.com/golang-jwt/jwt/v5"
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
    fmt.Println(error.Msg, error.Code)
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

  user, error := config.UserDatabase.GetUser(userID)

  if error.Err != nil {
    fmt.Println(error.Msg, error.Code)
    respondWithError(w, error.Code, error.Msg)
    return
  }
 
  respondWithJSON(w, http.StatusOK, user)
}

func usersGetAllEndpoint(w http.ResponseWriter, request *http.Request, config apiConfig) {
  users, error := config.UserDatabase.GetUsers()

  if error.Err != nil {
    fmt.Println(error.Msg, error.Code)
    respondWithError(w, error.Code, error.Msg)
    return
  }

  respondWithJSON(w, http.StatusOK, users)
}

func usersUpdateEndpoint(w http.ResponseWriter, request *http.Request, config apiConfig) {
  header := request.Header.Get("Authorization")
  headers := strings.Fields(header)
  tokenString := headers[1]
  
	decoder := json.NewDecoder(request.Body)
	req := ChirpyDatabase.UserRequest{}
	err := decoder.Decode(&req)

  claims := jwt.RegisteredClaims{}
  token, err := jwt.ParseWithClaims(
    tokenString,
    &claims,
    func(token *jwt.Token) (interface{}, error) { return []byte(config.jwtSecret), nil},
    )

  if err != nil {
    respondWithError(w, http.StatusUnauthorized, err.Error())
    return
  }

  if !token.Valid {
    respondWithError(w, http.StatusUnauthorized, err.Error())
    return
  }

  userID, _ := strconv.Atoi(claims.Subject)

  userRequest := ChirpyDatabase.UpdateUserRequest {
    Id: userID,
    Email: req.Email,
    Password: req.Password,
  }

  error := config.UserDatabase.UpdateUser(userRequest)
  if error.Err != nil {
    fmt.Println(error.Err)
    respondWithError(w, error.Code, error.Msg)
    return
  }

  res := ChirpyDatabase.UserResponse {
    Id: userID,
    Email: req.Email,
  }

  respondWithJSON(w, http.StatusOK, res)
}
