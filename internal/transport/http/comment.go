package http

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/BFamzz/comments-api/internal/comment"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

type CommentService interface {
	PostComment(context.Context, comment.Comment) (comment.Comment, error)
	GetComment(context.Context, string) (comment.Comment, error)
	UpdateComment(context.Context, string, comment.Comment) (comment.Comment, error)
	DeleteComment(context.Context, string) error
}

type Response struct {
	Message string
}

type PostCommentRequest struct {
	Slug string `json:"slug" validate:"required"`
	Author string `json:"author" validate:"required"`
	Body string `json:"body" validate:"required"`
}

func convertPostCommentRequestToComment(c PostCommentRequest) comment.Comment {
	return comment.Comment{
		Slug: c.Slug,
		Author: c.Author,
		Body: c.Body,
	}
}

func (h *Handler) PostComment(w http.ResponseWriter, r *http.Request) {
	var newComment PostCommentRequest
	if err := json.NewDecoder(r.Body).Decode(&newComment); err != nil {
		return
	}

	validator := validator.New()
	err := validator.Struct(newComment)
	if err != nil {
		http.Error(w, "not a valid comment", http.StatusBadRequest)
		return
	}

	convertedComment := convertPostCommentRequestToComment(newComment)

	postedComment, err := h.Service.PostComment(r.Context(), convertedComment)
	if err != nil {
		log.Print(err)
		return
	}

	if err := json.NewEncoder(w).Encode(postedComment); err != nil {
		panic(err)
	}
}

func (h *Handler) GetComment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	comment, err := h.Service.GetComment(r.Context(), id)
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(comment); err != nil {
		panic(err)
	}

}

func (h *Handler) UpdateComment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
	}

	var updatedComment comment.Comment
	if err := json.NewDecoder(r.Body).Decode(&updatedComment); err != nil {
		return
	}

	comment, err := h.Service.UpdateComment(r.Context(), id, updatedComment)
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(comment); err != nil {
		panic(err)
	}
}

func (h *Handler) DeleteComment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
	}

	err := h.Service.DeleteComment(r.Context(), id)
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(Response{Message: "Successfully deleted"}); err != nil {
		panic(err)
	}
}