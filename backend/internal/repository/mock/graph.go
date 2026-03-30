package mock

import (
	"context"
	"strings"

	graphv1 "github.com/synthify/backend/gen/synthify/graph/v1"
)

type GraphRepository struct {
	documents *DocumentRepository
}

func NewGraphRepository(documents *DocumentRepository) *GraphRepository {
	return &GraphRepository{documents: documents}
}

func (r *GraphRepository) GetGraph(_ context.Context, workspaceID string, documentID string) (*graphv1.Graph, error) {
	if workspaceID == "" {
		workspaceID = "00000000-0000-4000-8000-000000000001"
	}

	documents, err := r.documents.ListDocuments(context.Background(), workspaceID)
	if err != nil {
		return nil, err
	}

	activeDocuments := documents
	if documentID != "" {
		activeDocuments = make([]*graphv1.Document, 0, 1)
		for _, document := range documents {
			if document.DocumentId == documentID {
				activeDocuments = append(activeDocuments, document)
				break
			}
		}
	}

	nodes := []*graphv1.Node{
		{
			Id:              "cn_workspace_strategy",
			DocumentId:      "",
			CanonicalNodeId: "cn_workspace_strategy",
			Scope:           graphv1.GraphProjectionScope_GRAPH_PROJECTION_SCOPE_CANONICAL,
			Label:           "Workspace Strategy",
			Level:           1,
			Category:        graphv1.NodeCategory_NODE_CATEGORY_CONCEPT,
			Description:     "Workspace 全体に統合された上位概念",
			SourceChunkIds:  []string{},
		},
		{
			Id:              "cn_workspace_evidence",
			DocumentId:      "",
			CanonicalNodeId: "cn_workspace_evidence",
			Scope:           graphv1.GraphProjectionScope_GRAPH_PROJECTION_SCOPE_CANONICAL,
			Label:           "Collected Evidence",
			Level:           2,
			Category:        graphv1.NodeCategory_NODE_CATEGORY_EVIDENCE,
			Description:     "複数 document から集まった根拠",
			SourceChunkIds:  []string{},
		},
		{
			Id:              "cn_workspace_metrics",
			DocumentId:      "",
			CanonicalNodeId: "cn_workspace_metrics",
			Scope:           graphv1.GraphProjectionScope_GRAPH_PROJECTION_SCOPE_CANONICAL,
			Label:           "Shared Metrics",
			Level:           2,
			Category:        graphv1.NodeCategory_NODE_CATEGORY_METRIC,
			Description:     "Workspace 横断で参照される KPI / metrics",
			SourceChunkIds:  []string{},
		},
	}
	edges := []*graphv1.Edge{
		{
			Id:             "ed_strategy_to_evidence",
			DocumentId:     "",
			Scope:          graphv1.GraphProjectionScope_GRAPH_PROJECTION_SCOPE_CANONICAL,
			Source:         "cn_workspace_strategy",
			Target:         "cn_workspace_evidence",
			Type:           graphv1.EdgeType_EDGE_TYPE_SUPPORTS,
			SourceChunkIds: []string{},
		},
		{
			Id:             "ed_strategy_to_metrics",
			DocumentId:     "",
			Scope:          graphv1.GraphProjectionScope_GRAPH_PROJECTION_SCOPE_CANONICAL,
			Source:         "cn_workspace_strategy",
			Target:         "cn_workspace_metrics",
			Type:           graphv1.EdgeType_EDGE_TYPE_MEASURED_BY,
			SourceChunkIds: []string{},
		},
	}

	for index, document := range activeDocuments {
		docNodeID := "doc_node_" + document.DocumentId
		slug := slugFilename(document.Filename)
		chunkID := "chk_" + document.DocumentId

		nodes = append(nodes, &graphv1.Node{
			Id:              docNodeID,
			DocumentId:      document.DocumentId,
			CanonicalNodeId: "cn_workspace_strategy",
			Scope:           graphv1.GraphProjectionScope_GRAPH_PROJECTION_SCOPE_DOCUMENT,
			Label:           document.Filename,
			Level:           3,
			Category:        graphv1.NodeCategory_NODE_CATEGORY_EVIDENCE,
			Description:     "document source: " + document.Filename,
			SourceChunkIds:  []string{chunkID},
		})
		nodes = append(nodes, &graphv1.Node{
			Id:              "metric_" + slug,
			DocumentId:      document.DocumentId,
			CanonicalNodeId: "cn_workspace_metrics",
			Scope:           graphv1.GraphProjectionScope_GRAPH_PROJECTION_SCOPE_DOCUMENT,
			Label:           "Signal " + strings.ToUpper(string(rune('A'+index))),
			Level:           3,
			Category:        graphv1.NodeCategory_NODE_CATEGORY_METRIC,
			Description:     "document-specific metric extracted from " + document.Filename,
			SourceChunkIds:  []string{chunkID},
		})
		edges = append(edges, &graphv1.Edge{
			Id:             "ed_evidence_" + document.DocumentId,
			DocumentId:     document.DocumentId,
			Scope:          graphv1.GraphProjectionScope_GRAPH_PROJECTION_SCOPE_DOCUMENT,
			Source:         docNodeID,
			Target:         "cn_workspace_evidence",
			Type:           graphv1.EdgeType_EDGE_TYPE_SUPPORTS,
			SourceChunkIds: []string{chunkID},
		})
		edges = append(edges, &graphv1.Edge{
			Id:             "ed_metric_" + document.DocumentId,
			DocumentId:     document.DocumentId,
			Scope:          graphv1.GraphProjectionScope_GRAPH_PROJECTION_SCOPE_DOCUMENT,
			Source:         docNodeID,
			Target:         "metric_" + slug,
			Type:           graphv1.EdgeType_EDGE_TYPE_MEASURED_BY,
			SourceChunkIds: []string{chunkID},
		})
	}

	return &graphv1.Graph{
		WorkspaceId:   workspaceID,
		DocumentId:    "",
		CrossDocument: true,
		Nodes:         nodes,
		Edges:         edges,
	}, nil
}

func slugFilename(filename string) string {
	slug := strings.ToLower(filename)
	replacer := strings.NewReplacer(" ", "-", ".", "-", "_", "-")
	return replacer.Replace(slug)
}
