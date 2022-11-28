package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/emavillamayorpsh/rest-ws/models"
	"github.com/emavillamayorpsh/rest-ws/repository"
	"github.com/emavillamayorpsh/rest-ws/server"
	"github.com/golang-jwt/jwt"
	"github.com/segmentio/ksuid"
)

type InsertPostRequest struct {
	PostContent string `json:"post_content"`
}

type PostResponse struct {
	Id string `json:"id"`
	PostContent string `json:"post_content"`
}

func InsertPostHandler(s server.Server) http.HandlerFunc {
	return func (w http.ResponseWriter, r *http.Request) {
		// GET THE TOKEN FROM AUTHORIZATION
		tokenString := strings.TrimSpace(r.Header.Get("Authorization"))

		// CHECK IF TOKEN IS VALID
		token, err := jwt.ParseWithClaims(tokenString, &models.AppClaims{}, func(t *jwt.Token) (interface{}, error) {
			return []byte(s.Config().JWTSecret), nil
		})

		// IN CASE TOKEN INVALID RETURN ERROR
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		// DESTRUCTURES THE TOKEN IN ORDER TO GET THE KEY/VALUES OF IT
		if claims, ok := token.Claims.(*models.AppClaims); ok && token.Valid{
			var postRequest = InsertPostRequest{}
			if err := json.NewDecoder(r.Body).Decode(&postRequest); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			id, err := ksuid.NewRandom()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			post := models.Post{
				Id: id.String(),
				PostContent: postRequest.PostContent,
				UserId: claims.UserId,
			}

			err = repository.InsertPost(r.Context(), &post)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(PostResponse{
				Id: post.Id,
				PostContent: post.PostContent,
			})
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

	}
}