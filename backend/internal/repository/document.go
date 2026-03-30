package repository

import (
	"context"

	graphv1 "github.com/synthify/backend/gen/synthify/graph/v1"
)

type DocumentRepository interface {
	CreateDocument(ctx context.Context, workspaceID string, filename string, mimeType string, fileSize int64) (*graphv1.Document, error)
	GetDocument(ctx context.Context, workspaceID string, documentID string) (*graphv1.Document, error)
	ListDocuments(ctx context.Context, workspaceID string) ([]*graphv1.Document, error)
	StartProcessing(ctx context.Context, workspaceID string, documentID string) (*graphv1.Document, error)
}
