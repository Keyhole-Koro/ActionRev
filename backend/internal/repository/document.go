package repository

import (
	"context"
	"time"

	graphv1 "github.com/synthify/backend/gen/synthify/graph/v1"
)

type UploadTarget struct {
	Token       string
	WorkspaceID string
	Filename    string
	MimeType    string
	FileSize    int64
	ExpiresAt   time.Time
}

type DocumentRepository interface {
	CreateDocument(ctx context.Context, workspaceID string, filename string, mimeType string, fileSize int64) (*graphv1.Document, error)
	CreateUploadTarget(ctx context.Context, workspaceID string, filename string, mimeType string, fileSize int64) (*UploadTarget, error)
	ConsumeUploadTarget(ctx context.Context, token string) (*UploadTarget, error)
	GetDocument(ctx context.Context, workspaceID string, documentID string) (*graphv1.Document, error)
	ListDocuments(ctx context.Context, workspaceID string) ([]*graphv1.Document, error)
	StartProcessing(ctx context.Context, workspaceID string, documentID string) (*graphv1.Document, error)
}
