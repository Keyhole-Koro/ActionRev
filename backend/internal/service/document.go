package service

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/url"
	"strings"

	graphv1 "github.com/synthify/backend/gen/synthify/graph/v1"
	"github.com/synthify/backend/internal/repository"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type DocumentService struct {
	repo          repository.DocumentRepository
	publicBaseURL string
}

func NewDocumentService(repo repository.DocumentRepository, publicBaseURL string) *DocumentService {
	return &DocumentService{
		repo:          repo,
		publicBaseURL: strings.TrimRight(publicBaseURL, "/"),
	}
}

func (s *DocumentService) CreateDocument(ctx context.Context, req *graphv1.CreateDocumentRequest) (*graphv1.CreateDocumentResponse, error) {
	if err := validateDocumentUploadRequest(req.GetWorkspaceId(), req.GetFilename(), req.GetMimeType(), req.GetFileSize()); err != nil {
		return nil, err
	}

	document, err := s.repo.CreateDocument(ctx, req.GetWorkspaceId(), req.GetFilename(), req.GetMimeType(), req.GetFileSize())
	if err != nil {
		return nil, err
	}

	uploadTarget, err := s.repo.CreateUploadTarget(ctx, req.GetWorkspaceId(), req.GetFilename(), req.GetMimeType(), req.GetFileSize())
	if err != nil {
		return nil, err
	}

	return &graphv1.CreateDocumentResponse{
		Document:          document,
		UploadUrl:         s.buildMockUploadURL(uploadTarget.Token),
		UploadMethod:      "PUT",
		UploadContentType: req.GetMimeType(),
	}, nil
}

func (s *DocumentService) GetUploadURL(ctx context.Context, req *graphv1.GetUploadUrlRequest) (*graphv1.GetUploadUrlResponse, error) {
	if err := validateDocumentUploadRequest(req.GetWorkspaceId(), req.GetFilename(), req.GetMimeType(), req.GetFileSize()); err != nil {
		return nil, err
	}

	target, err := s.repo.CreateUploadTarget(ctx, req.GetWorkspaceId(), req.GetFilename(), req.GetMimeType(), req.GetFileSize())
	if err != nil {
		return nil, err
	}

	return &graphv1.GetUploadUrlResponse{
		UploadUrl:   s.buildMockUploadURL(target.Token),
		UploadToken: target.Token,
		ExpiresAt:   timestamppb.New(target.ExpiresAt),
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

func (s *DocumentService) ConsumeMockUpload(ctx context.Context, token string, contentType string, contentLength int64, body io.Reader) error {
	if strings.TrimSpace(token) == "" {
		return errors.New("upload token is required")
	}

	target, err := s.repo.ConsumeUploadTarget(ctx, token)
	if err != nil {
		return err
	}
	if target == nil {
		return errors.New("upload token not found")
	}

	if contentType != "" && target.MimeType != "" && !strings.EqualFold(contentType, target.MimeType) {
		return fmt.Errorf("content type mismatch: expected %s", target.MimeType)
	}
	if contentLength > 0 && target.FileSize > 0 && contentLength != target.FileSize {
		return fmt.Errorf("content length mismatch: expected %d bytes", target.FileSize)
	}
	if _, err := io.Copy(io.Discard, body); err != nil {
		return err
	}

	return nil
}

func (s *DocumentService) buildMockUploadURL(token string) string {
	baseURL := s.publicBaseURL
	if baseURL == "" {
		baseURL = "http://localhost:8080"
	}

	return fmt.Sprintf("%s/mock/uploads/%s", strings.TrimRight(baseURL, "/"), url.PathEscape(token))
}

func validateDocumentUploadRequest(workspaceID string, filename string, mimeType string, fileSize int64) error {
	if strings.TrimSpace(workspaceID) == "" {
		return errors.New("workspace id is required")
	}
	if strings.TrimSpace(filename) == "" {
		return errors.New("filename is required")
	}
	if strings.TrimSpace(mimeType) == "" {
		return errors.New("mime type is required")
	}
	if fileSize <= 0 {
		return errors.New("file size must be positive")
	}

	return nil
}
