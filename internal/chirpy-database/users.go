package ChirpyDatabase

import (
	"errors"
	"fmt"
	"net/http"
)

type User struct {
	Id             int    `json:"id"`
	Email          string `json:"email"`
	HashedPassword string `json:"hashed_password"`
	IsChirpyRed    bool   `json:"is_chirpy_red"`
	RefreshToken   RefreshToken
}

func (db *Database) CreateUser(email, hashedPassword string) Result {
	result := db.GetUserByEmail(email)
	if result.Body != nil {
		msg := fmt.Sprintf("Email \"%v\" already in use", email)
		return GetErrorResult(http.StatusConflict, errors.New(msg))
	}

	result = db.LoadDB()
	if result.Error != nil {
		return result
	}

	structure := (*result.Body).(DatabaseStructure)

	user := User{
		Id:             *db.nextUserID,
		Email:          email,
		HashedPassword: hashedPassword,
	}

	structure.Users[user.Id] = user

	result = db.WriteDB(structure)
	if result.Error != nil {
		return result
	}

	*db.nextUserID++

	return GetOKResult(http.StatusCreated, user)
}

func (db *Database) UpdateUserLogin(userID int, email, hashedPassword string) Result {
	result := db.GetUser(userID)
	if result.Error != nil {
		return result
	}

	user := (*result.Body).(User)
	user.HashedPassword = hashedPassword
	user.Email = email

	result = db.LoadDB()
	if result.Error != nil {
		return result
	}

	structure := (*result.Body).(DatabaseStructure)
	structure.Users[user.Id] = user

	result = db.WriteDB(structure)
	if result.Error != nil {
		return result
	}

	return GetOKResult(http.StatusCreated, user)
}

func (db *Database) GetUser(id int) Result {
	result := db.GetAllUsers()
	if result.Error != nil {
		return result
	}

	users := (*result.Body).(map[int]User)
	if user, ok := users[id]; ok {
		return GetOKResult(http.StatusOK, user)
	}

	msg := fmt.Sprintf("No user was found with id \"%v\"", id)
	return GetErrorResult(http.StatusNotFound, errors.New(msg))
}

func (db *Database) GetUserByEmail(email string) Result {
	result := db.GetAllUsers()
	if result.Error != nil {
		return result
	}

	users := (*result.Body).(map[int]User)
	for _, user := range users {
		if user.Email == email {
			return GetOKResult(http.StatusOK, user)
		}
	}

	msg := fmt.Sprintf("No user was found with email \"%v\"", email)
	return GetErrorResult(http.StatusNotFound, errors.New(msg))
}

func (db *Database) GetAllUsers() Result {
	result := db.LoadDB()
	if result.Error != nil {
		return result
	}

	structure := (*result.Body).(DatabaseStructure)
	if len(structure.Users) == 0 {
		msg := "No users found"
		return GetErrorResult(http.StatusNoContent, errors.New(msg))
	}

	return GetOKResult(http.StatusOK, structure.Users)
}

func (db *Database) UpgradeUser(userID int) Result {
	result := db.GetUser(userID)
	if result.Error != nil {
		return result
	}

	user := (*result.Body).(User)
  user.IsChirpyRed = true

	result = db.LoadDB()
	if result.Error != nil {
		return result
	}

	structure := (*result.Body).(DatabaseStructure)
	structure.Users[user.Id] = user

	result = db.WriteDB(structure)
	if result.Error != nil {
		return result
	}

	return GetOKResult(http.StatusNoContent, nil)
}
