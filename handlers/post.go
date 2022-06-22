package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/leonardoas10/go-rest-ws/helpers"
	"github.com/leonardoas10/go-rest-ws/models"
	"github.com/leonardoas10/go-rest-ws/repository"
	"github.com/leonardoas10/go-rest-ws/server"
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

func InsertPostHandler(s server.Server) http.HandlerFunc  {
	return func (w http.ResponseWriter, r *http.Request)  {
		token, err := helpers.GetToken(s,w,r)
	
		if err != nil {
			return
		}

		if claims, ok := token.Claims.(*models.AppClaims); ok && token.Valid {
			var postRequest = UpsertPostRequest{}
			if err := json.NewDecoder(r.Body).Decode(&postRequest); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
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
			var postMessage = models.WebsocketMessage{
				Type: "Post_Created",
				Payload: post,
			}
			s.Hub().Broadcast(postMessage, nil)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(PostResponse{
				Id: post.Id,
				PostContent: post.PostContent,
			})

		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func GetPostByIdHandler(s server.Server) http.HandlerFunc  {
	return func (w http.ResponseWriter, r *http.Request)  {
		params := mux.Vars(r)
		post, err := repository.GetPostById(r.Context(), params["id"])
		log.Println(post)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(post)
	}
}

func UpdatePostHandler(s server.Server) http.HandlerFunc  {
	return func (w http.ResponseWriter, r *http.Request)  {
		params := mux.Vars(r)
		token, err := helpers.GetToken(s,w,r)

		if err != nil {
			return
		}

		if claims, ok := token.Claims.(*models.AppClaims); ok && token.Valid {
			var postRequest = UpsertPostRequest{}
			if err := json.NewDecoder(r.Body).Decode(&postRequest); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
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
				Message: "Post updated",
			})

		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func DeletePostHandler(s server.Server) http.HandlerFunc  {
	return func (w http.ResponseWriter, r *http.Request)  {
		params := mux.Vars(r)
		token, err := helpers.GetToken(s,w,r)

		if err != nil {
			return
		}

		if claims, ok := token.Claims.(*models.AppClaims); ok && token.Valid {
			err = repository.DeletePost(r.Context(), params["id"], claims.UserId)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(PostUpdateResponse{
				Message: "Post deleted",
			})

		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func ListPostHandler(s server.Server) http.HandlerFunc  {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error
		pageStr := r.URL.Query().Get("page")
		var page = uint64(0)
		if pageStr != "" {
			page, err = strconv.ParseUint(pageStr, 10, 64)
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
