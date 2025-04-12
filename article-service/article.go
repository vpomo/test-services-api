package main

import (
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"
	"main/article-service/service"

	pb "main/internal/proto/article"
)

func main() {
	articleService := service.NewArticleService()

	grpcServer := grpc.NewServer()
	pb.RegisterArticleServiceServer(grpcServer, articleService)

	lis, err := net.Listen("tcp", "0.0.0.0::50052")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	go func() {
		log.Println("Starting ArticleService on port 50052")
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	// Graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	log.Println("Shutting down ArticleService...")
	grpcServer.GracefulStop()
	log.Println("ArticleService stopped")
}
