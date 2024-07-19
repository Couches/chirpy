package ChirpyDatabase

import (
	"cmp"
	"errors"
	"fmt"
	"net/http"
	"slices"
)

type Chirp struct {
	Id       int    `json:"id"`
	Body     string `json:"body"`
	AuthorID int    `json:"author_id"`
}

func (db *Database) CreateChirp(body string, authorID int) Result {
	result := db.LoadDB()
	if result.Error != nil {
		return result
	}

	structure := (*result.Body).(DatabaseStructure)

	chirp := Chirp{
		Id:       *db.nextChirpID,
		Body:     body,
		AuthorID: authorID,
	}

	fmt.Println(chirp)

	structure.Chirps[chirp.Id] = chirp

	result = db.WriteDB(structure)
	if result.Error != nil {
		return result
	}

	*db.nextChirpID++

	return GetOKResult(http.StatusCreated, chirp)
}

func (db *Database) DeleteChirp(chirpID, userID int) Result {
	result := db.LoadDB()
	if result.Error != nil {
		return result
	}

	structure := (*result.Body).(DatabaseStructure)

	if _, ok := structure.Chirps[chirpID]; !ok {
		msg := fmt.Sprintf("No chirp found with id \"%v\"", chirpID)
		return GetErrorResult(http.StatusUnauthorized, errors.New(msg))
	}

	chirp := structure.Chirps[chirpID]
	if chirp.AuthorID != userID {
		msg := fmt.Sprintf("Cannot delete others' chirps")
		return GetErrorResult(http.StatusForbidden, errors.New(msg))
	}

	delete(structure.Chirps, chirpID)

	result = db.WriteDB(structure)
	if result.Error != nil {
		return result
	}

	return GetOKResult(http.StatusNoContent, nil)
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

	chirps := []Chirp{}
	for _, chirp := range structure.Chirps {
		chirps = append(chirps, chirp)
	}

	return GetOKResult(http.StatusOK, chirps)
}

func (db *Database) GetChirpsFromAuthor(authorID int, sortBy string) Result {
	result := db.GetAllChirps()
	if result.Error != nil {
		return result
	}

	chirps := (*result.Body).([]Chirp)

	if authorID != 0 {
		filtered := []Chirp{}

		for _, chirp := range chirps {
			if chirp.AuthorID == authorID {
				filtered = append(filtered, chirp)
			}
		}

		chirps = filtered
	}

	slices.SortFunc(chirps, func(a, b Chirp) int {
		comparison := cmp.Compare(a.AuthorID, b.AuthorID)
		if sortBy == "desc" {
			comparison *= -1
		}
		return comparison
	})

	return GetOKResult(http.StatusOK, chirps)
}
