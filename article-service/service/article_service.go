package service

import (
	"context"
	"errors"
	"github.com/google/uuid"
	pb "main/internal/proto/article"
	"sync"
	"time"
)

type Article struct {
	ID        string
	Title     string
	Content   string
	AuthorID  string
	CreatedAt time.Time
	UpdatedAt time.Time
	Comments  []*Comment
}

type Comment struct {
	ID        string
	Content   string
	AuthorID  string
	ArticleID string
	CreatedAt time.Time
}

type ArticleService struct {
	pb.UnimplementedArticleServiceServer
	articles     map[string]*Article // in-memory
	comments     map[string]*Comment // in-memory
	storageMutex sync.RWMutex
}

func NewArticleService() *ArticleService {
	return &ArticleService{
		articles: make(map[string]*Article),
		comments: make(map[string]*Comment),
	}
}

func (s *ArticleService) CreateArticle(ctx context.Context, req *pb.CreateArticleRequest) (*pb.CreateArticleResponse, error) {
	if req.Title == "" || req.Content == "" || req.AuthorId == "" {
		return nil, errors.New("title, content and author_id are required")
	}

	now := time.Now()
	article := &Article{
		ID:        uuid.New().String(),
		Title:     req.Title,
		Content:   req.Content,
		AuthorID:  req.AuthorId,
		CreatedAt: now,
		UpdatedAt: now,
		Comments:  []*Comment{},
	}

	s.storageMutex.Lock()
	s.articles[article.ID] = article
	s.storageMutex.Unlock()

	return &pb.CreateArticleResponse{
		Article: &pb.Article{
			Id:       article.ID,
			Title:    article.Title,
			Content:  article.Content,
			AuthorId: article.AuthorID,
		},
	}, nil
}

func (s *ArticleService) GetArticle(ctx context.Context, req *pb.GetArticleRequest) (*pb.GetArticleResponse, error) {
	s.storageMutex.RLock()
	article, exists := s.articles[req.Id]
	s.storageMutex.RUnlock()

	if !exists {
		return nil, errors.New("article not found")
	}

	var pbComments []*pb.Comment
	for _, comment := range article.Comments {
		pbComments = append(pbComments, &pb.Comment{
			Id:        comment.ID,
			Content:   comment.Content,
			AuthorId:  comment.AuthorID,
			ArticleId: comment.ArticleID,
		})
	}

	return &pb.GetArticleResponse{
		Article: &pb.Article{
			Id:       article.ID,
			Title:    article.Title,
			Content:  article.Content,
			AuthorId: article.AuthorID,
			Comments: pbComments,
		},
	}, nil
}

func (s *ArticleService) UpdateArticle(ctx context.Context, req *pb.UpdateArticleRequest) (*pb.UpdateArticleResponse, error) {
	if req.Id == "" || req.Title == "" || req.Content == "" || req.AuthorId == "" {
		return nil, errors.New("id, title, content and author_id are required")
	}

	s.storageMutex.Lock()
	defer s.storageMutex.Unlock()

	article, exists := s.articles[req.Id]
	if !exists {
		return nil, errors.New("article not found")
	}

	if article.AuthorID != req.AuthorId {
		return nil, errors.New("only author can update the article")
	}

	article.Title = req.Title
	article.Content = req.Content
	article.UpdatedAt = time.Now()

	return &pb.UpdateArticleResponse{
		Article: &pb.Article{
			Id:       article.ID,
			Title:    article.Title,
			Content:  article.Content,
			AuthorId: article.AuthorID,
		},
	}, nil
}

func (s *ArticleService) AddComment(ctx context.Context, req *pb.AddCommentRequest) (*pb.AddCommentResponse, error) {
	if req.ArticleId == "" || req.Content == "" || req.AuthorId == "" {
		return nil, errors.New("article_id, content and author_id are required")
	}

	s.storageMutex.Lock()
	defer s.storageMutex.Unlock()

	_, exists := s.articles[req.ArticleId]
	if !exists {
		return nil, errors.New("article not found")
	}

	now := time.Now()
	comment := &Comment{
		ID:        uuid.New().String(),
		Content:   req.Content,
		AuthorID:  req.AuthorId,
		ArticleID: req.ArticleId,
		CreatedAt: now,
	}

	s.comments[comment.ID] = comment

	s.articles[req.ArticleId].Comments = append(s.articles[req.ArticleId].Comments, comment)

	return &pb.AddCommentResponse{
		Comment: &pb.Comment{
			Id:        comment.ID,
			Content:   comment.Content,
			AuthorId:  comment.AuthorID,
			ArticleId: comment.ArticleID,
		},
	}, nil
}
