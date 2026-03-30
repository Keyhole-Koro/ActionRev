package mock

import (
	"context"
	"fmt"
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

// syntheticNeighbors returns nodes and edges that are "hidden" from the base graph
// and only revealed when expanding a specific seed node.
var syntheticNeighbors = map[string][]struct {
	node *graphv1.Node
	edge *graphv1.Edge
}{
	"cn_workspace_strategy": {
		{
			node: &graphv1.Node{
				Id: "exp_strategy_vision", CanonicalNodeId: "exp_strategy_vision",
				Scope: graphv1.GraphProjectionScope_GRAPH_PROJECTION_SCOPE_CANONICAL,
				Label: "Long-term Vision", Level: 2,
				Category:    graphv1.NodeCategory_NODE_CATEGORY_CONCEPT,
				Description: "The multi-year horizon that the workspace strategy is driving toward.",
			},
			edge: &graphv1.Edge{
				Id: "ed_strategy_vision", Scope: graphv1.GraphProjectionScope_GRAPH_PROJECTION_SCOPE_CANONICAL,
				Source: "cn_workspace_strategy", Target: "exp_strategy_vision",
				Type: graphv1.EdgeType_EDGE_TYPE_HIERARCHICAL,
			},
		},
		{
			node: &graphv1.Node{
				Id: "exp_strategy_risk", CanonicalNodeId: "exp_strategy_risk",
				Scope: graphv1.GraphProjectionScope_GRAPH_PROJECTION_SCOPE_CANONICAL,
				Label: "Risk Factors", Level: 2,
				Category:    graphv1.NodeCategory_NODE_CATEGORY_CLAIM,
				Description: "Known risks that could impact strategic outcomes.",
			},
			edge: &graphv1.Edge{
				Id: "ed_strategy_risk", Scope: graphv1.GraphProjectionScope_GRAPH_PROJECTION_SCOPE_CANONICAL,
				Source: "cn_workspace_strategy", Target: "exp_strategy_risk",
				Type: graphv1.EdgeType_EDGE_TYPE_RELATED_TO,
			},
		},
		{
			node: &graphv1.Node{
				Id: "exp_strategy_action", CanonicalNodeId: "exp_strategy_action",
				Scope: graphv1.GraphProjectionScope_GRAPH_PROJECTION_SCOPE_CANONICAL,
				Label: "Priority Actions", Level: 2,
				Category:    graphv1.NodeCategory_NODE_CATEGORY_ACTION,
				Description: "Immediate next steps derived from the workspace strategy.",
			},
			edge: &graphv1.Edge{
				Id: "ed_strategy_action", Scope: graphv1.GraphProjectionScope_GRAPH_PROJECTION_SCOPE_CANONICAL,
				Source: "cn_workspace_strategy", Target: "exp_strategy_action",
				Type: graphv1.EdgeType_EDGE_TYPE_CAUSES,
			},
		},
	},
	"cn_workspace_evidence": {
		{
			node: &graphv1.Node{
				Id: "exp_evidence_gap", CanonicalNodeId: "exp_evidence_gap",
				Scope: graphv1.GraphProjectionScope_GRAPH_PROJECTION_SCOPE_CANONICAL,
				Label: "Evidence Gaps", Level: 3,
				Category:    graphv1.NodeCategory_NODE_CATEGORY_CLAIM,
				Description: "Areas where supporting evidence is insufficient or missing.",
			},
			edge: &graphv1.Edge{
				Id: "ed_evidence_gap", Scope: graphv1.GraphProjectionScope_GRAPH_PROJECTION_SCOPE_CANONICAL,
				Source: "cn_workspace_evidence", Target: "exp_evidence_gap",
				Type: graphv1.EdgeType_EDGE_TYPE_RELATED_TO,
			},
		},
		{
			node: &graphv1.Node{
				Id: "exp_evidence_conflict", CanonicalNodeId: "exp_evidence_conflict",
				Scope: graphv1.GraphProjectionScope_GRAPH_PROJECTION_SCOPE_CANONICAL,
				Label: "Conflicting Claims", Level: 3,
				Category:    graphv1.NodeCategory_NODE_CATEGORY_EVIDENCE,
				Description: "Evidence items that contradict each other across documents.",
			},
			edge: &graphv1.Edge{
				Id: "ed_evidence_conflict", Scope: graphv1.GraphProjectionScope_GRAPH_PROJECTION_SCOPE_CANONICAL,
				Source: "cn_workspace_evidence", Target: "exp_evidence_conflict",
				Type: graphv1.EdgeType_EDGE_TYPE_CONTRADICTS,
			},
		},
	},
	"cn_workspace_metrics": {
		{
			node: &graphv1.Node{
				Id: "exp_metrics_threshold", CanonicalNodeId: "exp_metrics_threshold",
				Scope: graphv1.GraphProjectionScope_GRAPH_PROJECTION_SCOPE_CANONICAL,
				Label: "Threshold Targets", Level: 3,
				Category:    graphv1.NodeCategory_NODE_CATEGORY_METRIC,
				Description: "Quantitative thresholds that define success for each KPI.",
			},
			edge: &graphv1.Edge{
				Id: "ed_metrics_threshold", Scope: graphv1.GraphProjectionScope_GRAPH_PROJECTION_SCOPE_CANONICAL,
				Source: "cn_workspace_metrics", Target: "exp_metrics_threshold",
				Type: graphv1.EdgeType_EDGE_TYPE_MEASURED_BY,
			},
		},
	},
}

func (r *GraphRepository) ExpandNeighbors(ctx context.Context, workspaceID string, seedNodeID string, maxDepth uint32, limitPerHop uint32, edgeTypeFilters []graphv1.EdgeType) (*graphv1.Graph, error) {
	if workspaceID == "" {
		workspaceID = "00000000-0000-4000-8000-000000000001"
	}
	if maxDepth == 0 {
		maxDepth = 1
	}
	if limitPerHop == 0 {
		limitPerHop = 10
	}

	fullGraph, err := r.GetGraph(ctx, workspaceID, "")
	if err != nil {
		return nil, err
	}

	// Build adjacency: nodeID -> list of (neighborID, edge)
	type adjEntry struct {
		neighborID string
		edge       *graphv1.Edge
	}
	adj := make(map[string][]adjEntry)
	for _, edge := range fullGraph.Edges {
		if len(edgeTypeFilters) == 0 || containsEdgeType(edgeTypeFilters, edge.Type) {
			adj[edge.Source] = append(adj[edge.Source], adjEntry{edge.Target, edge})
			adj[edge.Target] = append(adj[edge.Target], adjEntry{edge.Source, edge})
		}
	}

	// BFS from seed
	nodeByID := make(map[string]*graphv1.Node, len(fullGraph.Nodes))
	for _, n := range fullGraph.Nodes {
		nodeByID[n.Id] = n
	}

	visitedNodes := map[string]bool{seedNodeID: true}
	visitedEdges := map[string]bool{}
	frontier := []string{seedNodeID}
	resultNodes := []*graphv1.Node{}
	resultEdges := []*graphv1.Edge{}

	if n, ok := nodeByID[seedNodeID]; ok {
		resultNodes = append(resultNodes, n)
	}

	for depth := uint32(0); depth < maxDepth && len(frontier) > 0; depth++ {
		next := []string{}
		for _, nodeID := range frontier {
			count := uint32(0)
			for _, entry := range adj[nodeID] {
				if count >= limitPerHop {
					break
				}
				if !visitedNodes[entry.neighborID] {
					visitedNodes[entry.neighborID] = true
					if n, ok := nodeByID[entry.neighborID]; ok {
						resultNodes = append(resultNodes, n)
					}
					next = append(next, entry.neighborID)
					count++
				}
				edgeKey := entry.edge.Id
				if !visitedEdges[edgeKey] {
					visitedEdges[edgeKey] = true
					resultEdges = append(resultEdges, entry.edge)
				}
			}
		}
		frontier = next
	}

	// Append synthetic neighbors for the seed node (hidden from base graph)
	if synthetics, ok := syntheticNeighbors[seedNodeID]; ok {
		for i, s := range synthetics {
			if uint32(i) >= limitPerHop {
				break
			}
			if len(edgeTypeFilters) > 0 && !containsEdgeType(edgeTypeFilters, s.edge.Type) {
				continue
			}
			resultNodes = append(resultNodes, s.node)
			resultEdges = append(resultEdges, s.edge)
		}
	}

	// For document nodes, generate synthetic neighbors dynamically
	if strings.HasPrefix(seedNodeID, "doc_node_") {
		docID := strings.TrimPrefix(seedNodeID, "doc_node_")
		slug := slugFilename(docID)
		resultNodes = append(resultNodes, &graphv1.Node{
			Id: fmt.Sprintf("exp_%s_summary", slug), CanonicalNodeId: fmt.Sprintf("exp_%s_summary", slug),
			DocumentId:  docID,
			Scope:       graphv1.GraphProjectionScope_GRAPH_PROJECTION_SCOPE_DOCUMENT,
			Label:       "Auto-generated Summary",
			Level:       4,
			Category:    graphv1.NodeCategory_NODE_CATEGORY_CONCEPT,
			Description: "Key themes automatically extracted from this document.",
		})
		resultEdges = append(resultEdges, &graphv1.Edge{
			Id:     fmt.Sprintf("ed_exp_%s_summary", slug),
			Source: seedNodeID, Target: fmt.Sprintf("exp_%s_summary", slug),
			Scope: graphv1.GraphProjectionScope_GRAPH_PROJECTION_SCOPE_DOCUMENT,
			Type:  graphv1.EdgeType_EDGE_TYPE_HIERARCHICAL,
		})
	}

	return &graphv1.Graph{
		WorkspaceId:   workspaceID,
		CrossDocument: true,
		Nodes:         resultNodes,
		Edges:         resultEdges,
	}, nil
}

func containsEdgeType(filters []graphv1.EdgeType, t graphv1.EdgeType) bool {
	for _, f := range filters {
		if f == t {
			return true
		}
	}
	return false
}

func slugFilename(filename string) string {
	slug := strings.ToLower(filename)
	replacer := strings.NewReplacer(" ", "-", ".", "-", "_", "-")
	return replacer.Replace(slug)
}
