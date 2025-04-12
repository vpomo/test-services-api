package handler

import (
	"context"
	"encoding/json"
	"main/internal/proto/article"
	"net/http"

	"github.com/gorilla/mux"
	"google.golang.org/grpc"
)

// ArticleHandler handles article-related HTTP requests.
type ArticleHandler struct {
	client article.ArticleServiceClient
}

// CreateArticleRequest represents a request to create a new article.
type CreateArticleRequest struct {
	Title   string `json:"title" example:"My Article Title"`
	Content string `json:"content" example:"This is the content of the article."`
}

// UpdateArticleRequest represents a request to update an existing article.
type UpdateArticleRequest struct {
	Title   string `json:"title" example:"Updated Article Title"`
	Content string `json:"content" example:"This is the updated content of the article."`
}

// AddCommentRequest represents a request to add a comment to an article.
type AddCommentRequest struct {
	Content string `json:"content" example:"This is a comment."`
}

// NewArticleHandler creates a new ArticleHandler with the given gRPC client connection.
func NewArticleHandler(conn *grpc.ClientConn) *ArticleHandler {
	return &ArticleHandler{
		client: article.NewArticleServiceClient(conn),
	}
}

// CreateArticle godoc
// @Summary Create a new article
// @Description Creates a new article with the provided title and content.
// @ID create-article
// @Accept  json
// @Produce  json
// @Param   article     body    CreateArticleRequest  true  "Article details"
// @Success 200 {object} article.Article "Returns created article"
// @Failure 400 {string} string "Bad Request"
// @Failure 500 {string} string "Internal Server Error"
// @Router /articles [post]
func (h *ArticleHandler) CreateArticle(w http.ResponseWriter, r *http.Request) {
	var req CreateArticleRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userID := r.Context().Value("user_id").(string)

	grpcReq := &article.CreateArticleRequest{
		Title:    req.Title,
		Content:  req.Content,
		AuthorId: userID,
	}

	res, err := h.client.CreateArticle(context.Background(), grpcReq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(res.Article)
}

// GetArticle godoc
// @Summary Get an article by ID
// @Description Retrieves an article by its ID.
// @ID get-article
// @Produce  json
// @Param   id  path    string  true  "Article ID"
// @Success 200 {object} article.Article "Returns the article"
// @Failure 500 {string} string "Internal Server Error"
// @Router /articles/{id} [get]
func (h *ArticleHandler) GetArticle(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	articleID := vars["id"]

	grpcReq := &article.GetArticleRequest{
		Id: articleID,
	}

	res, err := h.client.GetArticle(context.Background(), grpcReq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(res.Article)
}

// UpdateArticle godoc
// @Summary Update an existing article
// @Description Updates an article with the provided title and content.
// @ID update-article
// @Accept  json
// @Produce  json
// @Param   id       path    string  true  "Article ID"
// @Param   article  body    UpdateArticleRequest  true  "Updated article details"
// @Success 200 {object} article.Article "Returns updated article"
// @Failure 400 {string} string "Bad Request"
// @Failure 500 {string} string "Internal Server Error"
// @Router /articles/{id} [put]
func (h *ArticleHandler) UpdateArticle(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	articleID := vars["id"]

	var req UpdateArticleRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userID := r.Context().Value("user_id").(string)

	grpcReq := &article.UpdateArticleRequest{
		Id:       articleID,
		Title:    req.Title,
		Content:  req.Content,
		AuthorId: userID,
	}

	res, err := h.client.UpdateArticle(context.Background(), grpcReq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(res.Article)
}

// AddComment godoc
// @Summary Add a comment to an article
// @Description Adds a comment to an article with the provided content.
// @ID add-comment
// @Accept  json
// @Produce  json
// @Param   id       path    string  true  "Article ID"
// @Param   comment  body    AddCommentRequest  true  "Comment details"
// @Success 200 {object} article.Comment "Returns added comment"
// @Failure 400 {string} string "Bad Request"
// @Failure 500 {string} string "Internal Server Error"
// @Router /articles/{id}/comments [post]
func (h *ArticleHandler) AddComment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	articleID := vars["id"]

	var req AddCommentRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userID := r.Context().Value("user_id").(string)

	grpcReq := &article.AddCommentRequest{
		ArticleId: articleID,
		Content:   req.Content,
		AuthorId:  userID,
	}

	res, err := h.client.AddComment(context.Background(), grpcReq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(res.Comment)
}
