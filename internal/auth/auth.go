package Auth

import (
	"net/http"

	"golang.org/x/crypto/bcrypt"
	ChirpyDatabase "github.com/Couches/chirpy-database"
)

func HashPassword(password string) ChirpyDatabase.Result {
  hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 10)
  if err != nil {
    return ChirpyDatabase.GetErrorResult(http.StatusInternalServerError, err)
  }

  return ChirpyDatabase.GetOKResult(1, hashedPassword)
}
