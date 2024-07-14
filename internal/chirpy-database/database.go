package ChirpyDatabase

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
)

type Database struct {
	path string
	mux  *sync.RWMutex
}

type DataEntity interface {
	GetID() int
	IsValid() bool
}

type DBStructure[T any] struct {
	Data map[int]T `json:"data"`
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
		fmt.Printf("Failed to create file \"%v\": %v\n", path, err)
		return db, err
	}

	defer file.Close()

	db = &Database{
		path: path,
		mux:  &sync.RWMutex{},
	}

	return db, nil
}

func (db *Database) writeDB(dbs any) error {
	file, err := os.Create(db.path)
	defer file.Close()

	if err != nil {
		fmt.Printf("Failed to create or truncate file \"%v\"\n", db.path)
		return err
	}

	body, err := json.Marshal(dbs)
	if err != nil {
		fmt.Printf("Failed to serialize data to JSON for writing\n")
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

func (db *Database) loadDB() (DBStructure[DataEntity], error) {
	dbs := DBStructure[DataEntity]{
		Data: map[int]DataEntity{},
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

func getError(msg string, code int, err error) CustomError {
	return CustomError{
		Msg:  msg,
		Code: code,
		Err:  err,
	}
}
