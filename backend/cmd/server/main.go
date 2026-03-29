package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"

	connect "connectrpc.com/connect"
	graphv1 "github.com/synthify/backend/gen/synthify/graph/v1"
	"github.com/synthify/backend/gen/synthify/graph/v1/graphv1connect"
)

type rootResponse struct {
	Name   string `json:"name"`
	Status string `json:"status"`
}

type graphServiceServer struct{}

func (s *graphServiceServer) GetGraph(
	_ context.Context,
	req *connect.Request[graphv1.GetGraphRequest],
) (*connect.Response[graphv1.GetGraphResponse], error) {
	workspaceID := req.Msg.WorkspaceId
	if workspaceID == "" {
		workspaceID = "ws_demo"
	}
	documentID := req.Msg.DocumentId
	if documentID == "" {
		documentID = "doc_demo"
	}

	graph := &graphv1.Graph{
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
	}

	return connect.NewResponse(&graphv1.GetGraphResponse{
		DocumentId: documentID,
		Graph:      graph,
	}), nil
}

func (s *graphServiceServer) ExpandNeighbors(
	_ context.Context,
	_ *connect.Request[graphv1.ExpandNeighborsRequest],
) (*connect.Response[graphv1.ExpandNeighborsResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, os.ErrInvalid)
}

func (s *graphServiceServer) FindPaths(
	_ context.Context,
	_ *connect.Request[graphv1.FindPathsRequest],
) (*connect.Response[graphv1.FindPathsResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, os.ErrInvalid)
}

func corsMiddleware(next http.Handler) http.Handler {
	allowedOrigins := splitCSV(os.Getenv("CORS_ALLOWED_ORIGINS"))
	if len(allowedOrigins) == 0 {
		allowedOrigins = []string{"http://localhost:5173"}
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if origin != "" && isAllowedOrigin(origin, allowedOrigins) {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Vary", "Origin")
			w.Header().Set("Access-Control-Allow-Methods", "GET,POST,OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type,Connect-Protocol-Version,Connect-Timeout-Ms,X-User-Agent")
		}

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func splitCSV(v string) []string {
	if v == "" {
		return nil
	}
	parts := strings.Split(v, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}

func isAllowedOrigin(origin string, allowed []string) bool {
	for _, candidate := range allowed {
		if candidate == origin {
			return true
		}
	}
	return false
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	})

	path, handler := graphv1connect.NewGraphServiceHandler(&graphServiceServer{})
	mux.Handle(path, handler)

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(rootResponse{
			Name:   "synthify-backend",
			Status: "running",
		})
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("synthify backend listening on :%s", port)
	if err := http.ListenAndServe(":"+port, corsMiddleware(mux)); err != nil {
		log.Fatal(err)
	}
}
