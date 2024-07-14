package ChirpyDatabase

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
  "golang.org/x/crypto/bcrypt"
)

type User struct {
	Id       int    `json:"id"`
	Valid    bool   `json:"valid"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserResponse struct {
	Id       int    `json:"id"`
	Valid    bool   `json:"valid"`
	Email    string `json:"email"`
}

func (u User) GetID() int {
	return u.Id
}

func (u User) IsValid() bool {
	return u.Valid
}

func (u User) ToUserResponse() UserResponse {
  return UserResponse{
    Id: u.Id,
    Valid: u.Valid,
    Email: u.Email,
  }
}

func (db *Database) CreateUser(req UserRequest) (UserResponse, CustomError) {
	user := User{}

	db.mux.Lock()
	defer db.mux.Unlock()

	dbs, err := db.loadUserDB()
	if err != nil {
		fmt.Printf("Failed to load database: %v\n", err)
	}

  for _, user := range(dbs.Data) {
    if req.Email == user.Email {
      msg := "Email already in use"
      error := CustomError {
        Err: errors.New(msg),
        Code: http.StatusConflict,
        Msg: msg,
      }
      return UserResponse{}, error
    }
  }

  pass, err := bcrypt.GenerateFromPassword([]byte(req.Password), 10)
  if err != nil {
    msg := "Could not hash user password"
    error := CustomError {
      Err: errors.New(msg),
      Code: http.StatusInternalServerError,
      Msg: msg,
    }
    return UserResponse{}, error 
  } 

	user = User{
		Id:    len(dbs.Data) + 1,
		Email: req.Email,
		Valid: true,
    Password: string(pass),
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
		return user.ToUserResponse(), error
	}

	return user.ToUserResponse(), CustomError{}
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

func (db *Database) GetUser(userID int) (UserResponse, CustomError) {
	dbs, err := db.loadUserDB()

	if err != nil {
		msg := "Failed to load database"
		fmt.Printf("%v\n", msg)

		return UserResponse{}, getError(msg, http.StatusInternalServerError, err)
	}

	user, ok := dbs.Data[userID]
	if !ok {
		msg := fmt.Sprintf("User with userID \"%v\" not found in database", userID)
		return UserResponse{}, getError(msg, http.StatusNotFound, errors.New(msg))
	}

	return user.ToUserResponse(), CustomError{}
}

func (db *Database) GetUserByEmail(email string) (User, CustomError) {
  users, error := db.GetUsers()
  if error.Err != nil {
    return User{}, error
  }

  for _, user := range users {
    if user.Email == email {
      return user, CustomError{}
    }
  }

  msg := fmt.Sprintf("User not found with email \"%v\"", email)
  err := errors.New(msg)
  return User{}, getError(msg, http.StatusNotFound, err)
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
