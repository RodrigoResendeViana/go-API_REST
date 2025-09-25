package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/RodrigoResendeViana/rss-aggregator/internal/database"
	"github.com/go-chi/chi"
	"github.com/google/uuid"
)

func (apiCfg *apiConfig) handlerCreateFeedFollow(w http.ResponseWriter, r *http.Request, user database.User) {
	type parameters struct {
		FeedID uuid.UUID `json:"feedIs"`
	}

	decoder := json.NewDecoder(r.Body)

	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error parsing JSON: %v", err))
		return
	}

	feed, err := apiCfg.DB.CreateFeedFollow(r.Context(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    user.ID,
		FeedID:    params.FeedID,
	})

	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Couldn't create feed follow %v", err))
		return

	}

	respondWithJSON(w, 201, databaseFeedFollowToFeedFollow(feed))
}

func (apiCfg *apiConfig) handlerGetFeedsUserFollows(w http.ResponseWriter, r *http.Request, user database.User) {
	feeds, err := apiCfg.DB.GetFeedsUserFollows(r.Context(), user.ID)

	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Couldn't get feeds a user follows %v", err))
		return
	}

	respondWithJSON(w, 200, databaseFeedFollowsToFeedFollows(feeds))
}

func (apiCfg *apiConfig) handlerDeleteFeedFollow(w http.ResponseWriter, r *http.Request, user database.User) {
	feedFollowIDStr := chi.URLParam(r, "feedFollowID")
	feedFollowID, err := uuid.Parse(feedFollowIDStr)

	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Couldn't parse feed follow ID %v", err))
		return
	}

	err = apiCfg.DB.DeleteFeedFollow(r.Context(), database.DeleteFeedFollowParams{
		ID:     feedFollowID,
		UserID: user.ID,
	})

	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Couldn't delete feed follow: %v", err))
		return
	}

	respondWithJSON(w, 200, struct{}{})
}
