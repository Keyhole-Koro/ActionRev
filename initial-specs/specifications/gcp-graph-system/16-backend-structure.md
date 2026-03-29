# 16. Backend Directory Structure

## 概要

Go バックエンドはレイヤードアーキテクチャを採用する。各レイヤーは一方向にのみ依存し、循環参照を禁止する。

```
handler → service → repository / infra
              ↑
           domain (全レイヤーから参照可)
```

proto 生成型はワイヤー形式として `handler` 層でのみ使用する。ビジネスロジックは `domain` 型で扱い、`handler` でマッピングする。

---

## ディレクトリ構造

```
backend/
├── cmd/
│   └── server/
│       └── main.go            # DI 組み立て・サーバー起動
│
├── internal/
│   ├── domain/                # ドメインモデル・値オブジェクト (依存なし)
│   │   ├── document.go
│   │   ├── node.go
│   │   ├── edge.go
│   │   ├── workspace.go
│   │   ├── user.go
│   │   └── errors.go          # ドメインエラー定義
│   │
│   ├── handler/               # Connect RPC ハンドラ (薄い層)
│   │   │                      # 責務: 入力バリデーション・proto↔domain 変換・service 呼び出し
│   │   ├── document.go
│   │   ├── graph.go
│   │   ├── workspace.go
│   │   ├── user.go
│   │   ├── billing.go
│   │   └── dev.go             # /dev/stats (dev ロール限定)
│   │
│   ├── service/               # ビジネスロジック層
│   │   ├── document.go        # ドキュメント登録・ステータス管理
│   │   ├── graph.go           # GetGraph・ExpandNeighbors・FindPaths
│   │   ├── workspace.go       # ワークスペース CRUD・メンバー管理
│   │   ├── user.go            # SyncUser・GetMe
│   │   ├── billing.go         # Stripe Checkout / Portal / Webhook 処理
│   │   └── normalization.go   # 正規化ツール管理・承認フロー
│   │
│   ├── pipeline/              # AI 抽出パイプライン (非同期ジョブから呼ばれる)
│   │   ├── pipeline.go        # パイプライン全体のオーケストレーション
│   │   ├── chunker.go         # セマンティックチャンク分割
│   │   ├── extractor.go       # Pass1 (chunk 単位) + Pass2 (ドキュメント統合) Gemini 呼び出し
│   │   ├── integrator.go      # ノード重複統合 (edit distance + embedding)
│   │   └── summarizer.go      # summary_html 生成
│   │
│   ├── repository/            # データアクセス層 (インターフェース + BigQuery/GCS 実装)
│   │   ├── interfaces.go      # リポジトリインターフェース定義
│   │   ├── bigquery/
│   │   │   ├── document.go
│   │   │   ├── node.go
│   │   │   ├── edge.go
│   │   │   ├── workspace.go
│   │   │   └── user.go
│   │   ├── spannergraph/
│   │   │   └── graph.go       # 近傍探索・経路検索・graph 同期
│   │   └── gcs/
│   │       └── upload.go      # 署名付き URL 発行・オブジェクト操作
│   │
│   ├── infra/                 # 外部サービスクライアント
│   │   ├── gemini/
│   │   │   ├── client.go      # Vertex AI Gemini 呼び出し
│   │   │   └── cache.go       # プロンプトキャッシュ (GEMINI_CACHE_ENABLED)
│   │   ├── firebase/
│   │   │   └── auth.go        # ID Token 検証・UID 取得
│   │   ├── tasks/
│   │   │   └── client.go      # Cloud Tasks エンキュー
│   │   ├── stripe/
│   │   │   ├── client.go      # Checkout / Portal セッション作成
│   │   │   └── webhook.go     # Webhook イベント処理
│   │   ├── discord/
│   │   │   └── webhook.go     # 通知送信 (8 イベント)
│   │   ├── spanner/
│   │   │   └── graph.go       # Spanner Graph クライアント
│   │   └── sandbox/
│   │       └── runner.go      # Cloud Run Jobs 起動・結果ポーリング
│   │
│   └── middleware/
│       ├── auth.go            # Firebase ID Token 検証・role チェック
│       ├── logging.go         # リクエストログ (Cloud Logging 形式)
│       └── recovery.go        # panic リカバリ
│
├── gen/                       # proto 生成コード (make generate で再生成)
│   └── actionrev/graph/v1/
│
├── Dockerfile
├── Dockerfile.dev             # Air ホットリロード
├── .air.toml
├── go.mod
├── go.sum
└── Makefile
```

---

## パッケージ設計方針

### domain

- 外部パッケージへの依存を持たない純粋な Go 型
- ドメインエラーは `errors.go` で sentinel error として定義する

```go
// errors.go
var (
    ErrNotFound         = errors.New("not found")
    ErrPermissionDenied = errors.New("permission denied")
    ErrStorageQuotaExceeded = errors.New("storage quota exceeded")
    ErrFileTooLarge     = errors.New("file too large")
)
```

### handler

- proto 型 ↔ domain 型の変換のみ行う
- ビジネスロジックを持たない
- `domain.ErrNotFound` → `connect.CodeNotFound` のマッピングを共通関数で処理する

```go
func toConnectError(err error) error {
    switch {
    case errors.Is(err, domain.ErrNotFound):
        return connect.NewError(connect.CodeNotFound, err)
    case errors.Is(err, domain.ErrPermissionDenied):
        return connect.NewError(connect.CodePermissionDenied, err)
    // ...
    }
}
```

### service

- `repository.interfaces.go` で定義したインターフェースのみに依存する (BigQuery 実装には直接依存しない)
- トランザクション相当の操作 (複数テーブル書き込み) はサービス層で調整する

### repository/interfaces.go

サービス層が依存するインターフェースをここに集約する。テスト時はモック実装に差し替え可能。

```go
type DocumentRepository interface {
    Create(ctx context.Context, doc domain.Document) error
    GetByID(ctx context.Context, id string) (domain.Document, error)
    UpdateStatus(ctx context.Context, id string, status domain.DocumentStatus) error
    ListByWorkspace(ctx context.Context, workspaceID string, opts ListOptions) ([]domain.Document, error)
}

type NodeRepository interface {
    BatchUpsert(ctx context.Context, nodes []domain.Node) error
    ListByDocument(ctx context.Context, documentID string) ([]domain.Node, error)
}
// ...

type GraphQueryRepository interface {
    ExpandNeighbors(ctx context.Context, seedNodeID string, depth int, edgeTypes []domain.EdgeType, limitPerHop int) (domain.Graph, error)
    FindPaths(ctx context.Context, sourceNodeID string, targetNodeID string, maxDepth int, edgeTypes []domain.EdgeType, limit int) (domain.Graph, []domain.GraphPath, error)
    SyncCanonicalGraph(ctx context.Context, nodes []domain.Node, edges []domain.Edge) error
}
```

### pipeline

- Cloud Tasks ジョブのエントリポイントは `handler/` に置き、実処理を `pipeline/` に委譲する
- パイプライン内部の各ステップは独立した関数として定義し、単体テスト可能にする

### infra/gemini/cache.go

```go
// プロンプトの SHA256 ハッシュをキーにレスポンスをファイルキャッシュ
// GEMINI_CACHE_ENABLED=false の場合はキャッシュを完全にバイパスする
type CachedClient struct {
    inner    GeminiClient
    cacheDir string
    enabled  bool
}
```

---

## 依存性注入 (DI)

`wire` や `fx` は使用せず、`cmd/server/main.go` でコンストラクタを手動で組み立てる。

```go
// main.go (概略)
func main() {
    cfg := config.Load()

    bqClient   := bigquery.NewClient(cfg)
    gcsClient  := gcs.NewClient(cfg)
    gemini     := gemini.NewCachedClient(cfg)
    firebaseAuth := firebase.NewAuth(cfg)
    tasksClient := tasks.NewClient(cfg)

    docRepo  := bqrepo.NewDocumentRepository(bqClient)
    nodeRepo := bqrepo.NewNodeRepository(bqClient)

    docService  := service.NewDocumentService(docRepo, gcsClient, tasksClient)
    graphService := service.NewGraphService(nodeRepo, edgeRepo)
    userService := service.NewUserService(userRepo)

    mux := http.NewServeMux()
    mux.Handle(documentv1connect.NewDocumentServiceHandler(
        handler.NewDocumentHandler(docService),
        connect.WithInterceptors(middleware.NewAuthInterceptor(firebaseAuth)),
    ))
    // ...
}
```

---

## Makefile ターゲット

```makefile
generate:   # proto → Go コード生成 (buf generate)
build:      # go build ./cmd/server
test:       # go test ./...
lint:       # golangci-lint run
run:        # ローカル起動 (docker compose up)
```

---

## 命名規則

| 対象 | 規則 | 例 |
|---|---|---|
| パッケージ名 | 単数形・小文字 | `service`, `handler`, `repository` |
| インターフェース | `〜er` または `〜Repository` | `DocumentRepository`, `GeminiClient` |
| コンストラクタ | `New〜` | `NewDocumentService` |
| エラー変数 | `Err〜` | `ErrNotFound` |
| ファイル名 | スネークケース | `document_repository.go` → ただし1ファイル1責務なら単語のみ `document.go` |
