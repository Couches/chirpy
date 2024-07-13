package ChirpDatabase

import (
	"encoding/json"
	"fmt"
	"os"
)

type DatabaseStructure[T any] struct {
  Contents map[int]T
}

type Database[T any] struct {
  Path string
}

func CreateDatabase[T any](path string) *Database[T] {
  fmt.Println("New database created successfully")
  return &Database[T]{}
}

func (db *Database[T]) Write(id int, data T) {
  fmt.Println("attempting to write data:")
  fmt.Println(data)
  dbs := db.ReadAll()
  file, err := os.Create(db.Path)
  if err != nil {
    fmt.Printf("There was an issue writing to the database")
    return
  }

  dbs.Contents[id] = data

  write_content, err := json.Marshal(dbs.Contents)
  if err != nil {
    fmt.Printf("There was an issue writing to the database")
    return
  }

  fmt.Println(write_content)
  file.Write(write_content)
}

func (db *Database[T]) ReadAll() DatabaseStructure[T] {
  file, err := os.ReadFile(db.Path)
  if err != nil {
    fmt.Printf("There was an issue reading the file: %v\n", db.Path)
    return DatabaseStructure[T]{}
  }

  contents := DatabaseStructure[T]{}
  err = json.Unmarshal(file, &contents)
  if err != nil {
    fmt.Printf("There was an issue decoding the database\n")
    return DatabaseStructure[T]{}
  }

  return contents
}

