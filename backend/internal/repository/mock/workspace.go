package mock

import (
	"context"
	"crypto/rand"
	"fmt"
	"sync"
	"time"

	graphv1 "github.com/synthify/backend/gen/synthify/graph/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type WorkspaceRepository struct {
	mu         sync.RWMutex
	workspaces map[string]*graphv1.Workspace
}

func NewWorkspaceRepository() *WorkspaceRepository {
	now := timestamppb.New(time.Now())

	return &WorkspaceRepository{
		workspaces: map[string]*graphv1.Workspace{
			"00000000-0000-4000-8000-000000000001": {
				WorkspaceId: "00000000-0000-4000-8000-000000000001",
				Name:        "Growth Strategy Review",
				CreatedAt:   now,
			},
		},
	}
}

func (r *WorkspaceRepository) CreateWorkspace(_ context.Context, name string) (*graphv1.Workspace, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	workspaceID, err := newWorkspaceID()
	if err != nil {
		return nil, err
	}
	workspace := &graphv1.Workspace{
		WorkspaceId: workspaceID,
		Name:        name,
		CreatedAt:   timestamppb.New(time.Now()),
	}
	r.workspaces[workspaceID] = workspace

	return cloneWorkspace(workspace), nil
}

func (r *WorkspaceRepository) GetWorkspace(_ context.Context, workspaceID string) (*graphv1.Workspace, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	workspace, ok := r.workspaces[workspaceID]
	if !ok {
		return nil, nil
	}

	return cloneWorkspace(workspace), nil
}

func (r *WorkspaceRepository) UpdateWorkspace(_ context.Context, workspaceID string, name string) (*graphv1.Workspace, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	workspace, ok := r.workspaces[workspaceID]
	if !ok {
		return nil, nil
	}

	workspace.Name = name
	return cloneWorkspace(workspace), nil
}

func (r *WorkspaceRepository) ListWorkspaces(_ context.Context) ([]*graphv1.Workspace, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	workspaces := make([]*graphv1.Workspace, 0, len(r.workspaces))
	for _, workspace := range r.workspaces {
		workspaces = append(workspaces, cloneWorkspace(workspace))
	}

	return workspaces, nil
}

func cloneWorkspace(workspace *graphv1.Workspace) *graphv1.Workspace {
	if workspace == nil {
		return nil
	}

	return &graphv1.Workspace{
		WorkspaceId: workspace.WorkspaceId,
		Name:        workspace.Name,
		CreatedAt:   workspace.CreatedAt,
	}
}

func newWorkspaceID() (string, error) {
	var value [16]byte
	if _, err := rand.Read(value[:]); err != nil {
		return "", err
	}

	value[6] = (value[6] & 0x0f) | 0x40
	value[8] = (value[8] & 0x3f) | 0x80

	return fmt.Sprintf(
		"%08x-%04x-%04x-%04x-%012x",
		value[0:4],
		value[4:6],
		value[6:8],
		value[8:10],
		value[10:16],
	), nil
}
