package Auth

import (
	"errors"
	"math"
	"net/http"
	"strconv"
	"time"

	ChirpyDatabase "github.com/Couches/chirpy-database"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) ChirpyDatabase.Result {
  hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 10)
  if err != nil {
    return ChirpyDatabase.GetErrorResult(http.StatusInternalServerError, err)
  }

  return ChirpyDatabase.GetOKResult(1, string(hashedPassword))
}


func ComparePassword(password, hashed_password string) ChirpyDatabase.Result {
  err := bcrypt.CompareHashAndPassword([]byte(hashed_password), []byte(password))

  if err != nil {
    return ChirpyDatabase.GetErrorResult(http.StatusUnauthorized, err)
  }

  return ChirpyDatabase.GetOKResult(http.StatusOK, nil)
}


func ValidateJWT(tokenString, tokenSecret string) ChirpyDatabase.Result {
  claims := jwt.RegisteredClaims {}

  token, err := jwt.ParseWithClaims(
    tokenString,
    &claims,
    func (token *jwt.Token) (interface{}, error) { return []byte(tokenSecret), nil },
    )

  if err != nil {
    return ChirpyDatabase.GetErrorResult(http.StatusUnauthorized, err)
  }

  if !token.Valid {
    err := errors.New("Invalid token")
    return ChirpyDatabase.GetErrorResult(http.StatusUnauthorized, err)
  }

  issuer, err := token.Claims.GetIssuer()
  if err != nil {
    return ChirpyDatabase.GetErrorResult(http.StatusUnauthorized, err)
  }

  if issuer != "chirpy" {
    err := errors.New("Invalid issuer")
    return ChirpyDatabase.GetErrorResult(http.StatusUnauthorized, err)
  }

  subject, err := token.Claims.GetSubject()
  if err != nil {
    return ChirpyDatabase.GetErrorResult(http.StatusInternalServerError, err)
  }

  return ChirpyDatabase.GetOKResult(http.StatusOK, subject)
}


func CreateJWT(userID int, tokenSecret string, expiresInSeconds int) ChirpyDatabase.Result {
  if expiresInSeconds <= 0 {
    expiresInSeconds = 86400
  }
  expiresInSeconds = int(math.Min(float64(expiresInSeconds), 86400))

  expiresIn := time.Duration(expiresInSeconds * int(time.Second))

  claims := jwt.RegisteredClaims {
    Issuer: "chirpy",
    IssuedAt: jwt.NewNumericDate(time.Now().UTC()),
    ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(expiresIn)),
    Subject: strconv.Itoa(userID),
  }

  token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
  signed, err := token.SignedString([]byte(tokenSecret))
  if err != nil {
    return ChirpyDatabase.GetErrorResult(http.StatusInternalServerError, err)
  }

  return ChirpyDatabase.GetOKResult(1, signed)
}
