package main

import (
	"encoding/json"
	"net/http"

	ChirpyDatabase "github.com/Couches/chirpy-database"
)

func respondWithError(w http.ResponseWriter, result ChirpyDatabase.Result) {
	http.Error(w, result.Error.Error(), result.Code)
}

func respondWithJSON(w http.ResponseWriter, result ChirpyDatabase.Result) {
	res, err := json.Marshal(result.Body)

	if err != nil {
    error := ChirpyDatabase.GetErrorResult(http.StatusInternalServerError, err)
		respondWithError(w, error)
		return
	}

	w.Header().Add("Content-Type", "text/json; charset=utf-8")
	w.WriteHeader(result.Code)
	w.Write(res)
}

func decodeRequestBody(r *http.Request, params interface{}) ChirpyDatabase.Result {
  decoder := json.NewDecoder(r.Body)
  req := params
  err := decoder.Decode(&req)

  if err != nil {
    error := ChirpyDatabase.GetErrorResult(http.StatusInternalServerError, err)
    return error
  }

  return ChirpyDatabase.GetOKResult(1, req)
}
