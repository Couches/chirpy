package ChirpDatabase

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
)

type Chirp struct {
	Id          int    `json:"id"`
	Valid       bool   `json:"valid"`
	CleanedBody string `json:"cleaned_body"`
}

type Database struct {
	path string
	mux  *sync.RWMutex
}

type DBStructure struct {
  Chirps map[int]Chirp `json:"chirps"`
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
    mux: &sync.RWMutex{},
  }

  return db, nil
}

func (db *Database) CreateChirp(body string) (Chirp, error) {
  chirp := Chirp{}
  err := json.Unmarshal([]byte(body), &chirp)
  if err != nil {
    fmt.Printf("Failed to deserialize Chirp\n")
    return chirp, err
  }

  return chirp, nil
}

func (db *Database) GetChirps() ([]Chirp, error) {
  dbs, err := db.loadDB()

  if err != nil {
    fmt.Printf("Failed to load database\n")
    return nil, err
  }

  chirps := make([]Chirp, 0, len(dbs.Chirps))
  for _, chirp := range dbs.Chirps {
    chirps = append(chirps, chirp)
  }

  return chirps, err
}

func (db *Database) loadDB() (DBStructure, error) {
  dbs := DBStructure{}
  file, err := os.ReadFile(db.path)
  if err != nil {
    fmt.Printf("Failed to read file \"%v\"\n", db.path)
    return dbs, err
  }

  err = json.Unmarshal(file, &dbs)

  if err != nil {
    fmt.Printf("Failed to deserialize JSON from file \"%v\"\n", db.path)
    return dbs, err
  }

  fmt.Println(dbs.Chirps)
  return dbs, nil
}

func (db *Database) writeDB(dbs DBStructure) error {

}
