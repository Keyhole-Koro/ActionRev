// migrate - BigQuery スキーママイグレーションランナー
//
// Usage:
//   go run ./cmd/migrate -project=my-project -dataset=graph -dir=../../migrations/bigquery
//   go run ./cmd/migrate -project=my-project -dry-run
package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"cloud.google.com/go/bigquery"
	"google.golang.org/api/iterator"
)

func main() {
	project := flag.String("project", os.Getenv("BIGQUERY_PROJECT_ID"), "GCP プロジェクト ID")
	dataset := flag.String("dataset", envOr("BIGQUERY_DATASET", "graph"), "BigQuery データセット ID")
	dir := flag.String("dir", "migrations/bigquery", "マイグレーションファイルのディレクトリ")
	dryRun := flag.Bool("dry-run", false, "pending を表示するだけ (適用しない)")
	flag.Parse()

	if *project == "" {
		slog.Error("-project が必要です")
		os.Exit(1)
	}

	ctx := context.Background()
	if err := run(ctx, *project, *dataset, *dir, *dryRun); err != nil {
		slog.Error("migration failed", "err", err)
		os.Exit(1)
	}
}

func run(ctx context.Context, project, dataset, dir string, dryRun bool) error {
	client, err := bigquery.NewClient(ctx, project)
	if err != nil {
		return fmt.Errorf("bigquery client: %w", err)
	}
	defer client.Close()

	m := &migrator{
		client:  client,
		project: project,
		dataset: dataset,
	}

	if err := m.ensureMigrationsTable(ctx); err != nil {
		return fmt.Errorf("ensure migrations table: %w", err)
	}

	applied, err := m.appliedVersions(ctx)
	if err != nil {
		return fmt.Errorf("fetch applied versions: %w", err)
	}

	files, err := loadMigrationFiles(dir)
	if err != nil {
		return fmt.Errorf("load migration files: %w", err)
	}

	pending := filterPending(files, applied)

	if len(pending) == 0 {
		slog.Info("no pending migrations")
		return nil
	}

	slog.Info("pending migrations", "count", len(pending))
	for _, f := range pending {
		slog.Info("  pending", "version", f.version, "name", f.name)
	}

	if dryRun {
		slog.Info("dry-run: skipping apply")
		return nil
	}

	for _, f := range pending {
		slog.Info("applying migration", "version", f.version, "name", f.name)
		if err := m.apply(ctx, f); err != nil {
			return fmt.Errorf("apply migration %d (%s): %w", f.version, f.name, err)
		}
		slog.Info("applied", "version", f.version)
	}

	slog.Info("all migrations applied", "count", len(pending))
	return nil
}

// ---------------------------------------------------------------------------
// migrator
// ---------------------------------------------------------------------------

type migrator struct {
	client  *bigquery.Client
	project string
	dataset string
}

const createMigrationsTable = `
CREATE TABLE IF NOT EXISTS ` + "`{project}.{dataset}.schema_migrations`" + ` (
  version     INT64     NOT NULL,
  description STRING,
  applied_at  TIMESTAMP NOT NULL
)`

func (m *migrator) ensureMigrationsTable(ctx context.Context) error {
	sql := m.interpolate(createMigrationsTable)
	q := m.client.Query(sql)
	job, err := q.Run(ctx)
	if err != nil {
		return err
	}
	_, err = job.Wait(ctx)
	return err
}

func (m *migrator) appliedVersions(ctx context.Context) (map[int]bool, error) {
	sql := m.interpolate("SELECT version FROM `{project}.{dataset}.schema_migrations`")
	q := m.client.Query(sql)
	it, err := q.Read(ctx)
	if err != nil {
		return nil, err
	}

	applied := map[int]bool{}
	for {
		var row struct {
			Version int64 `bigquery:"version"`
		}
		err := it.Next(&row)
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		applied[int(row.Version)] = true
	}
	return applied, nil
}

func (m *migrator) apply(ctx context.Context, f migrationFile) error {
	// プレースホルダーを置換
	sql := m.interpolate(f.sql)

	// セミコロンで分割して複数ステートメントを順に実行
	stmts := splitStatements(sql)
	for _, stmt := range stmts {
		stmt = strings.TrimSpace(stmt)
		if stmt == "" {
			continue
		}
		q := m.client.Query(stmt)
		job, err := q.Run(ctx)
		if err != nil {
			return fmt.Errorf("run statement: %w\nSQL: %s", err, stmt)
		}
		if _, err := job.Wait(ctx); err != nil {
			return fmt.Errorf("wait statement: %w\nSQL: %s", err, stmt)
		}
	}

	// schema_migrations に記録
	insertSQL := m.interpolate(fmt.Sprintf(
		"INSERT INTO `{project}.{dataset}.schema_migrations` (version, description, applied_at) VALUES (%d, '%s', CURRENT_TIMESTAMP())",
		f.version, escapeSingleQuote(f.name),
	))
	q := m.client.Query(insertSQL)
	job, err := q.Run(ctx)
	if err != nil {
		return err
	}
	_, err = job.Wait(ctx)
	return err
}

func (m *migrator) interpolate(sql string) string {
	sql = strings.ReplaceAll(sql, "{project}", m.project)
	sql = strings.ReplaceAll(sql, "{dataset}", m.dataset)
	return sql
}

// ---------------------------------------------------------------------------
// Migration file loading
// ---------------------------------------------------------------------------

type migrationFile struct {
	version int
	name    string
	sql     string
}

func loadMigrationFiles(dir string) ([]migrationFile, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("read dir %s: %w", dir, err)
	}

	var files []migrationFile
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".sql") {
			continue
		}

		parts := strings.SplitN(e.Name(), "_", 2)
		if len(parts) < 2 {
			continue
		}
		version, err := strconv.Atoi(parts[0])
		if err != nil {
			continue
		}
		name := strings.TrimSuffix(parts[1], ".sql")

		content, err := os.ReadFile(filepath.Join(dir, e.Name()))
		if err != nil {
			return nil, fmt.Errorf("read file %s: %w", e.Name(), err)
		}

		files = append(files, migrationFile{
			version: version,
			name:    name,
			sql:     string(content),
		})
	}

	sort.Slice(files, func(i, j int) bool {
		return files[i].version < files[j].version
	})
	return files, nil
}

func filterPending(files []migrationFile, applied map[int]bool) []migrationFile {
	var pending []migrationFile
	for _, f := range files {
		if !applied[f.version] {
			pending = append(pending, f)
		}
	}
	return pending
}

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

func splitStatements(sql string) []string {
	// コメント行 (-- ...) を除去してからセミコロンで分割
	var lines []string
	for _, line := range strings.Split(sql, "\n") {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "--") {
			continue
		}
		lines = append(lines, line)
	}
	return strings.Split(strings.Join(lines, "\n"), ";")
}

func escapeSingleQuote(s string) string {
	return strings.ReplaceAll(s, "'", "\\'")
}

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

// time パッケージを使用していることを示す (実際には BigQuery の CURRENT_TIMESTAMP を使用)
var _ = time.Now
