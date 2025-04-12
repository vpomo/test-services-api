package main

import (
	"context"
	httpSwagger "github.com/swaggo/http-swagger"
	"log"
	handler2 "main/api-gateway/handler"
	"main/api-gateway/middleware"
	"main/internal/proto/user"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
	"google.golang.org/grpc"
	_ "main/api-gateway/docs"
)

// @title Swagger Example API
// @version 1.0
// @description This is test task

// @contact.name API Support

// @host localhost:8080
// @BasePath /api
func main() {
	userConn, err := grpc.Dial("user-service:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect to user service: %v", err)
	}
	userClient := user.NewUserServiceClient(userConn)

	defer userConn.Close()

	articleConn, err := grpc.Dial("article-service:50052", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect to article service: %v", err)
	}
	defer articleConn.Close()

	userHandler := handler2.NewUserHandler(userConn)
	articleHandler := handler2.NewArticleHandler(articleConn)

	router := mux.NewRouter()

	router.PathPrefix("/swagger/").Handler(httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
		httpSwagger.DocExpansion("none"),
	))

	publicRouter := router.PathPrefix("/api").Subrouter()
	publicRouter.HandleFunc("/register", userHandler.Register).Methods("POST")
	publicRouter.HandleFunc("/login", userHandler.Login).Methods("POST")
	publicRouter.HandleFunc("/validate", userHandler.ValidateToken).Methods("GET")

	privateRouter := router.PathPrefix("/api").Subrouter()
	privateRouter.Use(middleware.AuthMiddleware(userClient))
	privateRouter.HandleFunc("/articles", articleHandler.CreateArticle).Methods("POST")
	privateRouter.HandleFunc("/articles/{id}", articleHandler.GetArticle).Methods("GET")
	privateRouter.HandleFunc("/articles/{id}", articleHandler.UpdateArticle).Methods("PUT")
	privateRouter.HandleFunc("/articles/{id}/comments", articleHandler.AddComment).Methods("POST")

	srv := &http.Server{
		Addr:         "0.0.0.0:8080",
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	go func() {
		log.Println("Starting server on port 8080")
		if err := srv.ListenAndServe(); err != nil {
			log.Printf("Server error: %v\n", err)
		}
	}()

	// Graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	srv.Shutdown(ctx)
	log.Println("Server gracefully stopped")
}
