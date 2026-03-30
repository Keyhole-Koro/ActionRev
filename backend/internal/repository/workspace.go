package repository

import (
	"context"

	graphv1 "github.com/synthify/backend/gen/synthify/graph/v1"
)

type WorkspaceRepository interface {
	CreateWorkspace(ctx context.Context, name string) (*graphv1.Workspace, error)
	GetWorkspace(ctx context.Context, workspaceID string) (*graphv1.Workspace, error)
	UpdateWorkspace(ctx context.Context, workspaceID string, name string) (*graphv1.Workspace, error)
	ListWorkspaces(ctx context.Context) ([]*graphv1.Workspace, error)
}
