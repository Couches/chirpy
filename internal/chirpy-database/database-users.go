package ChirpyDatabase

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
)

type User struct {
	Id    int    `json:"id"`
	Valid bool   `json:"valid"`
	Email string `json:"email"`
}

type UserRequest struct {
	Email string `json:"email"`
}

func (u User) GetID() int {
	return u.Id
}

func (u User) IsValid() bool {
	return u.Valid
}

func (db *Database) CreateUser(req UserRequest) (User, CustomError) {
	user := User{}

	db.mux.Lock()
	defer db.mux.Unlock()

	dbs, err := db.loadUserDB()
	if err != nil {
		fmt.Printf("Failed to load database: %v\n", err)
	}

	user = User{
		Id:    len(dbs.Data) + 1,
    Email: req.Email,
		Valid: true,
	}

	dbs.Data[user.Id] = user

	err = db.writeDB(dbs)
	if err != nil {
		msg := "Failed to write Chirp to database"
		error := CustomError{
			Err:  errors.New(msg),
			Code: http.StatusInternalServerError,
			Msg:  msg,
		}
		return user, error
	}

	return user, CustomError{}
}

func (db *Database) GetUsers() ([]User, CustomError) {
	dbs, err := db.loadUserDB()

	if err != nil {
		msg := "Failed to load database"
		fmt.Printf("%v\n", msg)

		return nil, getError(msg, http.StatusInternalServerError, err)
	}

	users := make([]User, 0, len(dbs.Data))
	for _, user := range dbs.Data {
    users = append(users, user)
	}

	return users, CustomError{}
}

func (db *Database) GetUser(userID int) (User, CustomError) {
	dbs, err := db.loadUserDB()

	if err != nil {
		msg := "Failed to load database"
		fmt.Printf("%v\n", msg)

		return User{}, getError(msg, http.StatusInternalServerError, err)
	}

	user, ok := dbs.Data[userID]
	if !ok {
		msg := fmt.Sprintf("User with userID \"%v\" not found in database", userID)
		return User{}, getError(msg, http.StatusNotFound, errors.New(msg))
	}

	return user, CustomError{}
}

func (db *Database) loadUserDB() (DBStructure[User], error) {
	dbs := DBStructure[User]{
		Data: map[int]User{},
	}
	file, err := os.ReadFile(db.path)
	if err != nil {
		fmt.Printf("Failed to read file \"%v\"\n", db.path)
		return dbs, err
	}

	fmt.Println(string(file))
	err = json.Unmarshal(file, &dbs)

	if err != nil {
		fmt.Printf("Failed to deserialize JSON from file \"%v\" with err %v\n", db.path, err)
		return dbs, err
	}

	return dbs, nil
}
