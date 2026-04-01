package service

import (
	"context"
	"errors"
	"strings"

	graphv1 "github.com/synthify/backend/gen/synthify/graph/v1"
	"github.com/synthify/backend/internal/repository"
)

type PaperNoteService struct {
	trees repository.TreeRepository
}

func NewPaperNoteService(trees repository.TreeRepository) *PaperNoteService {
	return &PaperNoteService{trees: trees}
}

func (s *PaperNoteService) ListNodeNotes(ctx context.Context, req *graphv1.ListNodeNotesRequest) (*graphv1.ListNodeNotesResponse, error) {
	workspaceID := strings.TrimSpace(req.GetWorkspaceId())
	paperNodeID := strings.TrimSpace(req.GetPaperNodeId())
	if workspaceID == "" {
		return nil, errors.New("workspace id is required")
	}
	if paperNodeID == "" {
		return nil, errors.New("paper node id is required")
	}

	notes, err := s.trees.ListNodeNotes(ctx, workspaceID, paperNodeID)
	if err != nil {
		return nil, err
	}

	return &graphv1.ListNodeNotesResponse{Notes: notes}, nil
}

func (s *PaperNoteService) CreateNodeNote(ctx context.Context, req *graphv1.CreateNodeNoteRequest) (*graphv1.PaperNote, error) {
	workspaceID := strings.TrimSpace(req.GetWorkspaceId())
	paperNodeID := strings.TrimSpace(req.GetPaperNodeId())
	title := strings.TrimSpace(req.GetTitle())
	body := strings.TrimSpace(req.GetBody())
	if workspaceID == "" {
		return nil, errors.New("workspace id is required")
	}
	if paperNodeID == "" {
		return nil, errors.New("paper node id is required")
	}
	if title == "" {
		return nil, errors.New("note title is required")
	}
	if body == "" {
		return nil, errors.New("note body is required")
	}

	return s.trees.CreateNodeNote(ctx, workspaceID, paperNodeID, req.GetKind(), title, body, req.GetPriority())
}
