package mock

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"sync"
	"time"

	graphv1 "github.com/synthify/backend/gen/synthify/graph/v1"
	"github.com/synthify/backend/internal/repository"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type DocumentRepository struct {
	mu            sync.RWMutex
	documents     map[string][]*graphv1.Document
	uploadTargets map[string]*repository.UploadTarget
}

const mockProcessingDelay = 2 * time.Second
const mockUploadTargetTTL = 15 * time.Minute

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
		uploadTargets: map[string]*repository.UploadTarget{},
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

func (r *DocumentRepository) CreateUploadTarget(_ context.Context, workspaceID string, filename string, mimeType string, fileSize int64) (*repository.UploadTarget, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	token, err := newUploadToken()
	if err != nil {
		return nil, err
	}

	target := &repository.UploadTarget{
		Token:       token,
		WorkspaceID: workspaceID,
		Filename:    filename,
		MimeType:    mimeType,
		FileSize:    fileSize,
		ExpiresAt:   time.Now().Add(mockUploadTargetTTL),
	}
	r.uploadTargets[token] = target

	return cloneUploadTarget(target), nil
}

func (r *DocumentRepository) ConsumeUploadTarget(_ context.Context, token string) (*repository.UploadTarget, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	target, ok := r.uploadTargets[token]
	if !ok {
		return nil, nil
	}

	if time.Now().After(target.ExpiresAt) {
		delete(r.uploadTargets, token)
		return nil, errors.New("upload token expired")
	}

	delete(r.uploadTargets, token)
	return cloneUploadTarget(target), nil
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

func cloneUploadTarget(target *repository.UploadTarget) *repository.UploadTarget {
	if target == nil {
		return nil
	}

	return &repository.UploadTarget{
		Token:       target.Token,
		WorkspaceID: target.WorkspaceID,
		Filename:    target.Filename,
		MimeType:    target.MimeType,
		FileSize:    target.FileSize,
		ExpiresAt:   target.ExpiresAt,
	}
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

func newUploadToken() (string, error) {
	bytes := make([]byte, 12)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	return "upl_" + hex.EncodeToString(bytes), nil
}
