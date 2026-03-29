package mock

import (
	"context"

	graphv1 "github.com/synthify/backend/gen/synthify/graph/v1"
)

type GraphRepository struct{}

func NewGraphRepository() *GraphRepository {
	return &GraphRepository{}
}

func (r *GraphRepository) GetGraph(_ context.Context, workspaceID string, documentID string) (*graphv1.Graph, error) {
	if workspaceID == "" {
		workspaceID = "ws_demo"
	}
	if documentID == "" {
		documentID = "doc_demo"
	}

	return &graphv1.Graph{
		WorkspaceId:   workspaceID,
		DocumentId:    documentID,
		CrossDocument: false,
		Nodes: []*graphv1.Node{
			{
				Id:              "nd_sales_strategy",
				DocumentId:      documentID,
				CanonicalNodeId: "cn_sales_strategy",
				Scope:           graphv1.GraphProjectionScope_GRAPH_PROJECTION_SCOPE_DOCUMENT,
				Label:           "販売戦略",
				Level:           1,
				Category:        graphv1.NodeCategory_NODE_CATEGORY_CONCEPT,
				Description:     "売上拡大のための上位方針",
				SourceChunkIds:  []string{"chk_001"},
			},
			{
				Id:              "nd_sns_strategy",
				DocumentId:      documentID,
				CanonicalNodeId: "cn_sns_strategy",
				Scope:           graphv1.GraphProjectionScope_GRAPH_PROJECTION_SCOPE_DOCUMENT,
				Label:           "SNS施策",
				Level:           2,
				Category:        graphv1.NodeCategory_NODE_CATEGORY_ACTION,
				Description:     "SNS経由の集客施策",
				SourceChunkIds:  []string{"chk_002"},
			},
			{
				Id:              "nd_cv_rate",
				DocumentId:      documentID,
				CanonicalNodeId: "cn_cv_rate",
				Scope:           graphv1.GraphProjectionScope_GRAPH_PROJECTION_SCOPE_DOCUMENT,
				Label:           "CV率3.2%",
				Level:           3,
				Category:        graphv1.NodeCategory_NODE_CATEGORY_METRIC,
				Description:     "主要KPIとして記載されたCV率",
				SourceChunkIds:  []string{"chk_003"},
			},
		},
		Edges: []*graphv1.Edge{
			{
				Id:             "ed_sales_to_sns",
				DocumentId:     documentID,
				Scope:          graphv1.GraphProjectionScope_GRAPH_PROJECTION_SCOPE_DOCUMENT,
				Source:         "nd_sales_strategy",
				Target:         "nd_sns_strategy",
				Type:           graphv1.EdgeType_EDGE_TYPE_HIERARCHICAL,
				SourceChunkIds: []string{"chk_002"},
			},
			{
				Id:             "ed_sns_to_cv",
				DocumentId:     documentID,
				Scope:          graphv1.GraphProjectionScope_GRAPH_PROJECTION_SCOPE_DOCUMENT,
				Source:         "nd_sns_strategy",
				Target:         "nd_cv_rate",
				Type:           graphv1.EdgeType_EDGE_TYPE_MEASURED_BY,
				SourceChunkIds: []string{"chk_003"},
			},
		},
	}, nil
}
