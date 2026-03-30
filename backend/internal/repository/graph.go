package repository

import (
	"context"

	graphv1 "github.com/synthify/backend/gen/synthify/graph/v1"
)

type GraphRepository interface {
	GetGraph(ctx context.Context, workspaceID string, documentID string) (*graphv1.Graph, error)
	ExpandNeighbors(ctx context.Context, workspaceID string, seedNodeID string, maxDepth uint32, limitPerHop uint32, edgeTypeFilters []graphv1.EdgeType) (*graphv1.Graph, error)
}
