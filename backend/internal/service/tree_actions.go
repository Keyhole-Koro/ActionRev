package service

import (
	"context"
	"errors"
	"strings"

	graphv1 "github.com/synthify/backend/gen/synthify/graph/v1"
	"github.com/synthify/backend/internal/repository"
)

type ActionRequestService struct {
	trees repository.TreeRepository
}

func NewActionRequestService(trees repository.TreeRepository) *ActionRequestService {
	return &ActionRequestService{trees: trees}
}

func (s *ActionRequestService) ListNodeActionRequests(ctx context.Context, req *graphv1.ListNodeActionRequestsRequest) (*graphv1.ListNodeActionRequestsResponse, error) {
	workspaceID := strings.TrimSpace(req.GetWorkspaceId())
	paperNodeID := strings.TrimSpace(req.GetPaperNodeId())
	if workspaceID == "" {
		return nil, errors.New("workspace id is required")
	}
	if paperNodeID == "" {
		return nil, errors.New("paper node id is required")
	}

	items, err := s.trees.ListNodeActionRequests(ctx, workspaceID, paperNodeID)
	if err != nil {
		return nil, err
	}

	return &graphv1.ListNodeActionRequestsResponse{ActionRequests: items}, nil
}

func (s *ActionRequestService) ResolveActionRequest(ctx context.Context, req *graphv1.ResolveActionRequestRequest) (*graphv1.ActionRequest, error) {
	workspaceID := strings.TrimSpace(req.GetWorkspaceId())
	actionRequestID := strings.TrimSpace(req.GetActionRequestId())
	if workspaceID == "" {
		return nil, errors.New("workspace id is required")
	}
	if actionRequestID == "" {
		return nil, errors.New("action request id is required")
	}

	item, err := s.trees.ResolveActionRequest(ctx, workspaceID, actionRequestID)
	if err != nil {
		return nil, err
	}
	if item == nil {
		return nil, errors.New("action request not found")
	}

	return item, nil
}

func (s *ActionRequestService) DismissActionRequest(ctx context.Context, req *graphv1.DismissActionRequestRequest) (*graphv1.ActionRequest, error) {
	workspaceID := strings.TrimSpace(req.GetWorkspaceId())
	actionRequestID := strings.TrimSpace(req.GetActionRequestId())
	if workspaceID == "" {
		return nil, errors.New("workspace id is required")
	}
	if actionRequestID == "" {
		return nil, errors.New("action request id is required")
	}

	item, err := s.trees.DismissActionRequest(ctx, workspaceID, actionRequestID)
	if err != nil {
		return nil, err
	}
	if item == nil {
		return nil, errors.New("action request not found")
	}

	return item, nil
}
