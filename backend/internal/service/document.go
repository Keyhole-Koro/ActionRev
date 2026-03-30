package service

import (
	"context"
	"errors"
	"strings"

	graphv1 "github.com/synthify/backend/gen/synthify/graph/v1"
	"github.com/synthify/backend/internal/repository"
)

type DocumentService struct {
	repo repository.DocumentRepository
}

func NewDocumentService(repo repository.DocumentRepository) *DocumentService {
	return &DocumentService{repo: repo}
}

func (s *DocumentService) CreateDocument(ctx context.Context, req *graphv1.CreateDocumentRequest) (*graphv1.CreateDocumentResponse, error) {
	if strings.TrimSpace(req.GetWorkspaceId()) == "" {
		return nil, errors.New("workspace id is required")
	}
	if strings.TrimSpace(req.GetFilename()) == "" {
		return nil, errors.New("filename is required")
	}

	document, err := s.repo.CreateDocument(ctx, req.GetWorkspaceId(), req.GetFilename(), req.GetMimeType(), req.GetFileSize())
	if err != nil {
		return nil, err
	}

	return &graphv1.CreateDocumentResponse{
		Document:          document,
		UploadUrl:         "",
		UploadMethod:      "MOCK",
		UploadContentType: req.GetMimeType(),
	}, nil
}

func (s *DocumentService) GetDocument(ctx context.Context, req *graphv1.GetDocumentRequest) (*graphv1.Document, error) {
	document, err := s.repo.GetDocument(ctx, req.GetWorkspaceId(), req.GetDocumentId())
	if err != nil {
		return nil, err
	}
	if document == nil {
		return nil, errors.New("document not found")
	}

	return document, nil
}

func (s *DocumentService) ListDocuments(ctx context.Context, req *graphv1.ListDocumentsRequest) (*graphv1.ListDocumentsResponse, error) {
	documents, err := s.repo.ListDocuments(ctx, req.GetWorkspaceId())
	if err != nil {
		return nil, err
	}

	return &graphv1.ListDocumentsResponse{Documents: documents}, nil
}

func (s *DocumentService) StartProcessing(ctx context.Context, req *graphv1.StartProcessingRequest) (*graphv1.StartProcessingResponse, error) {
	document, err := s.repo.StartProcessing(ctx, req.GetWorkspaceId(), req.GetDocumentId())
	if err != nil {
		return nil, err
	}
	if document == nil {
		return nil, errors.New("document not found")
	}

	return &graphv1.StartProcessingResponse{
		DocumentId: document.DocumentId,
		Status:     document.Status,
		JobId:      "job_" + document.DocumentId,
	}, nil
}
