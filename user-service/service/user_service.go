package service

import (
	"context"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	pb "main/internal/proto/user"
	"sync"
	"time"
)

type User struct {
	ID        string
	Username  string
	Email     string
	Password  string
	CreatedAt time.Time
}

type UserService struct {
	pb.UnimplementedUserServiceServer
	users      map[string]*User // in-memory
	usersMutex sync.RWMutex     // RWMutex
	jwtSecret  []byte
}

func NewUserService(jwtSecret string) *UserService {
	return &UserService{
		users:     make(map[string]*User),
		jwtSecret: []byte(jwtSecret),
	}
}

func (s *UserService) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	if req.Username == "" || req.Email == "" || req.Password == "" {
		return nil, errors.New("username, email and password are required")
	}

	s.usersMutex.RLock()
	for _, user := range s.users {
		if user.Email == req.Email {
			s.usersMutex.RUnlock()
			return nil, errors.New("email already exists")
		}
	}
	s.usersMutex.RUnlock()

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("failed to hash password")
	}

	user := &User{
		ID:        uuid.New().String(),
		Username:  req.Username,
		Email:     req.Email,
		Password:  string(hashedPassword),
		CreatedAt: time.Now(),
	}

	s.usersMutex.Lock()
	s.users[user.ID] = user
	s.usersMutex.Unlock()

	return &pb.RegisterResponse{UserId: user.ID}, nil
}

func (s *UserService) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	s.usersMutex.RLock()
	var user *User
	for _, u := range s.users {
		if u.Email == req.Email {
			user = u
			break
		}
	}
	s.usersMutex.RUnlock()

	if user == nil {
		return nil, errors.New("invalid email or password")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, errors.New("invalid email or password")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":  user.ID,
		"username": user.Username,
		"email":    user.Email,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
		"iat":      time.Now().Unix(),
	})

	tokenString, err := token.SignedString(s.jwtSecret)
	if err != nil {
		return nil, errors.New("failed to generate token")
	}

	return &pb.LoginResponse{Token: tokenString}, nil
}

func (s *UserService) ValidateToken(ctx context.Context, req *pb.ValidateTokenRequest) (*pb.ValidateTokenResponse, error) {
	token, err := jwt.Parse(req.Token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return s.jwtSecret, nil
	})

	if err != nil {
		return &pb.ValidateTokenResponse{Valid: false}, nil
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID, ok := claims["user_id"].(string)
		if !ok {
			return &pb.ValidateTokenResponse{Valid: false}, nil
		}

		s.usersMutex.RLock()
		user, exists := s.users[userID]
		s.usersMutex.RUnlock()

		if !exists {
			return &pb.ValidateTokenResponse{Valid: false}, nil
		}

		return &pb.ValidateTokenResponse{
			UserId:   user.ID,
			Username: user.Username,
			Valid:    true,
		}, nil
	}

	return &pb.ValidateTokenResponse{Valid: false}, nil
}
