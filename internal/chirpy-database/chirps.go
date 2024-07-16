package ChirpyDatabase

import (
	"errors"
	"fmt"
	"net/http"
)

type Chirp struct {
	Id    int    `json:"id"`
	Valid bool   `json:"valid"`
	Body  string `json:"body"`
}


func (db *Database) CreateChirp(body string) Result {
  result := db.LoadDB()
	if result.Error != nil {
		return result
	}

	structure := (*result.Body).(DatabaseStructure)

  chirp := Chirp {
    Id: len(structure.Chirps) + 1,
    Body: body,
  }

	structure.Chirps[chirp.Id] = chirp

	result = db.WriteDB(structure)
	if result.Error != nil {
		return result
	}

	return GetOKResult(http.StatusCreated, chirp)
}


func (db *Database) GetChirp(id int) Result {
	result := db.GetAllChirps()
	if result.Error != nil {
		return result
	}

  chirps := (*result.Body).(map[int]Chirp)
	if chirp, ok := chirps[id]; ok {
		return GetOKResult(http.StatusOK, chirp)
	}

	msg := fmt.Sprintf("No chirp was found with id \"%v\"", id)
	return GetErrorResult(http.StatusNotFound, errors.New(msg))
}


func (db *Database) GetAllChirps() Result {
	result := db.LoadDB()
	if result.Error != nil {
		return result
	}

	structure := (*result.Body).(DatabaseStructure)
	if len(structure.Chirps) == 0 {
		msg := "No chirps found"
		return GetErrorResult(http.StatusNoContent, errors.New(msg))
	}

	return GetOKResult(http.StatusOK, structure.Chirps)
}
