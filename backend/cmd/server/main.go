package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/synthify/backend/gen/synthify/graph/v1/graphv1connect"
	"github.com/synthify/backend/internal/handler"
	mockrepo "github.com/synthify/backend/internal/repository/mock"
	"github.com/synthify/backend/internal/service"
)

type rootResponse struct {
	Name   string `json:"name"`
	Status string `json:"status"`
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
			w.Header().Set("Access-Control-Allow-Methods", "GET,POST,PUT,OPTIONS")
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
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	})

	documentRepo := mockrepo.NewDocumentRepository()
	documentService := service.NewDocumentService(documentRepo, publicBackendBaseURL(port))
	documentHandler := handler.NewDocumentHandler(documentService)
	mockUploadHandler := handler.NewMockUploadHandler(documentService)
	graphRepo := mockrepo.NewGraphRepository(documentRepo)
	graphService := service.NewGraphService(graphRepo)
	graphHandler := handler.NewGraphHandler(graphService)
	workspaceRepo := mockrepo.NewWorkspaceRepository()
	workspaceService := service.NewWorkspaceService(workspaceRepo)
	workspaceHandler := handler.NewWorkspaceHandler(workspaceService)

	path, connectHandler := graphv1connect.NewGraphServiceHandler(graphHandler)
	mux.Handle(path, connectHandler)
	documentPath, documentConnectHandler := graphv1connect.NewDocumentServiceHandler(documentHandler)
	mux.Handle(documentPath, documentConnectHandler)
	mux.Handle("/mock/uploads/", mockUploadHandler)
	workspacePath, workspaceConnectHandler := graphv1connect.NewWorkspaceServiceHandler(workspaceHandler)
	mux.Handle(workspacePath, workspaceConnectHandler)

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(rootResponse{
			Name:   "synthify-backend",
			Status: "running",
		})
	})

	log.Printf("synthify backend listening on :%s", port)
	if err := http.ListenAndServe(":"+port, corsMiddleware(mux)); err != nil {
		log.Fatal(err)
	}
}

func publicBackendBaseURL(port string) string {
	if baseURL := strings.TrimSpace(os.Getenv("PUBLIC_BASE_URL")); baseURL != "" {
		return strings.TrimRight(baseURL, "/")
	}

	return "http://localhost:" + port
}
