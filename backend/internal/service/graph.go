package service

import (
	"context"
	"errors"

	graphv1 "github.com/synthify/backend/gen/synthify/graph/v1"
	"github.com/synthify/backend/internal/repository"
)

type GraphService struct {
	repo repository.GraphRepository
}

func NewGraphService(repo repository.GraphRepository) *GraphService {
	return &GraphService{repo: repo}
}

func (s *GraphService) GetGraph(ctx context.Context, req *graphv1.GetGraphRequest) (*graphv1.GetGraphResponse, error) {
	graph, err := s.repo.GetGraph(ctx, req.WorkspaceId, req.DocumentId)
	if err != nil {
		return nil, err
	}
	if graph == nil {
		return nil, errors.New("graph not found")
	}

	return &graphv1.GetGraphResponse{
		DocumentId: graph.DocumentId,
		Graph:      graph,
	}, nil
}

func (s *GraphService) ExpandNeighbors(ctx context.Context, req *graphv1.ExpandNeighborsRequest) (*graphv1.ExpandNeighborsResponse, error) {
	graph, err := s.repo.ExpandNeighbors(ctx, req.WorkspaceId, req.SeedNodeId, req.MaxDepth, req.LimitPerHop, req.EdgeTypeFilters)
	if err != nil {
		return nil, err
	}

	return &graphv1.ExpandNeighborsResponse{
		Graph:      graph,
		SeedNodeId: req.SeedNodeId,
		Depth:      req.MaxDepth,
	}, nil
}
