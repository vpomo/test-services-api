package handler

import (
	"context"
	"encoding/json"
	"main/internal/proto/user"
	"net/http"

	"google.golang.org/grpc"
)

// UserHandler handles user-related HTTP requests.
type UserHandler struct {
	client user.UserServiceClient
}

// RegisterRequest represents a request to register a new user.
type RegisterRequest struct {
	Username string `json:"username" example:"johndoe"`
	Email    string `json:"email" example:"johndoe@example.com"`
	Password string `json:"password" example:"password123"`
}

// LoginRequest represents a request to log in a user.
type LoginRequest struct {
	Email    string `json:"email" example:"johndoe@example.com"`
	Password string `json:"password" example:"password123"`
}

// NewUserHandler creates a new UserHandler with the given gRPC client connection.
func NewUserHandler(conn *grpc.ClientConn) *UserHandler {
	return &UserHandler{
		client: user.NewUserServiceClient(conn),
	}
}

// Register godoc
// @Summary Register a new user
// @Description Registers a new user with the provided username, email, and password.
// @ID register-user
// @Accept  json
// @Produce  json
// @Param   user     body    RegisterRequest  true  "User registration details"
// @Success 200 {object} map[string]string "Returns user_id"
// @Failure 400 {string} string "Bad Request"
// @Failure 500 {string} string "Internal Server Error"
// @Router /register [post]
func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	grpcReq := &user.RegisterRequest{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
	}

	res, err := h.client.Register(context.Background(), grpcReq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"user_id": res.UserId})
}

// Login godoc
// @Summary Log in a user
// @Description Logs in a user with the provided email and password.
// @ID login-user
// @Accept  json
// @Produce  json
// @Param   user     body    LoginRequest  true  "User login details"
// @Success 200 {object} map[string]string "Returns token"
// @Failure 400 {string} string "Bad Request"
// @Failure 401 {string} string "Unauthorized"
// @Router /login [post]
func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	grpcReq := &user.LoginRequest{
		Email:    req.Email,
		Password: req.Password,
	}

	res, err := h.client.Login(context.Background(), grpcReq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"token": res.Token})
}

// ValidateToken godoc
// @Summary Validate a user token
// @Description Validates the provided user token.
// @ID validate-token
// @Produce  json
// @Param   Authorization  header  string  true  "Authorization token"
// @Success 200 {string} string "Token is valid"
// @Failure 400 {string} string "Bad Request"
// @Failure 401 {string} string "Unauthorized"
// @Failure 500 {string} string "Internal Server Error"
// @Router /validate [get]
func (h *UserHandler) ValidateToken(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	if token == "" {
		http.Error(w, "Authorization token required", http.StatusBadRequest)
		return
	}

	valid, err := h.ValidateTokenInternal(token)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if !valid {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *UserHandler) ValidateTokenInternal(token string) (bool, error) {
	grpcReq := &user.ValidateTokenRequest{
		Token: token,
	}

	res, err := h.client.ValidateToken(context.Background(), grpcReq)
	if err != nil {
		return false, err
	}

	return res.Valid, nil
}
