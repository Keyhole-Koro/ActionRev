package main

import (
	"log/slog"
	"net/http"
	"os"

	"connectrpc.com/connect"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"

	"github.com/keyhole-koro/actionrev/internal/middleware"

	// 生成コード (buf generate 後に有効になる)
	// documentv1connect "github.com/keyhole-koro/actionrev/gen/actionrev/graph/v1/graphv1connect"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	cfg := loadConfig()

	// --- Infrastructure ---
	// bqClient     := mustBigQuery(cfg)
	// gcsClient    := mustGCS(cfg)
	// geminiClient := mustGemini(cfg)
	// firebaseAuth := mustFirebaseAuth(cfg)
	// tasksClient  := mustTasks(cfg)

	// --- Repositories ---
	// docRepo   := bqrepo.NewDocumentRepository(bqClient)
	// nodeRepo  := bqrepo.NewNodeRepository(bqClient)
	// edgeRepo  := bqrepo.NewEdgeRepository(bqClient)
	// wsRepo    := bqrepo.NewWorkspaceRepository(bqClient)
	// userRepo  := bqrepo.NewUserRepository(bqClient)
	// uploadRepo := gcsrepo.NewUploadRepository(gcsClient, cfg.GCSBucket)

	// --- Services ---
	// docService  := service.NewDocumentService(docRepo, uploadRepo, tasksClient)
	// graphService := service.NewGraphService(nodeRepo, edgeRepo)
	// wsService   := service.NewWorkspaceService(wsRepo)
	// userService := service.NewUserService(userRepo)

	// --- Auth Interceptor ---
	// authInterceptor := middleware.NewAuthInterceptor(firebaseAuth)
	_ = middleware.NewAuthInterceptor // suppress unused import until wired

	mux := http.NewServeMux()

	// mux.Handle(documentv1connect.NewDocumentServiceHandler(
	// 	handler.NewDocumentHandler(docService),
	// 	connect.WithInterceptors(authInterceptor),
	// ))
	_ = connect.WithInterceptors // suppress unused import

	// ヘルスチェック
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	addr := ":" + cfg.Port
	slog.Info("starting server", "addr", addr)

	if err := http.ListenAndServe(addr, h2c.NewHandler(mux, &http2.Server{})); err != nil {
		slog.Error("server error", "err", err)
		os.Exit(1)
	}
}

type config struct {
	Port              string
	GCPProjectID      string
	GCSBucket         string
	BigQueryProjectID string
	BigQueryDataset   string
	GeminiCacheDir    string
	GeminiCacheEnabled bool
}

func loadConfig() config {
	return config{
		Port:               getEnv("PORT", "8080"),
		GCPProjectID:       getEnv("GCP_PROJECT_ID", ""),
		GCSBucket:          getEnv("GCS_BUCKET", ""),
		BigQueryProjectID:  getEnv("BIGQUERY_PROJECT_ID", ""),
		BigQueryDataset:    getEnv("BIGQUERY_DATASET", "graph"),
		GeminiCacheDir:     getEnv("GEMINI_CACHE_DIR", "/app/.gemini-cache"),
		GeminiCacheEnabled: getEnv("GEMINI_CACHE_ENABLED", "false") == "true",
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
