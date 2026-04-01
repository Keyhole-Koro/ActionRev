package repository

import (
	"context"

	graphv1 "github.com/synthify/backend/gen/synthify/graph/v1"
)

type TreeRepository interface {
	GetWorkspaceTree(ctx context.Context, workspaceID string) (*graphv1.WorkspaceTree, error)
	ListPaperNodeChildren(ctx context.Context, workspaceID string, parentID string) ([]*graphv1.PaperNode, error)
	GetPaperNode(ctx context.Context, workspaceID string, paperNodeID string) (*graphv1.PaperNode, error)
	CreatePaperNode(ctx context.Context, workspaceID string, parentID string, title string, description string, content string, category graphv1.PaperNodeCategory, scope graphv1.PaperNodeScope, sourceDocumentIDs []string) (*graphv1.PaperNode, error)
	UpdatePaperNode(ctx context.Context, workspaceID string, paperNodeID string, title string, description string, content string, status graphv1.PaperNodeStatus) (*graphv1.PaperNode, error)
	ReorderPaperNode(ctx context.Context, workspaceID string, paperNodeID string, newParentID string, insertBeforeID string) (*graphv1.PaperNode, []string, error)
	ListNodeNotes(ctx context.Context, workspaceID string, paperNodeID string) ([]*graphv1.PaperNote, error)
	CreateNodeNote(ctx context.Context, workspaceID string, paperNodeID string, kind graphv1.PaperNoteKind, title string, body string, priority graphv1.NotePriority) (*graphv1.PaperNote, error)
	ListNodeActionRequests(ctx context.Context, workspaceID string, paperNodeID string) ([]*graphv1.ActionRequest, error)
	ResolveActionRequest(ctx context.Context, workspaceID string, actionRequestID string) (*graphv1.ActionRequest, error)
	DismissActionRequest(ctx context.Context, workspaceID string, actionRequestID string) (*graphv1.ActionRequest, error)
}
