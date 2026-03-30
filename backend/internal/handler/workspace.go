package handler

import (
	"context"
	"errors"
	"strings"

	connect "connectrpc.com/connect"
	graphv1 "github.com/synthify/backend/gen/synthify/graph/v1"
	"github.com/synthify/backend/internal/service"
)

type WorkspaceHandler struct {
	service *service.WorkspaceService
}

func NewWorkspaceHandler(service *service.WorkspaceService) *WorkspaceHandler {
	return &WorkspaceHandler{service: service}
}

func (h *WorkspaceHandler) CreateWorkspace(
	ctx context.Context,
	req *connect.Request[graphv1.CreateWorkspaceRequest],
) (*connect.Response[graphv1.Workspace], error) {
	workspace, err := h.service.CreateWorkspace(ctx, req.Msg)
	if err != nil {
		return nil, toWorkspaceConnectError(err)
	}

	return connect.NewResponse(workspace), nil
}

func (h *WorkspaceHandler) GetWorkspace(
	ctx context.Context,
	req *connect.Request[graphv1.GetWorkspaceRequest],
) (*connect.Response[graphv1.Workspace], error) {
	workspace, err := h.service.GetWorkspace(ctx, req.Msg)
	if err != nil {
		return nil, toWorkspaceConnectError(err)
	}

	return connect.NewResponse(workspace), nil
}

func (h *WorkspaceHandler) UpdateWorkspace(
	ctx context.Context,
	req *connect.Request[graphv1.UpdateWorkspaceRequest],
) (*connect.Response[graphv1.Workspace], error) {
	workspace, err := h.service.UpdateWorkspace(ctx, req.Msg)
	if err != nil {
		return nil, toWorkspaceConnectError(err)
	}

	return connect.NewResponse(workspace), nil
}

func (h *WorkspaceHandler) ListWorkspaces(
	ctx context.Context,
	req *connect.Request[graphv1.ListWorkspacesRequest],
) (*connect.Response[graphv1.ListWorkspacesResponse], error) {
	response, err := h.service.ListWorkspaces(ctx, req.Msg)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return connect.NewResponse(response), nil
}

func (h *WorkspaceHandler) AddWorkspaceMember(
	_ context.Context,
	_ *connect.Request[graphv1.AddWorkspaceMemberRequest],
) (*connect.Response[graphv1.AddWorkspaceMemberResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("AddWorkspaceMember is not implemented"))
}

func toWorkspaceConnectError(err error) error {
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
