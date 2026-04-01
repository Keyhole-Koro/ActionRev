package handler

import (
	"context"
	"strings"

	connect "connectrpc.com/connect"
	graphv1 "github.com/synthify/backend/gen/synthify/graph/v1"
	"github.com/synthify/backend/internal/service"
)

type WorkspaceTreeHandler struct {
	service *service.WorkspaceTreeService
}

func NewWorkspaceTreeHandler(service *service.WorkspaceTreeService) *WorkspaceTreeHandler {
	return &WorkspaceTreeHandler{service: service}
}

func (h *WorkspaceTreeHandler) GetWorkspaceTree(ctx context.Context, req *connect.Request[graphv1.GetWorkspaceTreeRequest]) (*connect.Response[graphv1.GetWorkspaceTreeResponse], error) {
	res, err := h.service.GetWorkspaceTree(ctx, req.Msg)
	if err != nil {
		return nil, toTreeConnectError(err)
	}

	return connect.NewResponse(res), nil
}

func (h *WorkspaceTreeHandler) ListPaperNodeChildren(ctx context.Context, req *connect.Request[graphv1.ListPaperNodeChildrenRequest]) (*connect.Response[graphv1.ListPaperNodeChildrenResponse], error) {
	res, err := h.service.ListPaperNodeChildren(ctx, req.Msg)
	if err != nil {
		return nil, toTreeConnectError(err)
	}

	return connect.NewResponse(res), nil
}

func (h *WorkspaceTreeHandler) GetPaperNode(ctx context.Context, req *connect.Request[graphv1.GetPaperNodeRequest]) (*connect.Response[graphv1.PaperNode], error) {
	res, err := h.service.GetPaperNode(ctx, req.Msg)
	if err != nil {
		return nil, toTreeConnectError(err)
	}

	return connect.NewResponse(res), nil
}

func (h *WorkspaceTreeHandler) CreatePaperNode(ctx context.Context, req *connect.Request[graphv1.CreatePaperNodeRequest]) (*connect.Response[graphv1.PaperNode], error) {
	res, err := h.service.CreatePaperNode(ctx, req.Msg)
	if err != nil {
		return nil, toTreeConnectError(err)
	}

	return connect.NewResponse(res), nil
}

func (h *WorkspaceTreeHandler) UpdatePaperNode(ctx context.Context, req *connect.Request[graphv1.UpdatePaperNodeRequest]) (*connect.Response[graphv1.PaperNode], error) {
	res, err := h.service.UpdatePaperNode(ctx, req.Msg)
	if err != nil {
		return nil, toTreeConnectError(err)
	}

	return connect.NewResponse(res), nil
}

func (h *WorkspaceTreeHandler) ReorderPaperNode(ctx context.Context, req *connect.Request[graphv1.ReorderPaperNodeRequest]) (*connect.Response[graphv1.ReorderPaperNodeResponse], error) {
	res, err := h.service.ReorderPaperNode(ctx, req.Msg)
	if err != nil {
		return nil, toTreeConnectError(err)
	}

	return connect.NewResponse(res), nil
}

type PaperNoteHandler struct {
	service *service.PaperNoteService
}

func NewPaperNoteHandler(service *service.PaperNoteService) *PaperNoteHandler {
	return &PaperNoteHandler{service: service}
}

func (h *PaperNoteHandler) ListNodeNotes(ctx context.Context, req *connect.Request[graphv1.ListNodeNotesRequest]) (*connect.Response[graphv1.ListNodeNotesResponse], error) {
	res, err := h.service.ListNodeNotes(ctx, req.Msg)
	if err != nil {
		return nil, toTreeConnectError(err)
	}

	return connect.NewResponse(res), nil
}

func (h *PaperNoteHandler) CreateNodeNote(ctx context.Context, req *connect.Request[graphv1.CreateNodeNoteRequest]) (*connect.Response[graphv1.PaperNote], error) {
	res, err := h.service.CreateNodeNote(ctx, req.Msg)
	if err != nil {
		return nil, toTreeConnectError(err)
	}

	return connect.NewResponse(res), nil
}

type ActionRequestHandler struct {
	service *service.ActionRequestService
}

func NewActionRequestHandler(service *service.ActionRequestService) *ActionRequestHandler {
	return &ActionRequestHandler{service: service}
}

func (h *ActionRequestHandler) ListNodeActionRequests(ctx context.Context, req *connect.Request[graphv1.ListNodeActionRequestsRequest]) (*connect.Response[graphv1.ListNodeActionRequestsResponse], error) {
	res, err := h.service.ListNodeActionRequests(ctx, req.Msg)
	if err != nil {
		return nil, toTreeConnectError(err)
	}

	return connect.NewResponse(res), nil
}

func (h *ActionRequestHandler) ResolveActionRequest(ctx context.Context, req *connect.Request[graphv1.ResolveActionRequestRequest]) (*connect.Response[graphv1.ActionRequest], error) {
	res, err := h.service.ResolveActionRequest(ctx, req.Msg)
	if err != nil {
		return nil, toTreeConnectError(err)
	}

	return connect.NewResponse(res), nil
}

func (h *ActionRequestHandler) DismissActionRequest(ctx context.Context, req *connect.Request[graphv1.DismissActionRequestRequest]) (*connect.Response[graphv1.ActionRequest], error) {
	res, err := h.service.DismissActionRequest(ctx, req.Msg)
	if err != nil {
		return nil, toTreeConnectError(err)
	}

	return connect.NewResponse(res), nil
}

func toTreeConnectError(err error) error {
	if err == nil {
		return nil
	}

	switch {
	case strings.Contains(err.Error(), "not found"):
		return connect.NewError(connect.CodeNotFound, err)
	case strings.Contains(err.Error(), "required"):
		return connect.NewError(connect.CodeInvalidArgument, err)
	default:
		return connect.NewError(connect.CodeInternal, err)
	}
}
