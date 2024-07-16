package ChirpyDatabase

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"time"
)

type RefreshToken struct {
  Value string `json:"value"`
  ExpirationTime time.Time `json:"expiration"`
  UserID int
}


func (db *Database) CreateRefreshToken(userID int) Result {
  var bytes [32]byte
  _, err := rand.Read(bytes[:])

  if err != nil {
    return GetErrorResult(http.StatusInternalServerError, err)
  }

  tokenString := hex.EncodeToString(bytes[:])
  expiresIn := 60 * 24 * time.Hour
  expirationTime := time.Now().Add(expiresIn)

  token := RefreshToken {
    Value: tokenString,
    ExpirationTime: expirationTime,
    UserID: userID,
  }

  result := db.LoadDB()
  if result.Error != nil {
    return result
  }

  structure := (*result.Body).(DatabaseStructure)

  structure.Tokens[tokenString] = token

  result = db.WriteDB(structure)
  if result.Error != nil {
    return result
  }

  return GetOKResult(1, token)
}


func (db *Database) GetRefreshToken(tokenString string) Result {
  result := db.LoadDB()
  if result.Error != nil {
    return result
  }

  structure := (*result.Body).(DatabaseStructure)

  if token, ok := structure.Tokens[tokenString]; ok {
    return GetOKResult(1, token)
  }

  msg := fmt.Sprintf("No refresh token found with value \"%v\"", tokenString)
  return GetErrorResult(http.StatusUnauthorized, errors.New(msg))
}

func (db *Database) DeleteRefreshToken(tokenString string) Result {
  result := db.LoadDB()
  if result.Error != nil {
    return result
  }

  structure := (*result.Body).(DatabaseStructure)

  if _, ok := structure.Tokens[tokenString]; !ok {
    msg := fmt.Sprintf("No refresh token found with value \"%v\"", tokenString)
    return GetErrorResult(http.StatusUnauthorized, errors.New(msg))
  }
  
  delete(structure.Tokens, tokenString)

  result = db.WriteDB(structure)
  if result.Error != nil {
    return result
  }

  return GetOKResult(http.StatusNoContent, nil)
}
