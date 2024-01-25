package main

import (
	"net/http"
	"fmt"

	"github.com/cosmosiwi/rssagg/internal/database"
	"github.com/cosmosiwi/rssagg/internal/auth"
)

type authedHandler func(http.ResponseWriter, *http.Request, database.User)

func (apiCfg *apiConfig) middlewareAuth(handler authedHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		apiKey, err := auth.GetAPIKey(r.Header)
		if err != nil {
			responseErr(w, 403, fmt.Sprintf("Auth error : %v", err))
			return
		}

		user, err := apiCfg.DB.GetUserByAPIKey(r.Context(), apiKey)
		if err != nil {
			responseErr(w, 404, fmt.Sprintf("Couldn't get user: %v", err))
		}

		handler(w, r, user)
	}
}
