package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/emavillamayorpsh/rest-ws/models"
	"github.com/emavillamayorpsh/rest-ws/repository"
	"github.com/emavillamayorpsh/rest-ws/server"
	"github.com/golang-jwt/jwt"
	"github.com/gorilla/mux"
	"github.com/segmentio/ksuid"
)

type UpsertPostRequest struct {
	PostContent string `json:"post_content"`
}

type PostResponse struct {
	Id string `json:"id"`
	PostContent string `json:"post_content"`
}

type PostUpdateResponse struct {
	Message string `json:"message"`
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
			var postRequest = UpsertPostRequest{}
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

func GetPostByIdHandler(s server.Server) http.HandlerFunc {
	return func (w http.ResponseWriter, r *http.Request) {
		// get id from request url
		params := mux.Vars(r)

		post, err := repository.GetPostById(r.Context(), params["id"])
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(post)
	}
}

func UpdatePostHandler(s server.Server) http.HandlerFunc {
	return func (w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
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
			var postRequest = UpsertPostRequest{}
			if err := json.NewDecoder(r.Body).Decode(&postRequest); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			post := models.Post{
				Id: params["id"],
				PostContent: postRequest.PostContent,
				UserId: claims.UserId,
			}

			err = repository.UpdatePost(r.Context(), &post)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(PostUpdateResponse{
				Message: "Post Updated",
			})
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

	}
}

func DeletePostHandler(s server.Server) http.HandlerFunc {
	return func (w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
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
			err = repository.DeletePost(r.Context(), params["id"], claims.UserId)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(PostUpdateResponse{
				Message: "Post Deleted",
			})
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

	}
}

func ListPostHandler(s server.Server) http.HandlerFunc {
	return func (w http.ResponseWriter, r *http.Request) {
		var err error

		// get "query param" from url
		pageStr :=  r.URL.Query().Get("page")
		var page = uint64(0)

		// validate that it has a page number
		if pageStr != "" {
			page, err =  strconv.ParseUint(pageStr,10, 64)
			// validate that it is a "valid" number
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
		}

		posts, err := repository.ListPost(r.Context(), page)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(posts)
	}
}