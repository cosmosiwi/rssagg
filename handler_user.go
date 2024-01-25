package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/cosmosiwi/rssagg/internal/database"
	"github.com/google/uuid"
)

func (apiCfg *apiConfig) HandlerCreateUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Name string `json:"name"`
	}
	decoder := json.NewDecoder(r.Body)

	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		responseErr(w, 400, fmt.Sprintf("Error parsing JSON: %v", err))
		return
	}

	user, err := apiCfg.DB.CreateUser(r.Context(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdateAt:  time.Now().UTC(),
		Name:      params.Name,
	})
	if err != nil {
		responseErr(w, 400, fmt.Sprintf("Couldn't create : %v", err))
	}

	responseWithJSON(w, 200, databaseUserToUser(user))
}
