package ChirpyDatabase

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
)

type ChirpDBStructure struct {
	Data map[int]Chirp `json:"data"`
}

type Chirp struct {
	Id    int    `json:"id"`
	Valid bool   `json:"valid"`
	Body  string `json:"body"`
}

type ChirpRequest struct {
	Body string `json:"body"`
}

func (c Chirp) GetID() int {
	return c.Id
}

func (c Chirp) IsValid() bool {
	return c.Valid
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

	dbs, err := db.loadChirpDB()
	if err != nil {
		fmt.Printf("Failed to load database: %v\n", err)
	}

	chirp = Chirp{
		Body:  cleanMessage(req.Body),
		Id:    len(dbs.Data) + 1,
		Valid: true,
	}

	dbs.Data[chirp.Id] = chirp

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
	dbs, err := db.loadChirpDB()

	if err != nil {
		msg := "Failed to load database"
		fmt.Printf("%v\n", msg)

		return nil, getError(msg, http.StatusInternalServerError, err)
	}

	chirps := make([]Chirp, 0, len(dbs.Data))
	for _, chirp := range dbs.Data {
    chirps = append(chirps, chirp)
	}

	return chirps, CustomError{}
}

func (db *Database) GetChirp(chirpID int) (Chirp, CustomError) {
	dbs, err := db.loadChirpDB()

	if err != nil {
		msg := "Failed to load database"
		fmt.Printf("%v\n", msg)

		return Chirp{}, getError(msg, http.StatusInternalServerError, err)
	}

	chirp, ok := dbs.Data[chirpID]
	if !ok {
		msg := fmt.Sprintf("Chirp with chirpID \"%v\" not found in database", chirpID)
		return Chirp{}, getError(msg, http.StatusNotFound, errors.New(msg))
	}

	return chirp, CustomError{}
}

func (db *Database) loadChirpDB() (DBStructure[Chirp], error) {
	dbs := DBStructure[Chirp]{
		Data: map[int]Chirp{},
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
