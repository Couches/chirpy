package ChirpyDatabase

import (
	"encoding/json"
	"net/http"
	"os"
	"sync"
)

type Database struct {
	path string
	mux  *sync.RWMutex
}

type DatabaseStructure struct {
	Chirps map[int]Chirp
	Users  map[int]User
}

func NewDB(path string) Result {
	file, err := os.Create(path)
	if err != nil {
		return GetErrorResult(http.StatusInternalServerError, err)
	}

	defer file.Close()

  db := &Database{
		path: path,
		mux:  &sync.RWMutex{},
	}

  return GetOKResult(1, *db)
}



func (db *Database) WriteDB(structure DatabaseStructure) Result {
	file, err := os.Create(db.path)
	if err != nil {
		return GetErrorResult(http.StatusInternalServerError, err)
	}

  defer file.Close()

	body, err := json.Marshal(structure)
	if err != nil {
		return GetErrorResult(http.StatusInternalServerError, err)
	}

	numWritten, err := file.Write(body)
	if err != nil {
		return GetErrorResult(http.StatusInternalServerError, err)
	}

	return GetOKResult(1, numWritten)
}



func (db *Database) LoadDB() Result {
	file, err := os.ReadFile(db.path)
	if err != nil {
		return GetErrorResult(http.StatusInternalServerError, err)
	}

  structure := DatabaseStructure{
    Users: map[int]User{},
    Chirps: map[int]Chirp{},
  }
	err = json.Unmarshal(file, &structure)
	if err != nil {
    if err.Error() != "unexpected end of JSON input" {
		  return GetErrorResult(http.StatusInternalServerError, err)
    }
	}

  return GetOKResult(1, structure)
}
