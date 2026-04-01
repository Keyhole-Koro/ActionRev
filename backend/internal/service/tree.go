package service

import (
	"context"
	"errors"
	"strings"

	graphv1 "github.com/synthify/backend/gen/synthify/graph/v1"
	"github.com/synthify/backend/internal/repository"
)

type WorkspaceTreeService struct {
	workspaces repository.WorkspaceRepository
	documents  repository.DocumentRepository
	trees      repository.TreeRepository
}

func NewWorkspaceTreeService(workspaces repository.WorkspaceRepository, documents repository.DocumentRepository, trees repository.TreeRepository) *WorkspaceTreeService {
	return &WorkspaceTreeService{
		workspaces: workspaces,
		documents:  documents,
		trees:      trees,
	}
}

func (s *WorkspaceTreeService) GetWorkspaceTree(ctx context.Context, req *graphv1.GetWorkspaceTreeRequest) (*graphv1.GetWorkspaceTreeResponse, error) {
	workspaceID := strings.TrimSpace(req.GetWorkspaceId())
	if workspaceID == "" {
		return nil, errors.New("workspace id is required")
	}

	workspace, err := s.workspaces.GetWorkspace(ctx, workspaceID)
	if err != nil {
		return nil, err
	}
	if workspace == nil {
		return nil, errors.New("workspace not found")
	}

	tree, err := s.trees.GetWorkspaceTree(ctx, workspaceID)
	if err != nil {
		return nil, err
	}
	if tree == nil {
		return nil, errors.New("workspace tree not found")
	}

	documents, err := s.documents.ListDocuments(ctx, workspaceID)
	if err != nil {
		return nil, err
	}
	tree.Documents = documents

	return &graphv1.GetWorkspaceTreeResponse{
		Workspace: workspace,
		Tree:      tree,
	}, nil
}

func (s *WorkspaceTreeService) ListPaperNodeChildren(ctx context.Context, req *graphv1.ListPaperNodeChildrenRequest) (*graphv1.ListPaperNodeChildrenResponse, error) {
	workspaceID := strings.TrimSpace(req.GetWorkspaceId())
	parentID := strings.TrimSpace(req.GetParentId())
	if workspaceID == "" {
		return nil, errors.New("workspace id is required")
	}
	if parentID == "" {
		return nil, errors.New("parent id is required")
	}

	nodes, err := s.trees.ListPaperNodeChildren(ctx, workspaceID, parentID)
	if err != nil {
		return nil, err
	}

	return &graphv1.ListPaperNodeChildrenResponse{Nodes: nodes}, nil
}

func (s *WorkspaceTreeService) GetPaperNode(ctx context.Context, req *graphv1.GetPaperNodeRequest) (*graphv1.PaperNode, error) {
	workspaceID := strings.TrimSpace(req.GetWorkspaceId())
	paperNodeID := strings.TrimSpace(req.GetPaperNodeId())
	if workspaceID == "" {
		return nil, errors.New("workspace id is required")
	}
	if paperNodeID == "" {
		return nil, errors.New("paper node id is required")
	}

	node, err := s.trees.GetPaperNode(ctx, workspaceID, paperNodeID)
	if err != nil {
		return nil, err
	}
	if node == nil {
		return nil, errors.New("paper node not found")
	}

	return node, nil
}

func (s *WorkspaceTreeService) CreatePaperNode(ctx context.Context, req *graphv1.CreatePaperNodeRequest) (*graphv1.PaperNode, error) {
	workspaceID := strings.TrimSpace(req.GetWorkspaceId())
	parentID := strings.TrimSpace(req.GetParentId())
	title := strings.TrimSpace(req.GetTitle())
	if workspaceID == "" {
		return nil, errors.New("workspace id is required")
	}
	if parentID == "" {
		return nil, errors.New("parent id is required")
	}
	if title == "" {
		return nil, errors.New("paper node title is required")
	}

	return s.trees.CreatePaperNode(
		ctx,
		workspaceID,
		parentID,
		title,
		strings.TrimSpace(req.GetDescription()),
		strings.TrimSpace(req.GetContent()),
		req.GetCategory(),
		req.GetScope(),
		req.GetSourceDocumentIds(),
	)
}

func (s *WorkspaceTreeService) UpdatePaperNode(ctx context.Context, req *graphv1.UpdatePaperNodeRequest) (*graphv1.PaperNode, error) {
	workspaceID := strings.TrimSpace(req.GetWorkspaceId())
	paperNodeID := strings.TrimSpace(req.GetPaperNodeId())
	if workspaceID == "" {
		return nil, errors.New("workspace id is required")
	}
	if paperNodeID == "" {
		return nil, errors.New("paper node id is required")
	}

	node, err := s.trees.UpdatePaperNode(
		ctx,
		workspaceID,
		paperNodeID,
		strings.TrimSpace(req.GetTitle()),
		strings.TrimSpace(req.GetDescription()),
		strings.TrimSpace(req.GetContent()),
		req.GetStatus(),
	)
	if err != nil {
		return nil, err
	}
	if node == nil {
		return nil, errors.New("paper node not found")
	}

	return node, nil
}

func (s *WorkspaceTreeService) ReorderPaperNode(ctx context.Context, req *graphv1.ReorderPaperNodeRequest) (*graphv1.ReorderPaperNodeResponse, error) {
	workspaceID := strings.TrimSpace(req.GetWorkspaceId())
	paperNodeID := strings.TrimSpace(req.GetPaperNodeId())
	newParentID := strings.TrimSpace(req.GetNewParentId())
	if workspaceID == "" {
		return nil, errors.New("workspace id is required")
	}
	if paperNodeID == "" {
		return nil, errors.New("paper node id is required")
	}
	if newParentID == "" {
		return nil, errors.New("new parent id is required")
	}

	node, siblingIDs, err := s.trees.ReorderPaperNode(ctx, workspaceID, paperNodeID, newParentID, strings.TrimSpace(req.GetInsertBeforeId()))
	if err != nil {
		return nil, err
	}
	if node == nil {
		return nil, errors.New("paper node not found")
	}

	return &graphv1.ReorderPaperNodeResponse{
		Node:       node,
		SiblingIds: siblingIDs,
	}, nil
}
