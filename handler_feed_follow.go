package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/cosmosiwi/rssagg/internal/database"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func (apiCfg *apiConfig) handlerCreateFeedFollow(w http.ResponseWriter, r *http.Request, user database.User) {
	type parameters struct {
		FeedID uuid.UUID `json:"feed_id"`
	}

	decoder := json.NewDecoder(r.Body)

	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		responseErr(w, 400, fmt.Sprintf("Error parsing json: %v", err))
	}

	feedFollow, err := apiCfg.DB.CreateFeedFollow(r.Context(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdateAt:  time.Now().UTC(),
		UserID:    user.ID,
		FeedID:    params.FeedID,
	})

	if err != nil {
		responseErr(w, 400, fmt.Sprintf("couldn't create feed follow: %v", err))
	}

	responseWithJSON(w, 200, databaseFeedFollowToFeedFollow(feedFollow))

}

func (apiCfg *apiConfig) handlerGetFeedFollows(w http.ResponseWriter, r *http.Request, user database.User){
	feedFollows, err := apiCfg.DB.GetFeedFollows(r.Context(),user.ID)
	if err != nil {
		responseErr(w, 400, fmt.Sprintf("Couldn't get feedFollows"))
	}

	responseWithJSON(w, 201, databaseFeedFollowsToFeedFollows(feedFollows))
}

func (apiCfg *apiConfig) handlerDeleteFeedFollow(w http.ResponseWriter, r *http.Request, user database.User){
	feedFollowIDStr := chi.URLParam(r, "feedFollowID")
	feedFollowID, err := uuid.Parse(feedFollowIDStr)
	if err != nil {
		responseErr(w, 400, fmt.Sprintf("Couldn't parse feed follow id: %v", err))
		return
	}

	err = apiCfg.DB.DeleteFeedFollows(r.Context(), database.DeleteFeedFollowsParams{
		ID: feedFollowID,
		UserID: user.ID,
	})
	if err != nil {
		responseErr(w, 400, fmt.Sprintf("Couldn't delete feed follow: %v", err))
		return
	}
	responseWithJSON(w, 200, struct{}{})
}