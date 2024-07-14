package ChirpDatabase

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"
)

type Chirp struct {
	Id    int    `json:"id"`
	Valid bool   `json:"valid"`
	Body  string `json:"body"`
}

type Database struct {
	path string
	mux  *sync.RWMutex
}

type DBStructure struct {
	Chirps map[int]Chirp `json:"chirps"`
}

type ChirpRequest struct {
	Body string `json:"body"`
}

type CustomError struct {
	Err  error
	Code int
	Msg  string
}

func NewDatabase(path string) (*Database, error) {
	db := &Database{}
	file, err := os.Create(path)
	if err != nil {
		fmt.Println("Failed to create file \"%v\": %v", path, err)
		return db, err
	}

	defer file.Close()

	db = &Database{
		path: path,
		mux:  &sync.RWMutex{},
	}

	return db, nil
}

func (db *Database) CreateChirp(req ChirpRequest) (Chirp, CustomError) {
	chirp := Chirp{}

	if len(req.Body) > 140 {
		msg := "Chirp is too long"
		error := CustomError{
			Err:  errors.New(msg),
			Code: http.StatusBadRequest,
			Msg:  msg,
		}
		return chirp, error
	}

	db.mux.Lock()
	defer db.mux.Unlock()

	dbs, err := db.loadDB()
	if err != nil {
		fmt.Printf("Failed to load database: %v\n", err)
	}

	chirp = Chirp{
		Body: cleanMessage(req.Body),
		Id:          len(dbs.Chirps) + 1,
		Valid:       true,
	}

	dbs.Chirps[chirp.Id] = chirp

	err = db.writeDB(dbs)
	if err != nil {
		msg := "Failed to write Chirp to database"
		error := CustomError{
			Err:  errors.New(msg),
			Code: http.StatusInternalServerError,
			Msg:  msg,
		}
		return chirp, error
	}

	return chirp, CustomError{}
}

func (db *Database) GetChirps() ([]Chirp, CustomError) {
	dbs, err := db.loadDB()

	if err != nil {
		msg := "Failed to load database"
		fmt.Printf("%v\n", msg)

		return nil, getError(msg, http.StatusInternalServerError, err)
	}

	chirps := make([]Chirp, 0, len(dbs.Chirps))
	for _, chirp := range dbs.Chirps {
		chirps = append(chirps, chirp)
	}

	return chirps, CustomError{}
}

func (db *Database) GetChirp(chirpID int) (Chirp, CustomError) {
  dbs, err := db.loadDB()

	if err != nil {
		msg := "Failed to load database"
		fmt.Printf("%v\n", msg)

		return Chirp{}, getError(msg, http.StatusInternalServerError, err)
	}

  chirp, ok := dbs.Chirps[chirpID]
  if !ok {
    msg := fmt.Sprintf("Chirp with chirpID \"%v\" not found in database", chirpID)
    return Chirp{}, getError(msg, http.StatusNotFound, errors.New(msg))
  }

  return chirp, CustomError{}
}

func (db *Database) loadDB() (DBStructure, error) {
	dbs := DBStructure{
		Chirps: map[int]Chirp{},
	}
	file, err := os.ReadFile(db.path)
	if err != nil {
		fmt.Printf("Failed to read file \"%v\"\n", db.path)
		return dbs, err
	}

	fmt.Println(string(file))
	err = json.Unmarshal(file, &dbs)

	if err != nil {
		fmt.Printf("Failed to deserialize JSON from file \"%v\"\n", db.path)
		return dbs, err
	}

	return dbs, nil
}

func (db *Database) writeDB(dbs DBStructure) error {
	file, err := os.Create(db.path)
	defer file.Close()

	if err != nil {
		fmt.Printf("Failed to create or truncate file \"%v\"\n", db.path)
		return err
	}

	body, err := json.Marshal(dbs)
	if err != nil {
		fmt.Printf("Failed to serialize Chirps to JSON for writing\n")
		return err
	}

	written, err := file.Write(body)
	if err != nil {
		fmt.Printf("Failed to write to file \"%v\"\n", db.path)
		return err
	}

	fmt.Printf("Successfully wrote %v bytes to file \"%v\"\n", written, db.path)
	return nil
}

func cleanMessage(body string) string {
	filteredWords := map[string]bool{
		"kerfuffle": true,
		"sharbert":  true,
		"fornax":    true,
	}

	arr := strings.Fields(body)
	for i, word := range arr {
		if _, ok := filteredWords[strings.ToLower(word)]; ok {
			arr[i] = "****"
		}
	}

	return strings.Join(arr, " ")
}

func getError(msg string, code int, err error) CustomError {
	return CustomError{
		Msg:  msg,
		Code: code,
		Err:  err,
	}
}
