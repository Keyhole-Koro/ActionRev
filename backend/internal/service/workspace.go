package service

import (
	"context"
	"errors"
	"strings"

	graphv1 "github.com/synthify/backend/gen/synthify/graph/v1"
	"github.com/synthify/backend/internal/repository"
)

type WorkspaceService struct {
	repo repository.WorkspaceRepository
}

func NewWorkspaceService(repo repository.WorkspaceRepository) *WorkspaceService {
	return &WorkspaceService{repo: repo}
}

func (s *WorkspaceService) CreateWorkspace(ctx context.Context, req *graphv1.CreateWorkspaceRequest) (*graphv1.Workspace, error) {
	name := strings.TrimSpace(req.GetName())
	if name == "" {
		return nil, errors.New("workspace name is required")
	}

	return s.repo.CreateWorkspace(ctx, name)
}

func (s *WorkspaceService) GetWorkspace(ctx context.Context, req *graphv1.GetWorkspaceRequest) (*graphv1.Workspace, error) {
	workspace, err := s.repo.GetWorkspace(ctx, req.GetWorkspaceId())
	if err != nil {
		return nil, err
	}
	if workspace == nil {
		return nil, errors.New("workspace not found")
	}

	return workspace, nil
}

func (s *WorkspaceService) UpdateWorkspace(ctx context.Context, req *graphv1.UpdateWorkspaceRequest) (*graphv1.Workspace, error) {
	name := strings.TrimSpace(req.GetName())
	if req.GetWorkspaceId() == "" {
		return nil, errors.New("workspace id is required")
	}
	if name == "" {
		return nil, errors.New("workspace name is required")
	}

	workspace, err := s.repo.UpdateWorkspace(ctx, req.GetWorkspaceId(), name)
	if err != nil {
		return nil, err
	}
	if workspace == nil {
		return nil, errors.New("workspace not found")
	}

	return workspace, nil
}

func (s *WorkspaceService) ListWorkspaces(ctx context.Context, _ *graphv1.ListWorkspacesRequest) (*graphv1.ListWorkspacesResponse, error) {
	workspaces, err := s.repo.ListWorkspaces(ctx)
	if err != nil {
		return nil, err
	}

	return &graphv1.ListWorkspacesResponse{
		Workspaces: workspaces,
	}, nil
}
