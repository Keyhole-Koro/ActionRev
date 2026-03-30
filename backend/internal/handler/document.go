package handler

import (
	"context"
	"net/http"
	"strings"

	connect "connectrpc.com/connect"
	graphv1 "github.com/synthify/backend/gen/synthify/graph/v1"
	"github.com/synthify/backend/internal/service"
)

type DocumentHandler struct {
	service *service.DocumentService
}

func NewDocumentHandler(service *service.DocumentService) *DocumentHandler {
	return &DocumentHandler{service: service}
}

func (h *DocumentHandler) CreateDocument(ctx context.Context, req *connect.Request[graphv1.CreateDocumentRequest]) (*connect.Response[graphv1.CreateDocumentResponse], error) {
	response, err := h.service.CreateDocument(ctx, req.Msg)
	if err != nil {
		return nil, toDocumentConnectError(err)
	}

	return connect.NewResponse(response), nil
}

func (h *DocumentHandler) GetUploadUrl(ctx context.Context, req *connect.Request[graphv1.GetUploadUrlRequest]) (*connect.Response[graphv1.GetUploadUrlResponse], error) {
	response, err := h.service.GetUploadURL(ctx, req.Msg)
	if err != nil {
		return nil, toDocumentConnectError(err)
	}

	return connect.NewResponse(response), nil
}

func (h *DocumentHandler) GetDocument(ctx context.Context, req *connect.Request[graphv1.GetDocumentRequest]) (*connect.Response[graphv1.Document], error) {
	document, err := h.service.GetDocument(ctx, req.Msg)
	if err != nil {
		return nil, toDocumentConnectError(err)
	}

	return connect.NewResponse(document), nil
}

func (h *DocumentHandler) ListDocuments(ctx context.Context, req *connect.Request[graphv1.ListDocumentsRequest]) (*connect.Response[graphv1.ListDocumentsResponse], error) {
	response, err := h.service.ListDocuments(ctx, req.Msg)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return connect.NewResponse(response), nil
}

func (h *DocumentHandler) StartProcessing(ctx context.Context, req *connect.Request[graphv1.StartProcessingRequest]) (*connect.Response[graphv1.StartProcessingResponse], error) {
	response, err := h.service.StartProcessing(ctx, req.Msg)
	if err != nil {
		return nil, toDocumentConnectError(err)
	}

	return connect.NewResponse(response), nil
}

func toDocumentConnectError(err error) error {
	if err == nil {
		return nil
	}

	switch {
	case strings.Contains(err.Error(), "not found"):
		return connect.NewError(connect.CodeNotFound, err)
	case strings.Contains(err.Error(), "required"),
		strings.Contains(err.Error(), "must be positive"),
		strings.Contains(err.Error(), "mismatch"),
		strings.Contains(err.Error(), "expired"):
		return connect.NewError(connect.CodeInvalidArgument, err)
	default:
		return connect.NewError(connect.CodeInternal, err)
	}
}

type MockUploadHandler struct {
	service *service.DocumentService
}

func NewMockUploadHandler(service *service.DocumentService) *MockUploadHandler {
	return &MockUploadHandler{service: service}
}

func (h *MockUploadHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		w.Header().Set("Allow", http.MethodPut)
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	token := strings.TrimPrefix(r.URL.Path, "/mock/uploads/")
	if token == "" || token == r.URL.Path {
		http.Error(w, "upload token is required", http.StatusBadRequest)
		return
	}

	if err := h.service.ConsumeMockUpload(r.Context(), token, r.Header.Get("Content-Type"), r.ContentLength, r.Body); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
