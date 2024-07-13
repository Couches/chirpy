package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

type chirpRequest struct {
	Body string `json:"body"`
}

type chirp struct {
	Id          int    `json:"id"`
	Valid       bool   `json:"valid"`
	CleanedBody string `json:"cleaned_body"`
}

func chirpsEndpoint(w http.ResponseWriter, request *http.Request, config apiConfig) {
	decoder := json.NewDecoder(request.Body)
	req := chirpRequest{}
	err := decoder.Decode(&req)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	if len(req.Body) > 140 {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long")
		return
	}

	contents := config.Database.ReadAll().Contents

	payload := chirp{
		Id:          len(contents) + 1,
		Valid:       true,
		CleanedBody: cleanMessage(req.Body),
	}

  config.Database.Write(payload.Id, payload)

	respondWithJSON(w, http.StatusOK, payload)
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	http.Error(w, msg, code)
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	res, err := json.Marshal(payload)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	w.Header().Add("Content-Type", "text/json; charset=utf-8")
	w.WriteHeader(code)
	w.Write(res)
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
