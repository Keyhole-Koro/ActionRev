package mock

import (
	"context"
	"sync"
	"time"

	graphv1 "github.com/synthify/backend/gen/synthify/graph/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type DocumentRepository struct {
	mu        sync.RWMutex
	documents map[string][]*graphv1.Document
}

const mockProcessingDelay = 2 * time.Second

func NewDocumentRepository() *DocumentRepository {
	return &DocumentRepository{
		documents: map[string][]*graphv1.Document{
			"00000000-0000-4000-8000-000000000001": {
				{
					DocumentId:  "doc_demo",
					WorkspaceId: "00000000-0000-4000-8000-000000000001",
					Filename:    "growth-strategy-review.pdf",
					MimeType:    "application/pdf",
					FileSize:    1024 * 1024,
					Status:      graphv1.DocumentLifecycleState_DOCUMENT_LIFECYCLE_STATE_COMPLETED,
					CreatedAt:   timestamppb.New(time.Now().Add(-1 * time.Hour)),
					UpdatedAt:   timestamppb.New(time.Now().Add(-30 * time.Minute)),
				},
			},
		},
	}
}

func (r *DocumentRepository) CreateDocument(_ context.Context, workspaceID string, filename string, mimeType string, fileSize int64) (*graphv1.Document, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := timestamppb.New(time.Now())
	document := &graphv1.Document{
		DocumentId:  "doc_" + now.AsTime().Format("20060102150405.000000000"),
		WorkspaceId: workspaceID,
		Filename:    filename,
		MimeType:    mimeType,
		FileSize:    fileSize,
		Status:      graphv1.DocumentLifecycleState_DOCUMENT_LIFECYCLE_STATE_UPLOADED,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	r.documents[workspaceID] = append([]*graphv1.Document{document}, r.documents[workspaceID]...)
	return cloneDocument(document), nil
}

func (r *DocumentRepository) GetDocument(_ context.Context, workspaceID string, documentID string) (*graphv1.Document, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, document := range r.documents[workspaceID] {
		if document.DocumentId == documentID {
			normalizeDocumentStatus(document)
			return cloneDocument(document), nil
		}
	}

	return nil, nil
}

func (r *DocumentRepository) ListDocuments(_ context.Context, workspaceID string) ([]*graphv1.Document, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	documents := r.documents[workspaceID]
	result := make([]*graphv1.Document, 0, len(documents))
	for _, document := range documents {
		normalizeDocumentStatus(document)
		result = append(result, cloneDocument(document))
	}

	return result, nil
}

func (r *DocumentRepository) StartProcessing(_ context.Context, workspaceID string, documentID string) (*graphv1.Document, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, document := range r.documents[workspaceID] {
		if document.DocumentId == documentID {
			document.Status = graphv1.DocumentLifecycleState_DOCUMENT_LIFECYCLE_STATE_PROCESSING
			document.UpdatedAt = timestamppb.New(time.Now())
			return cloneDocument(document), nil
		}
	}

	return nil, nil
}

func cloneDocument(document *graphv1.Document) *graphv1.Document {
	if document == nil {
		return nil
	}

	return &graphv1.Document{
		DocumentId:  document.DocumentId,
		WorkspaceId: document.WorkspaceId,
		Filename:    document.Filename,
		MimeType:    document.MimeType,
		FileSize:    document.FileSize,
		Status:      document.Status,
		CreatedAt:   document.CreatedAt,
		UpdatedAt:   document.UpdatedAt,
	}
}

func normalizeDocumentStatus(document *graphv1.Document) {
	if document == nil {
		return
	}

	if document.Status != graphv1.DocumentLifecycleState_DOCUMENT_LIFECYCLE_STATE_PROCESSING {
		return
	}

	updatedAt := document.GetUpdatedAt()
	if updatedAt == nil {
		return
	}

	if time.Since(updatedAt.AsTime()) >= mockProcessingDelay {
		document.Status = graphv1.DocumentLifecycleState_DOCUMENT_LIFECYCLE_STATE_COMPLETED
		document.UpdatedAt = timestamppb.New(time.Now())
	}
}
